package twitter_share

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"tvushare/model"
	"tvushare/pkg/cfg"
)

const (
	OAUTH_TOKEN_URL    string = "https://api.x.com/2/oauth2/token"
	MEDIA_ENDPOINT_URL string = "https://upload.twitter.com/1.1/media/upload.json"
	POST_TWEET_URL     string = "https://api.x.com/2/tweets"
	MinSize            int64  = 5 * 1024 * 1024
)

type UserInfoData struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	UserName        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type UserInfo struct {
	Data UserInfoData `json:"data"`
}

type MediaUploadInitRes struct {
	MediaId          int64  `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	ExpiresAfterSecs int64  `json:"expires_after_secs"`
	MediaKey         string `json:"media_key"`
}

type MediaUploadProcessingInfo struct {
	State          string `json:"state"`
	CheckAfterSecs int    `json:"check_after_secs"`
}

type MediaUploadFinalizeRes struct {
	MediaId          int64                      `json:"media_id"`
	MediaIdString    string                     `json:"media_id_string"`
	ExpiresAfterSecs int64                      `json:"expires_after_secs"`
	Size             int64                      `json:"size"`
	ProcessingInfo   *MediaUploadProcessingInfo `json:"processing_info"`
}

type UploadData struct {
	Command       string `json:"command"`
	TotalBytes    int    `json:"total_bytes"`
	MediaType     string `json:"media_type"`
	MediaCategory string `json:"media_category"`
	//AdditionalOwners string `json:"additional_owners"`
}

type UploadMediaBase struct {
	client        *http.Client
	videoSize     int64
	videoPath     string
	mediaInit     *MediaUploadInitRes
	mediaFinalize *MediaUploadFinalizeRes
}

func findMyUserV1(client *http.Client) (*UserInfo, error) {
	res, err := client.Get("https://api.twitter.com/2/users/me?user.fields=id,location,profile_image_url,username&expansions=affiliation.user_id")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	var info UserInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &info, nil
}

// POST https://upload.twitter.com/1.1/media/upload.json?command=INIT&total_bytes=10240&media_type=image/jpeg
func (ub *UploadMediaBase) uploadInit() error {
	ps := fmt.Sprintf("command=INIT&total_bytes=%v&media_type=%v&additional_owners=%v", ub.videoSize, "video%2Fmp4", "1826887606740942849")
	fmt.Println(ps)
	payload := strings.NewReader(ps)
	res, err := ub.client.Post(MEDIA_ENDPOINT_URL, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	var mediaInit MediaUploadInitRes
	err = json.Unmarshal(body, &mediaInit)
	if err != nil {
		return err
	}
	ub.mediaInit = &mediaInit
	return nil
}

func (ub *UploadMediaBase) uploadAppend() error {
	// Uploads media in chunks and appends to chunks uploaded
	segmentId := 0
	var bytesSent int64 = 0
	file, err := os.Open(ub.videoPath)
	if err != nil {
		return err
	}
	defer file.Close()
	//reader := bufio.NewReader(file)
	var errR error
	var n int
	for bytesSent < ub.videoSize {
		chunk := make([]byte, MinSize)
		n, errR = file.Read(chunk)
		if errR != nil || n == 0 {
			break
		}
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("command", "APPEND")
		_ = writer.WriteField("media_id", ub.mediaInit.MediaIdString)
		_ = writer.WriteField("segment_index", strconv.Itoa(segmentId))
		part, err := writer.CreateFormFile("media", fmt.Sprintf("chunk_%d", segmentId))
		if err != nil {
			return err
		}
		if _, err = io.Copy(part, bytes.NewReader(chunk[:n])); err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(writer.FormDataContentType())
		res, err := ub.client.Post(MEDIA_ENDPOINT_URL, writer.FormDataContentType(), payload)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		if err != nil || res.StatusCode < http.StatusOK || res.StatusCode > 299 {
			return fmt.Errorf("upload chunk(%v) failed, status: %v, err: %v", segmentId, res.StatusCode, err)
		}
		segmentId += 1
		bytesSent += int64(n)
	}
	if (errR != nil && err != io.EOF) || bytesSent != ub.videoSize {
		return fmt.Errorf("upload chunks failed; err: %v, bytesSent: %v", errR, bytesSent)
	}
	return nil
}

func (ub *UploadMediaBase) uploadFinalize() error {
	// Finalizes uploads and starts video processing
	payload := strings.NewReader(fmt.Sprintf("command=FINALIZE&media_id=%v", ub.mediaInit.MediaId))
	res, err := ub.client.Post(MEDIA_ENDPOINT_URL, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	var mediaFinalize MediaUploadFinalizeRes
	err = json.Unmarshal(body, &mediaFinalize)
	if err != nil {
		return err
	}
	ub.mediaFinalize = &mediaFinalize
	return nil
}

func (ub *UploadMediaBase) checkStatus() error {
	// Checks video processing status
	if ub.mediaFinalize.ProcessingInfo == nil {
		return nil
	}

	if ub.mediaFinalize.ProcessingInfo.State == "succeeded" {
		return nil
	}
	if ub.mediaFinalize.ProcessingInfo.State == "failed" {
		return fmt.Errorf("upload media state: %v", ub.mediaFinalize.ProcessingInfo.State)
	}
	time.Sleep(time.Duration(ub.mediaFinalize.ProcessingInfo.CheckAfterSecs) * time.Second)
	url := fmt.Sprintf("%s?command=STATUS&media_id=%v", MEDIA_ENDPOINT_URL, ub.mediaInit.MediaId)
	res, err := ub.client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	var temp MediaUploadFinalizeRes
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return err
	}
	ub.mediaFinalize.ProcessingInfo = temp.ProcessingInfo
	return ub.checkStatus()
}

func (ub *UploadMediaBase) tweet(title string) error {
	params := map[string]interface{}{
		"text": title,
		"media": map[string]interface{}{
			"media_ids": []string{ub.mediaInit.MediaIdString},
		},
	}
	paramsData, err := json.Marshal(params)

	fmt.Println(string(paramsData))
	res, err := ub.client.Post(POST_TWEET_URL, "application/json", bytes.NewBuffer(paramsData))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("tweet abnormal, status: %v", res.StatusCode)
	}

	return nil
}

func exchangeToken(code string) (*model.OauthToken, error) {

	payload := strings.NewReader(fmt.Sprintf("code=%s&grant_type=authorization_code&client_id=%s&redirect_uri=%s&code_verifier=challenge", code, cfg.Service.TwitterAppClientId, cfg.Service.RedirectUri))

	client := &http.Client{}
	req, err := http.NewRequest("POST", OAUTH_TOKEN_URL, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", cfg.Service.TwitterAppAuthBasic))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	var token model.OauthToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &token, nil
}

func findMyUserV2(access_token string) (*UserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/me?user.fields=id,location,profile_image_url,username&expansions=affiliation.user_id", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", access_token))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	var info UserInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &info, nil
}
