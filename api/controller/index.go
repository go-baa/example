package controller

import (
	"github.com/go-baa/example/api/model"
	"github.com/go-baa/log"
	"gopkg.in/baa.v1"
)

type index struct{}

// IndexController ...
var IndexController = index{}

// Index list articles
func (index) Index(c *baa.Context) {
	page := c.ParamInt("page")
	pagesize := 10

	rows, total, err := model.ArticleModel.Search(page, pagesize)
	if err != nil {
		output(c, 1, err.Error(), nil)
		return
	}

	log.Debugf("rows: %#v, total: %d\n", rows, total)

	output(c, 0, "", map[string]interface{}{
		"total": total,
		"items": rows,
	})
}

// Show show article
func (index) Show(c *baa.Context) {
	id := c.ParamInt("id")

	row, err := model.ArticleModel.Get(id)
	if err != nil {
		output(c, 1, err.Error(), nil)
		return
	}

	log.Debugf("row: %#v\n", row)

	output(c, 0, "", row)
}

// NormalReturn json output struct
type NormalReturn struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func output(c *baa.Context, code int, message string, data interface{}) {
	ret := new(NormalReturn)
	ret.Code = code
	ret.Message = message
	if data != nil {
		ret.Data = data
	}
	c.JSON(200, ret)
}
