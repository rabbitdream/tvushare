package controller

import "github.com/gogf/gf/v2/net/ghttp"

func JsonExit(r *ghttp.Request, res interface{}) {
	r.Response.WriteJson(res)
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8")
	r.Exit()
}
