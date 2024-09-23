package util

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, localFilePath string) (int64, error) {
	file, err := os.OpenFile(localFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()
	rsp, err := http.Get(url)
	defer func() {
		_ = rsp.Body.Close()
	}()
	return io.Copy(file, rsp.Body)
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else {
			return false
		}
	}
	return true
}
