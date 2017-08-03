package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/go-baa/setting"
	"github.com/go-baa/example/blog/modules/util"

	"gopkg.in/baa.v1"
)

// Funcs 模板函数库
func Funcs(b *baa.Baa) template.FuncMap {
	return map[string]interface{}{
		"version": func() string {
			return "v0.1 beta"
		},
		"config":        config,
		"assets":        assets,
		"assetsUrl":     assetsURL,
		"secondToTime":  util.SecondToTime,
		"stylesheetTag": stylesheetTag,
		"scriptTag":     scriptTag,
		"time":          timeToLayout,
		"baaEq":         util.Eq,
		"baaNeq":        util.Neq,
		"baaGt":         util.Gt,
		"baaLt":         util.Lt,
		"baaEmpty":      util.Empty,
		"url": func(name string, params ...interface{}) string {
			n := make([]interface{}, len(params))
			for k, v := range params {
				n[k] = fmt.Sprint(v)
			}
			return b.URLFor(name, n...)
		},
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"html": func(params ...interface{}) template.HTML {
			return template.HTML(fmt.Sprint(params...))
		},
		"strcut": func(length int, dot string, str string) template.HTML {
			return template.HTML(util.StrNatCut(str, length, dot))
		},
		"trim": func(s string) string {
			return strings.TrimSpace(s)
		},
		"trimToLine": func(s string) string {
			s = strings.Replace(s, "\n", "", -1)
			s = strings.Replace(s, "\t", "", -1)
			return strings.TrimSpace(s)
		},
	}
}

func config(args ...string) string {
	key := args[0]
	dft := ""
	if len(args) == 2 {
		dft = args[1]
	}
	if len(key) == 0 {
		return dft
	}
	return setting.Config.MustString(key, dft)
}

// 解析VUE构建文件
var manifest map[string]string

func loadManifest() {
	data, err := util.ReadFile(setting.Config.MustString("assets.buildPath", "") + "rev-manifest.json")
	if err == nil {
		json.Unmarshal(data, &manifest)
	}
}

func assets(path string) string {
	if val, ok := manifest[path]; ok {
		return setting.Config.MustString("assets.baseUri", "") + val
	}
	return ""
}

func assetsURL(str ...string) string {
	path := strings.Join(str, "")
	if val, ok := manifest["hash"]; ok {
		if strings.Contains(path, "?") {
			return path + "&" + val
		}
		return path + "?" + val
	}
	return path
}

func stylesheetTag(name string) template.HTML {
	uri := assets(name)
	tag := ""
	if len(uri) > 0 {
		tag = "<link rel=\"stylesheet\" href=\"" + uri + "\">"
	}
	return template.HTML(tag)
}

func scriptTag(name string) template.HTML {
	uri := assets(name)
	tag := ""
	if len(uri) > 0 {
		tag = "<script src=\"" + uri + "\"></script>"
	}
	return template.HTML(tag)
}

func timeToLayout(layout string, t time.Time) string {
	if t.Unix() > 0 {
		return t.Format(layout)
	}
	return ""
}

func init() {
	loadManifest()
}
