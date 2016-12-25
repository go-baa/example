// Package setting 提供了一些运行时选项的配置，提供了默认值和通过配置文件设置的初始化，运行运行时改变
package setting

import "gopkg.in/baa.v1"

var (
	// AppName 应用名称
	AppName string
	// AppVersion 应用版本
	AppVersion string
	// AppURL 应用URL
	AppURL string
	// Debug 调试模式
	Debug bool
)

func init() {
	if baa.Env != baa.PROD {
		Debug = true
	}
	v, err := Config.GetBool("debug")
	if err == nil {
		Debug = v
	}
	v, err = Config.GetBool("app.debug")
	if err == nil {
		Debug = v
	}
	AppName = Config.MustString("app.name", "appName")
	AppVersion = Config.MustString("app.version", "unknown")
	AppURL = Config.MustString("app.url", "")
}
