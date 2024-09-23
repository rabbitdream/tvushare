package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// 允许跨域中间件
func CORS(r *ghttp.Request) {
	opt := ghttp.CORSOptions{
		AllowOrigin:      "*",
		AllowCredentials: "true",
		MaxAge:           63072000,
		AllowHeaders:     "Session",
	}
	r.Response.CORS(opt)
	r.Middleware.Next()
}
