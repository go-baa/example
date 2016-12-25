package backend

import (
	"github.com/go-baa/example/blog/controller/base"
	"github.com/go-baa/example/blog/model"
	"net/http"

	"github.com/go-baa/example/blog/modules/log"

	"github.com/baa-middleware/session"
	"gopkg.in/baa.v1"
)

type authController struct{}

// AuthController 验证控制器
var AuthController = authController{}

// Login 登录页面
func (t authController) Login(c *baa.Context) {
	c.HTML(200, "backend/login")
}

// CheckLogin 登录
func (t authController) CheckLogin(c *baa.Context) {
	ret := base.NewReturn()
	username := c.Query("username")
	password := c.Query("password")
	adminInfo, err := model.AdminModel.Login(username, password)
	if err != nil {
		log.Errorf("Login verify error: %v\n", err)
		ret.Code = 1
		ret.Message = "Username or password wrong."
	} else {
		base.SetUser(c, adminInfo)
		ret.Data["auth"] = adminInfo
		c.SetCookie("auth_username", username, 86400*30)
	}

	c.JSON(200, ret)
}

// Logout 退出
func (t authController) Logout(c *baa.Context) {
	s := c.Get("session").(*session.Session)
	s.Delete("auth")
	c.Redirect(http.StatusTemporaryRedirect, "/admin")
	return
}
