package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"tvushare/common/constants"
	. "tvushare/pkg/log"
	"tvushare/pkg/tvucrypto"
	"tvushare/pkg/twitter_share"
	"tvushare/pkg/youtube_share"
)

type StateInfo struct {
	Platform   string `json:"platform"`
	TvuAccount string `json:"tvu_account"`
}

func GetYoutubeAuthUrl(ctx context.Context, tvuAccount string) (string, error) {
	stage := StateInfo{
		Platform:   constants.YOUTUBE,
		TvuAccount: tvuAccount,
	}
	b, err := json.Marshal(&stage)
	if err != nil {
		Logger.Debugf(ctx, "[GetYoutubeAuthUrl] json marshal err: %v", err)
		return "", err
	}
	return youtube_share.GenAuthUrl(tvucrypto.StringEncrypt(string(b), tvucrypto.StateEncryptkey)), nil
}

func GetTwitterAuthUrl(ctx context.Context, tvuAccount string) (string, error) {
	stage := StateInfo{
		Platform:   constants.TWITTER,
		TvuAccount: tvuAccount,
	}
	b, err := json.Marshal(&stage)
	if err != nil {
		Logger.Debugf(ctx, "[GetYoutubeAuthUrl] json marshal err: %v", err)
		return "", err
	}
	return twitter_share.GetAuthUrlV1(base64.StdEncoding.EncodeToString([]byte(tvucrypto.StringEncrypt(string(b), tvucrypto.StateEncryptkey))))
}
