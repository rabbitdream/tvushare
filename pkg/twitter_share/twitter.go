package twitter_share

import (
	"context"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/gogf/gf/v2/util/guid"
	"tvushare/model"
	"tvushare/pkg/cfg"
	. "tvushare/pkg/log"
)

const TwitterAuthCodeUrlTemplate string = "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=tweet.read tweet.write users.read offline.access&state=%s&code_challenge=challenge&code_challenge_method=plain"

type RequestCache struct {
	Config        *oauth1.Config `json:"config"`
	RequestSecret string         `json:"request_secret"`
}

var requestCacheMap map[string]RequestCache

func InitAppV1() error {
	requestCacheMap = make(map[string]RequestCache)
	return nil
}

func GetAuthUrlV1(state string) (string, error) {
	id := guid.S()
	twitterConfig := &oauth1.Config{
		ConsumerKey:    cfg.Service.TwitterApiKey,
		ConsumerSecret: cfg.Service.TwitterApiSecret,
		CallbackURL:    fmt.Sprintf("%s?id=%s&state=%s", cfg.Service.RedirectUri, id, state),
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	requestToken, requestSecret, err := twitterConfig.RequestToken()
	if err != nil {
		return "", err
	}
	authorizationURL, err := twitterConfig.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}
	requestCacheMap[requestToken] = RequestCache{
		Config:        twitterConfig,
		RequestSecret: requestSecret,
	}
	return authorizationURL.String(), err
}

func GetTwitterTokenV1(requestToken, verifier string) (*model.OauthToken, error) {
	config := requestCacheMap[requestToken].Config
	requestSecret := requestCacheMap[requestToken].RequestSecret
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	token := &model.OauthToken{
		AccessToken: accessToken,
		TokenSecret: accessSecret,
	}
	return token, nil
}

func GetUserInfoV1(token *model.OauthToken) (*UserInfo, error) {
	config := oauth1.NewConfig(cfg.Service.TwitterApiKey, cfg.Service.TwitterApiSecret)
	tok := oauth1.NewToken(token.AccessToken, token.TokenSecret)
	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, tok)
	return findMyUserV1(httpClient)

}

func UploadVideo(ctx context.Context, token *model.OauthToken, title, videoPath string, videoSize int64) error {
	twitterConfig := oauth1.NewConfig(cfg.Service.TwitterApiKey, cfg.Service.TwitterApiSecret)
	tok := oauth1.NewToken(token.AccessToken, token.TokenSecret)
	httpClient := twitterConfig.Client(oauth1.NoContext, tok)
	ub := &UploadMediaBase{
		client:    httpClient,
		videoSize: videoSize,
		videoPath: videoPath,
	}
	Logger.Debug(ctx, "[UploadVideo] uploadInit start ...")
	err := ub.uploadInit()
	if err != nil {
		Logger.Errorf(ctx, "[UploadVideo] uploadInit err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "[UploadVideo] uploadAppend start, %v...", ub.mediaInit)
	err = ub.uploadAppend()
	if err != nil {
		Logger.Errorf(ctx, "[UploadVideo] uploadAppend err: %v", err)
		return err
	}
	Logger.Debug(ctx, "[UploadVideo] uploadFinalize start...")
	err = ub.uploadFinalize()
	if err != nil {
		Logger.Errorf(ctx, "[UploadVideo] uploadFinalize err: %v", err)
		return err
	}
	Logger.Debug(ctx, "[UploadVideo] checkStatus start...")
	err = ub.checkStatus()
	if err != nil {
		Logger.Errorf(ctx, "[UploadVideo] uploadFinalize err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "[UploadVideo] finish: %v...", ub.mediaFinalize)
	err = ub.tweet(title)
	if err != nil {
		Logger.Errorf(ctx, "[UploadVideo] tweet err: %v", err)
		return err
	}
	Logger.Debug(ctx, "[UploadVideo] success...")
	return nil
}

func GenAuthUrlV2(state string) string {
	return fmt.Sprintf(TwitterAuthCodeUrlTemplate, cfg.Service.TwitterAppClientId, cfg.Service.RedirectUri, state)
}

func GetTwitterTokenV2(code string) (*model.OauthToken, error) {
	return exchangeToken(code)
}

func GetUserInfoV2(access_token string) (*UserInfo, error) {
	return findMyUserV2(access_token)
}
