package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	_mrand "math/rand"
	"strings"
	"time"
)

//时间
func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//时间
func GetDate() string {
	return time.Now().Format("2006-01-02")
}

//保存到数据库时用的加密
func mD5Value(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//保存到数据库时用的加密
func EnCode(str string) string {
	str = mD5Value(str)
	if len(str) > 8 {
		str = str[:3] + "c" + str[3:]
		str = str[:5] + "b" + str[5:]
	}
	return str
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func GetRandUser() string {
	user := GetGuid()
	if len(user) > 15 {
		return user[:15]
	}
	return user
}

func RandGet(quality string, info []map[string]string) map[string]string {
	if quality != "" {
		if strings.Index(quality, "A") < 0 && strings.Index(quality, "B") < 0 && strings.Index(quality, "C") < 0 && strings.Index(quality, "D") < 0 {
			quality = ""
		}
	}
	_mrand.Seed(time.Now().UnixNano() / 100)
	rst := make(map[string]string)
	for {
		index := _mrand.Intn(len(info))
		rst = info[index]
		if quality == "" || strings.Index(quality, rst["quality"]) >= 0 {
			break
		}
	}
	return rst
}
