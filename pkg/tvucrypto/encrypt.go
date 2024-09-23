package tvucrypto

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	. "tvushare/pkg/log"
)

const StateEncryptkey string = "vfw212u9y8d2fwfl"

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func StringEncrypt(data, key string) string {
	cryptedData, err := AesEncrypt([]byte(data), []byte(key))
	if err != nil {
		Logger.Errorf(context.TODO(), "AesEncrypt failed err: %v , data is: %v, key is: %v ", err, data, key)
		return ""
	}
	cryptedStr := base64.StdEncoding.EncodeToString(cryptedData)
	return cryptedStr
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("input data format error!")
	}
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func StringDecrypt(data, key string) string {
	cryptedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}

	originData, err := AesDecrypt(cryptedData, []byte(key))
	if err != nil {
		return ""
	}
	return string(originData)

}
