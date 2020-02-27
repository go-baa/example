package backend

import (
	"github.com/go-baa/baa"
	"github.com/go-baa/example/blog/controller/base"
	"github.com/go-baa/example/blog/model"
)

type adminController struct{}

// AdminController 管理员控制器
var AdminController = adminController{}

// Password 修改密码页面
func (t adminController) Password(c *baa.Context) {
	c.HTML(200, "backend/password")
}

// ChangePwd 修改密码
func (t adminController) ChangePwd(c *baa.Context) {
	auth := base.GetUser(c)
	ret := base.NewReturn()
	originPassword := c.QueryTrim("origin_password")
	newPassword := c.QueryTrim("new_password")
	rePassword := c.QueryTrim("re_password")
	if len(originPassword) == 0 || len(newPassword) == 0 {
		ret.Code = 1
		ret.Message = "密码不能为空"
		c.JSON(200, ret)
		return
	}

	if newPassword != rePassword {
		ret.Code = 1
		ret.Message = "新密码两次输入不一致"
		c.JSON(200, ret)
		return
	}
	if originPassword == newPassword {
		ret.Code = 0
		ret.Message = "新旧密码一致，没有发生变化。"
		c.JSON(200, ret)
		return
	}

	err := model.AdminModel.ChangePassword(auth.ID, originPassword, newPassword)
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	}

	c.JSON(200, ret)
}
