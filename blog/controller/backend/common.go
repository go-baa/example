package backend

import (
	"github.com/go-baa/example/blog/controller/base"

	"gopkg.in/baa.v1"
)

type commonController struct{}

// CommonController 公共控制器
var CommonController = commonController{}

// UploadResponse 上传响应
type UploadResponse struct {
	Success  bool   `json:"success"`
	Msg      string `json:"msg"`
	FilePath string `json:"file_path"`
}

// Upload 上传图片
func (t commonController) Upload(c *baa.Context) {
	res := new(UploadResponse)

	_, imgURI, err := base.UploadFile("post", "upload_file", c, "")
	if err != nil {
		res.Msg = err.Error()
		c.JSON(200, res)
		return
	}

	res.Success = true
	res.FilePath = imgURI

	c.JSON(200, res)
}
