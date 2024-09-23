package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"tvushare/controller/param"
	. "tvushare/pkg/log"
	"tvushare/service"
)

func ShareMedia(r *ghttp.Request) {
	ctx := context.Background()
	var req param.ShareMediaReq
	res := param.ErrorInfo{
		ErrorCode: 0,
		ErrorMsg:  "success",
	}
	Logger.Debugf(ctx, "[Export] Body: %v", r.Header.Get("Content-Type"))
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		msg := fmt.Sprintf("[Export] 传参异常: err=%v", err.Error())
		Logger.Error(context.TODO(), msg)
		res.ErrorCode = 500
		res.ErrorMsg = "The source parameter is incorrect. Contact the administrator for confirmation"
		JsonExit(r, res)
	}
	res = service.ShareMedia(ctx, &req)
	JsonExit(r, res)
}
