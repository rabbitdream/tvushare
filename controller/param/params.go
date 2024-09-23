package param

import "tvushare/model"

type ErrorInfo struct {
	ErrorCode int    `json:"errorCode" binding:"required"` // 错误码
	ErrorMsg  string `json:"errorMsg" binding:"required"`  // 错误消息
}

type AuthCodeURLRes struct {
	ErrorInfo
	AuthUrl string `json:"authUrl"`
}

type SocialAccountsReq struct {
	ErrorInfo
}

type ShareMediaReq struct {
	Title              string `json:"title"`
	Platform           string `json:"platform"`
	MediaUrl           string `json:"mediaUrl"`
	FileName           string `json:"fileName"`
	ServerName         string `json:"serverName"`
	Recipient          string `json:"recipient"`
	SocialMediaTokenId int    `json:"socialMediaTokenId"`
}

type AuthPlatformsRes struct {
	ErrorInfo
	SocialMediaAccountList []model.TvuSocialMediaAccount `json:"socialMediaAccountList"`
}
