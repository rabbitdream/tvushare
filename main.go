package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"tvushare/pkg/cfg"
	. "tvushare/pkg/log"
	"tvushare/pkg/twitter_share"
	"tvushare/pkg/youtube_share"
	"tvushare/router"
)

func main() {
	ctx := context.Background()
	fmt.Println("tvushare ...")
	_ = cfg.InitCfg()
	InitLogger()
	err := youtube_share.InitApp(ctx)
	if err != nil {
		Logger.Errorf(ctx, "youtube init app err: %v", err)
		return
	}
	err = twitter_share.InitAppV1()
	if err != nil {
		Logger.Errorf(ctx, "twitter init app err: %v", err)
		return
	}
	runServer()
}

func runServer() {
	s := g.Server()
	s.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
		403: func(r *ghttp.Request) {
			g.Log().Debug(context.TODO(), "403")
			r.Response.Writeln("403")
		},
		404: func(r *ghttp.Request) {
			g.Log().Debug(context.TODO(), "404")
			r.Response.Writeln("404")
		},
		500: func(r *ghttp.Request) {
			g.Log().Debug(context.TODO(), "500")
			r.Response.Writeln("500")
		},
	})
	// 绑定路由
	router.InitRouter(s)
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS11,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
	s.EnableHTTPS("server.crt", "server.key", tlsConfig)
	s.SetHTTPSPort(cfg.Service.HTTPSListen)
	s.SetPort(cfg.Service.Listen)
	s.Run()
}
