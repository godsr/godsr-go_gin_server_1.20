package service

import (
	"crypto/sha256"
	"encoding/hex"
	"github/godsr/go_gin_server/util"
)

func HashSALT(password string) string {
	// 비밀번호 암호화
	hash := sha256.New()
	hashValue := password + util.Conf("HASH_SALT") //소금
	hash.Write([]byte(hashValue))
	md := hash.Sum(nil)
	hashPw := hex.EncodeToString(md)

	return hashPw
}
