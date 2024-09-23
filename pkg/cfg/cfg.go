package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"os"
	"path/filepath"
)

var (
	Service struct {
		ProgramPath string
		// 服务端口
		Listen              int
		HTTPSListen         int
		RedirectUri         string
		TwitterAppClientId  string
		TwitterAppAuthBasic string
		TwitterApiKey       string
		TwitterApiSecret    string

		TempVideoPath string
	}

	Redis struct {
		Addr          string
		MaxActive     int
		MaxIdle       int
		IdleTimeout   int
		Wait          bool
		LockPollingMs int64
	}
)

func InitCfg() error {
	var ctx = gctx.New()
	// 读取配置值
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	Service.ProgramPath = filepath.Dir(ex)
	Service.Listen = g.Cfg().MustGet(ctx, "server.listen", 8600).Int()
	Service.HTTPSListen = g.Cfg().MustGet(ctx, "server.https_listen", 10600).Int()
	Service.RedirectUri = g.Cfg().MustGet(ctx, "server.redirect_uri", "http://localhost:8090/social-accounts").String()
	Service.TwitterAppClientId = g.Cfg().MustGet(ctx, "server.twitter_app_client_id", "").String()
	Service.TwitterAppAuthBasic = g.Cfg().MustGet(ctx, "server.twitter_app_auth_basic", "").String()
	Service.TwitterApiKey = g.Cfg().MustGet(ctx, "server.twitter_api_key", "").String()
	Service.TwitterApiSecret = g.Cfg().MustGet(ctx, "server.twitter_api_secret", "").String()

	Redis.Addr = g.Cfg().MustGet(ctx, "redis.addr", "").String()
	Redis.MaxIdle = g.Cfg().MustGet(ctx, "redis.max_idle", 30).Int()
	Redis.MaxActive = g.Cfg().MustGet(ctx, "redis.max_active", 5000).Int()
	Redis.IdleTimeout = g.Cfg().MustGet(ctx, "redis.idle_timeout", 20).Int()
	Redis.Wait = g.Cfg().MustGet(ctx, "redis.wait", true).Bool()
	Redis.LockPollingMs = g.Cfg().MustGet(ctx, "redis.lock_polling_ms", 200).Int64()
	return nil
}
