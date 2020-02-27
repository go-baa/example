package router

import (
	"net/http"

	"github.com/baa-middleware/session"
	"github.com/go-baa/baa"
	"github.com/go-baa/example/blog/controller/backend"
	"github.com/go-baa/example/blog/controller/frontend"
	"github.com/go-baa/example/blog/model"
)

// Router 路由
func Router(b *baa.Baa) {
	b.SetAutoHead(true)
	b.SetAutoTrailingSlash(true)

	// 前台路由
	b.Get("/", frontend.IndexController.Index)
	b.Get("/p/:page", frontend.ContentController.Search)
	b.Get("/c/:id", frontend.ContentController.Show)
	b.Get("/about", frontend.IndexController.About)

	// 后台路由
	b.Get("/admin/login", backend.AuthController.Login).Name("auth_login")
	b.Get("/admin/logout", backend.AuthController.Logout).Name("auth_logout")
	b.Post("/admin/login", backend.AuthController.CheckLogin)
	b.Group("/admin", func() {
		b.Get("/password", backend.AdminController.Password).Name("change_pwd")
		b.Post("/password", backend.AdminController.ChangePwd)
		b.Get("/", backend.IndexController.Index).Name("content_index")

		// content
		b.Group("/content", func() {
			b.Get("/add", backend.ContentController.Add).Name("content_add")
			b.Post("/", backend.ContentController.Create).Name("content_create")
			b.Group("/:cid", func() {
				b.Get("/edit", backend.ContentController.Edit)
				b.Put("/", backend.ContentController.Update)
				b.Delete("/", backend.ContentController.Delete)
			})
			b.Get("/search", backend.ContentController.Search)
		})

		// upload
		b.Post("/upload", backend.CommonController.Upload)
	}, loginAuth)

	// ping
	b.Get("/ping", func(c *baa.Context) {
		c.Text(200, []byte("PONG"))
	})

	// error
	b.SetNotFound(func(c *baa.Context) {
		c.Set("code", 404)
		c.HTML(404, "base/error/404")
	})
}

func loginAuth(c *baa.Context) {
	s := c.Get("session").(*session.Session)
	rawAuth := s.Get("auth")
	if rawAuth == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		return
	}

	auth := rawAuth.(*model.AdminInfo)
	c.Set("auth", auth)
}
