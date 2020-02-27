package backend

import (
	"github.com/go-baa/baa"
)

type indexController struct{}

// IndexController 首页控制器
var IndexController = indexController{}

// Index 首页
func (t indexController) Index(c *baa.Context) {
	c.HTML(200, "backend/index")
}
