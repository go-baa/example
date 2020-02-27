package frontend

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/example/blog/controller/base"
	"github.com/go-baa/example/blog/model"
)

type contentController struct{}

// ContentController 内容控制器
var ContentController = contentController{}

// Search 列表
func (t contentController) Search(c *baa.Context) {
	page := c.ParamInt("page")
	if page == 0 {
		page = 1
	}
	pagesize := 5
	ret := base.NewReturn()
	rows, total, err := model.ContentModel.Search("", pagesize*(page-1), pagesize)
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	} else {
		ret.Data["items"] = rows
		ret.Data["cur_page"] = page
		ret.Data["more"] = 0
		if page*pagesize < total {
			ret.Data["more"] = 1
		}
	}

	c.JSON(200, ret)
}

// Show 详情
func (t contentController) Show(c *baa.Context) {
	cid := c.ParamInt("id")
	content, err := model.ContentModel.Get(cid)
	if err != nil {
		c.NotFound()
		return
	}
	info, _ := content.Format()
	c.Set("data", info)

	c.HTML(200, "frontend/post")
}
