package service

import (
	"context"
	"tvushare/controller/param"
	"tvushare/model"
	. "tvushare/pkg/log"
)

func GetAuthPlatforms(ctx context.Context, account string, res *param.AuthPlatformsRes) {
	list, err := model.FindAuthListByTvuAccount(ctx, account)
	if err != nil {
		Logger.Errorf(ctx, "[GetAuthPlatforms] FindAuthListByTvuAccount err: %v", err)
		res.ErrorCode = -1
		res.ErrorMsg = "find auth social media accounts failed."
	}
	res.SocialMediaAccountList = list
	return
}
