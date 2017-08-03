package backend

import (
	"github.com/go-baa/example/blog/controller/base"
	"github.com/go-baa/example/blog/model"
	"gopkg.in/baa.v1"
)

type contentController struct{}

// ContentController 内容控制器
var ContentController = contentController{}

// Search 列表
func (t contentController) Search(c *baa.Context) {
	start := c.QueryInt("start")
	length := c.QueryInt("length")
	keyword := c.Query("keyword")
	var limit, offset int
	if length > 0 {
		limit = length
		offset = start
	} else {
		limit = 15
	}
	ret := base.NewReturn()
	rows, total, err := model.ContentModel.Search(keyword, offset, limit)
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	} else {
		ret.Data["data"] = rows
		ret.Data["total"] = total
	}

	c.JSON(200, ret)
}

// Add 添加页面
func (t contentController) Add(c *baa.Context) {
	c.HTML(200, "backend/add")
}

// Create 添加请求
func (t contentController) Create(c *baa.Context) {
	auth := base.GetUser(c)
	title := c.Query("title")
	content := c.Query("content")
	ret := base.NewReturn()
	params := map[string]interface{}{"title": title, "content": content}
	_, err := model.ContentModel.Create(auth.ID, params)
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	}
	c.JSON(200, ret)
}

// Edit 编辑页面
func (t contentController) Edit(c *baa.Context) {
	cid := c.ParamInt("cid")
	row, err := model.ContentModel.Get(cid)
	if err != nil {
		c.NotFound()
		return
	}

	c.Set("data", row)
	c.HTML(200, "backend/edit")
}

// Update 更新请求
func (t contentController) Update(c *baa.Context) {
	ret := base.NewReturn()
	cid := c.ParamInt("cid")
	title := c.Query("title")
	content := c.Query("content")
	_, err := model.ContentModel.Update(cid, map[string]interface{}{"title": title, "content": content})
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	}

	c.JSON(200, ret)
}

// Delete 删除
func (t contentController) Delete(c *baa.Context) {
	cid := c.ParamInt("cid")
	ret := base.NewReturn()
	err := model.ContentModel.Delete(cid)
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	}

	c.JSON(200, ret)
}
