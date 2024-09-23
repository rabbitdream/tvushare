package service

import (
	"context"
	"encoding/base64"
	"os"
	"path"
	"strings"
	"tvushare/common/constants"
	"tvushare/controller/param"
	"tvushare/model"
	"tvushare/pkg/cfg"
	. "tvushare/pkg/log"
	"tvushare/pkg/twitter_share"
	"tvushare/pkg/util"
	"tvushare/pkg/youtube_share"
)

func ShareMedia(ctx context.Context, req *param.ShareMediaReq) param.ErrorInfo {
	res := param.ErrorInfo{
		ErrorCode: 0,
		ErrorMsg:  "success",
	}
	Platform := strings.ToUpper(req.Platform)
	// TIKTOK: 支持URL
	var localDirPath, localVideoPath string
	var videoSize int64
	var err error
	if Platform != constants.TIKTOK {
		localDirPath = path.Join(cfg.Service.TempVideoPath, base64.StdEncoding.EncodeToString([]byte(req.Recipient)))
		Logger.Debugf(ctx, "[ShareMedia] localDirPath: %s", localDirPath)
		if !util.IsExist(localDirPath) {
			err = os.MkdirAll(localDirPath, os.ModePerm)
			if err != nil {
				Logger.Errorf(ctx, "[ShareMedia] mkdir err: %v", err)
				res.ErrorCode = -1
				res.ErrorMsg = "create temp video directory failed"
				return res
			}
		}
		localVideoPath = path.Join(localDirPath, req.MediaUrl[strings.LastIndex(req.MediaUrl, "/"):])
		Logger.Debugf(ctx, "[ShareMedia] localVideoPath: %s", localVideoPath)
		videoSize, err = util.DownloadFile(req.MediaUrl, localVideoPath)
		if err != nil {
			Logger.Errorf(ctx, "[ShareMedia] download http file err: %v", err)
			res.ErrorCode = -1
			res.ErrorMsg = "download http file failed"
			return res
		}
	}
	token, err := model.FindOauthTokenById(req.SocialMediaTokenId)
	if err != nil {
		Logger.Debugf(ctx, "[ShareMedia] FindOauthTokenById err: %v", err)
		res.ErrorCode = -1
		res.ErrorMsg = "find oauth2 token failed"
		return res
	}
	switch Platform {
	case constants.YOUTUBE:
		err = ShareYoutube(ctx, localVideoPath, req.Title, token)
		if err != nil {
			Logger.Errorf(ctx, "[ShareMedia] share media video to youtube_share err: %v", err)
			res.ErrorCode = -1
			res.ErrorMsg = "share media video to youtube failed."
		}
		break
	case constants.TWITTER:
		err = ShareTwitter(ctx, req.Title, localVideoPath, videoSize, token)
		if err != nil {
			Logger.Errorf(ctx, "[ShareMedia] share media video to twitter_share err: %v", err)
			res.ErrorCode = -1
			res.ErrorMsg = "share media video to twitter failed."
		}
		break
	default:
		res.ErrorCode = -1
		res.ErrorMsg = "no support for the platform"
		return res
	}
	if err != nil {
		res.ErrorCode = -1
		res.ErrorMsg = "youtube_share share failed"
	}
	return res
}

func ShareYoutube(ctx context.Context, localVideoPath, title string, token *model.OauthToken) error {
	return youtube_share.UploadVideo(ctx, token, title, localVideoPath)
}

func ShareTwitter(ctx context.Context, title, localVideoPath string, videoSize int64, token *model.OauthToken) error {
	return twitter_share.UploadVideo(ctx, token, title, localVideoPath, videoSize)
}
