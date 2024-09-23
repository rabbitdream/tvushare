package tvucrypto

import (
	"crypto/rsa"
	"sync"
)

type Claim struct {
	Alg      string `json:"alg"`
	SecretNo int    `json:"secret_no"`
}

var cryptoMethods = map[string]func() CryptoMethod{}
var cryptoMethodLock = new(sync.RWMutex)

type CryptoMethod interface {
	CreateSecret(secretNo, bitsNo int) error
	GetPrivateKey(secretNo int) (*rsa.PrivateKey, error)
	GetPublicKey(secretNo int) (*rsa.PublicKey, error)
	Encrypt(claim Claim, msg []byte) (string, error)
	Decrypt(content string) (Claim, []byte, error)
	Alg() string
}

func RegisterCryptoMethod(alg string, f func() CryptoMethod) {
	cryptoMethodLock.Lock()
	defer cryptoMethodLock.Unlock()
	cryptoMethods[alg] = f
}

func GetCryptoMethod(alg string) CryptoMethod {
	cryptoMethodLock.RLock()
	defer cryptoMethodLock.RUnlock()
	if methodF, ok := cryptoMethods[alg]; ok {
		return methodF()
	}
	return nil
}
