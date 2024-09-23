package tvucrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type CryptoMethodRSA struct {
	Name string
	Hash Hash
}

var (
	cryptoMethodPKCS1 *CryptoMethodRSA
)

func init() {
	// PKCS1
	cryptoMethodPKCS1 = &CryptoMethodRSA{Name: "PKCS1", Hash: PKCS1}
	RegisterCryptoMethod(cryptoMethodPKCS1.Alg(), func() CryptoMethod {
		return cryptoMethodPKCS1
	})
}

func (m *CryptoMethodRSA) Alg() string {
	return m.Name
}

func (m *CryptoMethodRSA) CreateSecret(secretNo, bitsNo int) error {
	// 判断secret文件夹是否存在,不存在则创建
	_, err := os.Stat("./secret")
	if os.IsNotExist(err) {
		err := os.Mkdir("./secret", 0777)
		if err != nil {
			return err
		}
	}
	// ------------------RSA生成私钥文件步骤--------------------
	// 1. 生成RSA私钥对
	var bits int
	flag.IntVar(&bits, "key flag", bitsNo, "aaa")
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	// 2. 将私钥公钥对象转换为DER编码形式
	var derPrivateKey, derPublicKey []byte
	switch m.Hash {
	case PKCS1:
		derPrivateKey = x509.MarshalPKCS1PrivateKey(privateKey)
		derPublicKey = x509.MarshalPKCS1PublicKey(publicKey)
	default:
		return errors.New(fmt.Sprintf("[CreateSecret] this hash not exist: %v", m.Hash))
	}
	// 3. 创建私钥pem文件
	privatePath := "./secret/" + m.Alg() + "_PRIVATE_" + strconv.Itoa(secretNo) + ".pem"
	privateFile, err := os.Create(privatePath)
	if err != nil {
		return err
	}
	// 4. 对私钥信息进行编码，写入私钥文件
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPrivateKey,
	}
	err = pem.Encode(privateFile, block)
	if err != nil {
		return err
	}
	// 5. 创建公钥pem文件
	publicPath := "./secret/" + m.Alg() + "_PUBLIC_" + strconv.Itoa(secretNo) + ".pem"
	publicFile, err := os.Create(publicPath)
	if err != nil {
		return err
	}
	// 6. 对公钥信息进行编码，写入公钥文件
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPublicKey,
	}
	err = pem.Encode(publicFile, block)
	if err != nil {
		return err
	}
	return nil
}

func (m *CryptoMethodRSA) GetPrivateKey(secretNo int) (*rsa.PrivateKey, error) {
	privatePath := "./secret/" + m.Alg() + "_PRIVATE_" + strconv.Itoa(secretNo) + ".pem"
	privateFile, err := os.Open(privatePath)
	if err != nil {
		return nil, err
	}
	defer privateFile.Close()
	content, err := ioutil.ReadAll(privateFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	var privateKey *rsa.PrivateKey
	switch m.Hash {
	case PKCS1:
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("[GetPrivateKey] this hash not exist: %v", m.Hash))
	}
	return privateKey, nil
}

func (m *CryptoMethodRSA) GetPublicKey(secretNo int) (*rsa.PublicKey, error) {
	publicPath := "./secret/" + m.Alg() + "_PUBLIC_" + strconv.Itoa(secretNo) + ".pem"
	publicFile, err := os.Open(publicPath)
	if err != nil {
		return nil, err
	}
	defer publicFile.Close()
	content, err := ioutil.ReadAll(publicFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	var publicKey *rsa.PublicKey
	switch m.Hash {
	case PKCS1:
		publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("[GetPublicKey] this hash not exist: %v", m.Hash))
	}
	return publicKey, nil
}

func (m *CryptoMethodRSA) Encrypt(claim Claim, msg []byte) (string, error) {
	// 获取公钥
	publicKey, err := m.GetPublicKey(claim.SecretNo)
	if err != nil {
		return "", err
	}
	// 加密msg
	var msgAfterEncrypt []byte
	msgLength := len(msg)
	step := publicKey.Size() - 11
	switch m.Hash {
	case PKCS1:
		for start := 0; start < msgLength; start += step {
			finish := start + step
			if finish > msgLength {
				finish = msgLength
			}
			chunkMsgAfterEncrypt, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg[start:finish])
			if err != nil {
				return "", err
			}
			msgAfterEncrypt = append(msgAfterEncrypt, chunkMsgAfterEncrypt...)
		}
	default:
		return "", errors.New(fmt.Sprintf("[Encrypt] this hash not exist: %v", m.Hash))

	}
	// base64编码claim
	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}
	head := base64.StdEncoding.EncodeToString(claimBytes)
	// 组合加密内容
	content := head + "." + base64.StdEncoding.EncodeToString([]byte(msgAfterEncrypt))
	return content, nil
}

//func (m *CryptoMethodRSA) Decrypt(content string) (Claim, []byte, error) {
//	contentSlice := strings.Split(content, ".")
//	var claim Claim
//	//claimBytes := []byte(contentSlice[0])
//	msgBytes, _ := base64.StdEncoding.DecodeString(contentSlice[1])
//	claimBytes, _ := base64.StdEncoding.DecodeString(contentSlice[0])
//	err := json.Unmarshal(claimBytes, &claim)
//	if err != nil {
//		return Claim{}, nil, err
//	}
//	privateKey, err := m.GetPrivateKey(claim.SecretNo)
//	if err != nil {
//		return Claim{}, nil, err
//	}
//
//	var msg []byte
//	switch m.Hash {
//	case PKCS1:
//		msg, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, msgBytes)
//		if err != nil {
//			return Claim{}, nil, err
//		}
//	default:
//		return Claim{}, nil, errors.New(fmt.Sprintf("[Encrypt] this hash not exist: %v" , m.Hash))
//	}
//	return claim, msg, nil
//}

func (m *CryptoMethodRSA) Decrypt(content string) (Claim, []byte, error) {
	contentSlice := strings.Split(content, ".")
	var claim Claim
	//claimBytes := []byte(contentSlice[0])
	msgBytes, _ := base64.StdEncoding.DecodeString(contentSlice[1])
	claimBytes, _ := base64.StdEncoding.DecodeString(contentSlice[0])
	err := json.Unmarshal(claimBytes, &claim)
	if err != nil {
		return Claim{}, nil, err
	}
	privateKey, err := m.GetPrivateKey(claim.SecretNo)
	if err != nil {
		return Claim{}, nil, err
	}
	msgLength := len(msgBytes)
	step := privateKey.PublicKey.Size()
	var msg []byte
	switch m.Hash {
	case PKCS1:
		for start := 0; start < msgLength; start += step {
			finish := start + step
			if finish > msgLength {
				finish = msgLength
			}
			chunkMsg, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, msgBytes[start:finish])
			if err != nil {
				return Claim{}, nil, err
			}
			msg = append(msg, chunkMsg...)
		}
	default:
		return Claim{}, nil, errors.New(fmt.Sprintf("[Encrypt] this hash not exist: %v", m.Hash))
	}
	return claim, msg, nil
}
