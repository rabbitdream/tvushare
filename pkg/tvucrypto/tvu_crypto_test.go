package tvucrypto

import (
	"testing"
)

type M3u8FileJson struct {
	Version             string  `json:"version"`
	FilePath            string  `json:"file_path"`
	BeginTimestamp      int64   `json:"begin_timestamp"`
	EndTimestamp        int64   `json:"end_timestamp"`
	BeginOffset         float64 `json:"begin_offset"`
	ChunkBeginTimestamp int64   `json:"chunk_begin_timestamp"`
	ChunkEndTimestamp   int64   `json:"chunk_end_timestamp"`
}

/*
base64({"alg": "RS256", "secretno": 1}) + '.' + base64(encrypt(content of the original JSON file)).
use private key to encrypt ; use public key to decrypt.
*/

/*
RSA算法
DSA算法
ECC算法
DH算法
*/

/*
1. 指定加密方式，获取该加密方式对应的相关方法
2. 生成该加密方式的秘钥文件（需指定秘钥文件编号和秘钥长度）
3. 公钥加密私钥解密,根据加密规则拼接加密内容与解析密文。
*/

/*
Q:
	1. rsa有密文长度限制。可采取非对称加密加密堆成加密密文。RSA(AES(data))

*/

func TestGenerateRSAKey(t *testing.T) {
	M := GetCryptoMethod("PKCS1")
	err := M.CreateSecret(0, 4096)
	// 也可采用openssl
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("success.")
	}
}

func TestTVUCrypto(t *testing.T) {
	claim := Claim{
		Alg:      "PKCS1",
		SecretNo: 0,
	}
	//m3u8Json := M3u8FileJson{
	//	Version:             "1.0",
	//	FilePath:            "https://mma-s3-test.s3.ap-northeast-2.amazonaws.com/ts/55524C2045585420000001744CBDC6C7/20201010-02-55-39/test_m3u8.m3u8",
	//	BeginTimestamp:      int64(1599015702061),
	//	EndTimestamp:        int64(1599015702061),
	//	BeginOffset:         0.000000,
	//	ChunkBeginTimestamp: int64(1599015702061),
	//	ChunkEndTimestamp:   int64(1599015702061),
	//}
	// 获取加密方法
	M := GetCryptoMethod(claim.Alg)
	//// 序列化
	//dataJson, err := json.Marshal(m3u8Json)
	//if err != nil {
	//	t.Log(err.Error())
	//}
	//// 加密
	//content, err := M.Encrypt(claim, dataJson)
	//if err != nil {
	//	t.Log(err.Error())
	//} else {
	//	t.Logf("success---------------->, content: %s", content)
	//}
	//err = ioutil.WriteFile("./test.json", []byte(content), 0644)
	//if err != nil {
	//	t.Log(err.Error())
	//}

	//s := "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"
	s := "eyJhbGciOiJQS0NTMSIsInNlY3JldF9ubyI6MH0=.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"
	//content := base64.StdEncoding.EncodeToString([]byte(s))
	// 解密
	cm, msg, err := M.Decrypt(s)
	if err != nil {
		t.Log(err.Error())
	}
	t.Logf("decrypt, claim: %v", cm)
	t.Logf("decrypt, msg: %v", string(msg))
}

func TestStringDecrypt(t *testing.T) {
	s := StringDecrypt("VOhqDakjN6B8rvZSClLoZshRNqcBOdg0D+DcogVj1WcQELrHWkjFbKM6J6+eKEw07moa+sv9qsF0Zv9+Je0y8w==", StateEncryptkey)
	t.Log(s)
}
