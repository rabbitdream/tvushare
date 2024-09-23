package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"tvushare/common/constants"
	"tvushare/controller/param"
	"tvushare/model"
	. "tvushare/pkg/log"
	"tvushare/pkg/tvucrypto"
	"tvushare/pkg/twitter_share"
	"tvushare/pkg/youtube_share"
)

func SaveSocialAccounts(ctx context.Context, state, code, oauth_token, verifier string, res *param.SocialAccountsReq) {
	// 解析state
	b, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		Logger.Errorf(ctx, "[SaveSocialAccounts] param state failed: %v", err)
		res.ErrorCode = -1
		res.ErrorMsg = fmt.Sprintf("param state failed: %v", err)
		return
	}
	Logger.Debugf(ctx, "[SaveSocialAccounts] state: %v", string(b))
	str := tvucrypto.StringDecrypt(string(b), tvucrypto.StateEncryptkey)
	var stateInfo StateInfo
	err = json.Unmarshal([]byte(str), &stateInfo)
	if err != nil {
		Logger.Errorf(ctx, "[SaveSocialAccounts] param state failed: %v", err)
		res.ErrorCode = -1
		res.ErrorMsg = fmt.Sprintf("param state failed: %v", err)
		return
	}
	Logger.Debugf(ctx, "[SaveSocialAccounts] state: %v, code: %v", stateInfo, code)
	switch stateInfo.Platform {
	case constants.YOUTUBE:
		err = BindYoutube(ctx, stateInfo.TvuAccount, code)
	case constants.TWITTER:
		//err = BindTwitterV2(ctx, stateInfo.TvuAccount, code)
		err = BindTwitter(ctx, stateInfo.TvuAccount, oauth_token, verifier)
	case constants.TIKTOK:
	case constants.FACEBOOK:
	case constants.INSTAGRAM:
	default:
		Logger.Debugf(ctx, "[GetAuthURL] no support for the platform(%s) ...", stateInfo.Platform)
		res.ErrorCode = -1
		res.ErrorMsg = fmt.Sprintf("no support for the platform(%s)", stateInfo.Platform)
		return
	}
}

func BindYoutube(ctx context.Context, tvuAccount, code string) error {
	Logger.Debug(ctx, "[BindYoutube] start ...")
	token, err := youtube_share.GetToken(code)
	if err != nil {
		Logger.Errorf(ctx, "[BindYoutube] get token err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "[BindYoutube] token: %v", token)
	info, err := youtube_share.GetUserInfo(ctx, token)
	if err != nil {
		Logger.Errorf(ctx, "[BindYoutube] get user info err: %v", err)
		return err
	}
	tvuSocialMediaAccount := &model.TvuSocialMediaAccount{
		TvuAccount:           tvuAccount,
		SocialMediaUserName:  info.Snippet.Title,
		SocialMediaAvatarUrl: info.Snippet.Thumbnails.High.Url,
		SocialMediaPlatform:  constants.YOUTUBE,
		SocialMediaUserId:    info.Id,
	}
	return SaveDB(ctx, token, tvuSocialMediaAccount)
}

func BindTwitter(ctx context.Context, tvuAccount, oauth_token, verifier string) error {
	Logger.Debug(ctx, "[BindTwitter] start ...")
	token, err := twitter_share.GetTwitterTokenV1(oauth_token, verifier)
	if err != nil {
		Logger.Errorf(ctx, "[BindTwitter] get token err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "[BindTwitter] token: %v", token)
	info, err := twitter_share.GetUserInfoV1(token)
	if err != nil {
		Logger.Errorf(ctx, "[BindTwitter] get user info err: %v", err)
		return err
	}
	tvuSocialMediaAccount := &model.TvuSocialMediaAccount{
		TvuAccount:           tvuAccount,
		SocialMediaUserName:  info.Data.Name,
		SocialMediaAvatarUrl: info.Data.ProfileImageUrl,
		SocialMediaPlatform:  constants.TWITTER,
		SocialMediaUserId:    info.Data.Id,
	}
	return SaveDB(ctx, token, tvuSocialMediaAccount)
}

func BindTwitterV2(ctx context.Context, tvuAccount, code string) error {
	Logger.Debug(ctx, "[BindTwitter] start ...")
	token, err := twitter_share.GetTwitterTokenV2(code)
	if err != nil {
		Logger.Errorf(ctx, "[BindTwitter] get token err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "[BindTwitter] token: %v", token)
	info, err := twitter_share.GetUserInfoV2(token.AccessToken)
	if err != nil {
		Logger.Errorf(ctx, "[BindTwitter] get user info err: %v", err)
		return err
	}
	tvuSocialMediaAccount := &model.TvuSocialMediaAccount{
		TvuAccount:           tvuAccount,
		SocialMediaUserName:  info.Data.Name,
		SocialMediaAvatarUrl: info.Data.ProfileImageUrl,
		SocialMediaPlatform:  constants.TWITTER,
		SocialMediaUserId:    info.Data.Id,
	}
	return SaveDB(ctx, token, tvuSocialMediaAccount)
}

func SaveDB(ctx context.Context, token *model.OauthToken, tvuSocialMediaAccount *model.TvuSocialMediaAccount) error {
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		Logger.Errorf(ctx, "[BindYoutube] start transaction err: %v", err)
		return err
	}
	tokenId, err := model.SaveOauthToken(token, tx)
	if err != nil {
		Logger.Errorf(ctx, "[BindYoutube] SaveOauthToken err: %v", err)
		tx.Rollback()
		return err
	}
	tvuSocialMediaAccount.SocialMediaTokenId = tokenId
	_, err = model.SaveTvuSocialMediaAccount(tvuSocialMediaAccount, tx)
	if err != nil {
		Logger.Errorf(ctx, "[BindYoutube] SaveTvuSocialMediaAccount err: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
