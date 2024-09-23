package router

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"tvushare/controller"
	"tvushare/middleware"
)

func InitRouter(s *ghttp.Server) {
	s.Group("/", func(root *ghttp.RouterGroup) {
		root.Middleware(middleware.CORS)
		root.Bind(
			new(controller.Hello),
		)
		root.GET("/authurl", controller.GetAuthURL)
		root.GET("/auth-platforms", controller.AuthPlatforms)
		root.GET("/social-accounts", controller.SaveSocialAccounts)
		root.POST("/share-media", controller.ShareMedia)
	})
}
