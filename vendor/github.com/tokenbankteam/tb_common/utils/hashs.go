package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func HashId(userId int64, max int64) int64 {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strconv.FormatInt(userId, 10)))
	cipherStr := hex.EncodeToString(md5Ctx.Sum(nil))
	val, err := strconv.ParseInt(cipherStr[0:2]+cipherStr[len(cipherStr)-2:], 16, 64)
	if err != nil {
		panic(err)
	}
	return val%max + 1
}
