package twitter_share

import (
	"context"
	"os"
	"testing"
	"tvushare/pkg/log"
)

func Test_UploadInit(t *testing.T) {
	log.InitLogger()
	fi, _ := os.Stat("CROWN.mp4")
	err := UploadVideo(context.Background(), nil, "first tweet video", "CROWN.mp4", fi.Size())
	t.Log(err)
	return
}
