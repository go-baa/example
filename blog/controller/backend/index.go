package backend

import (
	"gopkg.in/baa.v1"
)

type indexController struct{}

// IndexController 首页控制器
var IndexController = indexController{}

// Index 首页
func (t indexController) Index(c *baa.Context) {
	c.HTML(200, "backend/index")
}
