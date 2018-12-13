package md5

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/axgle/mahonia"
)

func Hash(signString string) string {
	hash := md5.New()
	hash.Write([]byte(signString))
	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func ToUtf8(src string) string {
	srcCoder := mahonia.NewDecoder("GBK")
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder("UTF8")
	_, data, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(data)
}

