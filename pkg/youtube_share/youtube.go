package youtube_share

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"os"
	"tvushare/model"
	. "tvushare/pkg/log"
)

var YoutubeConfig *oauth2.Config

func InitApp(ctx context.Context) error {
	b, err := ioutil.ReadFile("./config/youtube_client_secret.json")
	if err != nil {
		Logger.Errorf(ctx, "Unable to read youtube_share client secret file: %v", err)
		return err
	}
	// If modifying the scope, delete your previously saved credentials
	// at ~/.credentials/youtube_share-go.json
	YoutubeConfig, err = google.ConfigFromJSON(b, youtube.YoutubeScope, youtube.YoutubeReadonlyScope, youtube.YoutubeUploadScope) // 读取用户信息，上传视频权限
	//config, err := google.ConfigFromJSON([]byte("AIzaSyAuk5vkU2qCuw2MT_q_8qRqf3wrcy0eMts"), scope)
	if err != nil {
		Logger.Errorf(ctx, "Unable to parse youtube client secret file to config: %v", err)
		return err
	}
	return nil
}

func GenAuthUrl(state string) string {
	return YoutubeConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func GetToken(code string) (*model.OauthToken, error) {
	return exchangeTokenV2(YoutubeConfig, code)
}

func GetUserInfo(ctx context.Context, token *model.OauthToken) (*youtube.Channel, error) {
	tok := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
	}
	client := YoutubeConfig.Client(ctx, tok)
	service, err := youtube.New(client)
	if err != nil {
		Logger.Errorf(ctx, "[GetUserInfo] creating youtube client err: %v", err)
		return nil, err
	}
	call := service.Channels.List([]string{"snippet,contentDetails,statistics"})
	call = call.Mine(true)
	res, err := call.Do()
	if err != nil {
		Logger.Errorf(ctx, "[GetUserInfo] get user info err: %v", err)
		return nil, err
	}
	Logger.Debugf(ctx, "[GetUserInfo] res: %v", res.Items)
	return res.Items[0], nil
}

func UploadVideo(ctx context.Context, token *model.OauthToken, title, videoPath string) error {
	tok := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
	}
	client := YoutubeConfig.Client(ctx, tok)
	service, err := youtube.New(client)
	if err != nil {
		Logger.Errorf(ctx, "creating YouTube client err: %v", err)
		return err
	}
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title: title,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "public"},
	}
	call := service.Videos.Insert([]string{"snippet,status"}, upload)
	file, err := os.Open(videoPath)
	if err != nil {
		Logger.Errorf(ctx, "open local temp file: %s, err: %v", videoPath, err)
		return err
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		Logger.Errorf(ctx, "youtube upload video err: %v", err)
		return err
	}
	Logger.Debugf(ctx, "Youtube upload successful! Video ID: %v\n", response.Id)
	return nil
}
