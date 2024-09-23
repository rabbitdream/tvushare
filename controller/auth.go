package controller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
	"tvushare/common/constants"
	"tvushare/controller/param"
	. "tvushare/pkg/log"
	"tvushare/service"
)

func GetAuthURL(r *ghttp.Request) {
	ctx := context.Background()
	res := param.AuthCodeURLRes{
		ErrorInfo: param.ErrorInfo{
			ErrorCode: 0,
			ErrorMsg:  "success",
		},
	}
	platform := r.GetQuery("platform").String()
	tvuAccount := r.GetQuery("account").String()
	Logger.Debugf(ctx, "[GetAuthURL] platform: %s, tvuAccount: %s", platform, tvuAccount)
	var err error
	switch strings.ToUpper(platform) {
	case constants.YOUTUBE:
		Logger.Debugf(ctx, "[GetAuthURL] YOUTUBE ...")
		res.AuthUrl, err = service.GetYoutubeAuthUrl(ctx, tvuAccount)
		break
	case constants.TWITTER:
		Logger.Debugf(ctx, "[GetAuthURL] TWITTER ...")
		res.AuthUrl, err = service.GetTwitterAuthUrl(ctx, tvuAccount)
		break
	case constants.TIKTOK:
		Logger.Debugf(ctx, "[GetAuthURL] TIKTOK ...")
		break
	case constants.FACEBOOK:
		Logger.Debugf(ctx, "[GetAuthURL] FACEBOOK ...")
		break
	case constants.INSTAGRAM:
		Logger.Debugf(ctx, "[GetAuthURL] INSTAGRAM ...")
		break
	default:
		Logger.Debugf(ctx, "[GetAuthURL] no support for the platform(%s) ...", platform)
		res.ErrorCode = -1
		res.ErrorMsg = "no support for the platform"
		JsonExit(r, res)
	}
	if err != nil {
		Logger.Errorf(ctx, "[GetAuthURL] get url err: %v", err)
		res.ErrorCode = -1
		res.ErrorMsg = fmt.Sprintf("get auth code url err: %v", err)
	}
	Logger.Debugf(ctx, "[GetAuthURL] auth code url: %v", res.AuthUrl)
	JsonExit(r, res)
}

func SaveSocialAccounts(r *ghttp.Request) {
	ctx := context.Background()
	res := param.SocialAccountsReq{}
	res.ErrorCode = 0
	res.ErrorMsg = "success"
	Logger.Debugf(ctx, "[SaveSocialAccounts] url:%v", r.URL.String())
	state := r.GetQuery("state").String()
	code := r.GetQuery("code").String()
	oauth_token := r.GetQuery("oauth_token").String()
	verifier := r.GetQuery("oauth_verifier").String()
	Logger.Debugf(ctx, "[SaveSocialAccounts] state:%v, code: %v, oauth_token: %v, verifier: %v", state, code, oauth_token, verifier)
	service.SaveSocialAccounts(ctx, state, code, oauth_token, verifier, &res)
	JsonExit(r, res)
}

func AuthPlatforms(r *ghttp.Request) {
	ctx := context.Background()
	var res param.AuthPlatformsRes
	res.ErrorCode = 0
	res.ErrorMsg = "success"
	Logger.Debugf(ctx, "[AuthPlatforms] url:%v", r.URL.String())
	account := r.GetQuery("account").String()
	service.GetAuthPlatforms(ctx, account, &res)
	JsonExit(r, res)
}
