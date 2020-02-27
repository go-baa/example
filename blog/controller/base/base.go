package base

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/baa-middleware/session"
	"github.com/go-baa/baa"
	"github.com/go-baa/example/blog/model"
	"github.com/go-baa/example/blog/modules/util"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
)

// NormalReturn 标准返回格式
type NormalReturn struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewReturn 实例化一个新的返回结构
func NewReturn() *NormalReturn {
	re := new(NormalReturn)
	re.Data = make(map[string]interface{})
	return re
}

// GetUser 获取session用户
func GetUser(c *baa.Context) *model.AdminInfo {
	s := c.Get("session").(*session.Session)
	return s.Get("auth").(*model.AdminInfo)
}

// SetUser 设置用户session
func SetUser(c *baa.Context, user *model.AdminInfo) error {
	s := c.Get("session").(*session.Session)
	return s.Set("auth", user)
}

var (
	globalUploadBasePath  string // 全局上传根路径
	globalUploadBaseURI   string // 全局上传根URI
	globalUploadExtension string // 全局文件后缀限制
	globalUploadMaxsize   int64  // 全局文件大小限制
)

// UploadFile 从Http流中上传一个文件, 返回上传后的文件地址
// 默认上传获取到的第一个文件，如果指定 fieldName 仅上传指定的文件
// 允许指定一个附件路径，默认会上传到upload目录，如果有addonPath则附加
func UploadFile(uploadType, fieldName string, c *baa.Context, addonPath string) (string, string, error) {
	maxSize := setting.Config.MustInt64("upload."+uploadType+".maxsize", globalUploadMaxsize)
	err := c.Req.ParseMultipartForm(maxSize)
	if err != nil {
		return "", "", fmt.Errorf("超过上传限制，最大允许 %d m, %s", maxSize, err)
	}

	// 如果没有指定 文件字段，取第一个获取到的文件
	if fieldName == "" {
		for k := range c.Req.MultipartForm.File {
			fieldName = k
			break
		}
	}
	files := c.Req.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return "", "", fmt.Errorf("没有文件被上传")
	}

	ext := strings.ToLower(filepath.Ext(files[0].Filename))
	allowExt := strings.Split(setting.Config.MustString("upload."+uploadType+".extension", globalUploadExtension), ";")
	if ext == "" || util.InSlice(allowExt, ext, "string") == false {
		return "", "", fmt.Errorf("只允许上传指定的格式: %s", strings.Join(allowExt, ";"))
	}

	file, err := files[0].Open()
	if err != nil {
		return "", "", fmt.Errorf("文件读取失败: %s", err)
	}
	defer file.Close()

	uploadPath := strings.Trim(setting.Config.MustString("upload."+uploadType+".path", ""), "/")
	if uploadPath == "" {
		uploadPath = uploadType
	}
	uploadPath = globalUploadBasePath + "/" + uploadPath
	uploadPath, err = filepath.Abs(uploadPath)
	if err != nil {
		return "", "", fmt.Errorf("上传目录转化失败: %s", err)
	}
	if addonPath != "" {
		addonPath = "/" + strings.Trim(addonPath, "/")
	}
	err = util.MkdirAll(uploadPath + addonPath)
	if err != nil {
		return "", "", fmt.Errorf("上传目录创建失败: %s", err)
	}
	dstPath := uploadPath + addonPath + "/" + util.RandFileName() + ext
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", "", fmt.Errorf("文件创建失败: %s", err)
	}
	defer dst.Close()
	size, err := io.Copy(dst, file)
	if err != nil {
		return "", "", fmt.Errorf("文件写入失败: %s", err)
	}

	if setting.Debug {
		log.Infof("upload a file: %s fileSzie: %d saved to %s\n", files[0].Filename, size, dstPath)
	}

	uploadURI := strings.Trim(setting.Config.MustString("upload."+uploadType+".uri", ""), "/")
	if uploadURI == "" {
		uploadURI = uploadType
	}
	uploadURI = globalUploadBaseURI + "/" + uploadURI
	if len(uploadURI) == 0 {
		return dstPath, "", nil
	}
	uploadRelativeURI := dstPath[len(uploadPath):]
	return dstPath, uploadURI + uploadRelativeURI, nil
}

func init() {
	// 处理全局的上传配置
	globalUploadBasePath = setting.Config.MustString("upload.basePath", "")
	if len(globalUploadBasePath) > 1 {
		globalUploadBasePath = strings.TrimRight(globalUploadBasePath, "/")
	}
	globalUploadBaseURI = setting.Config.MustString("upload.baseUri", "")
	if len(globalUploadBaseURI) > 1 {
		globalUploadBaseURI = strings.TrimRight(globalUploadBaseURI, "/")
	}
	globalUploadExtension = setting.Config.MustString("upload.extension", "")
	globalUploadMaxsize = setting.Config.MustInt64("upload.maxsize", 1048576) // 1m
	gob.Register(NormalReturn{})
}
