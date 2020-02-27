package frontend

import (
	"github.com/go-baa/baa"
)

type indexController struct{}

// IndexController 首页控制器
var IndexController = indexController{}

// Index 首页
func (t indexController) Index(c *baa.Context) {
	c.HTML(200, "frontend/index")
}

// About 关于
func (t indexController) About(c *baa.Context) {
	c.HTML(200, "frontend/about")
}
