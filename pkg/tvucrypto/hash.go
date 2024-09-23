package tvucrypto

type Hash uint

const (
	PKCS1 Hash = 1 + iota
	OAEP
)
