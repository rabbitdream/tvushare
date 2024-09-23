package log

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

var Logger *glog.Logger

func InitLogger() {
	Logger = g.Log().Line()
}
