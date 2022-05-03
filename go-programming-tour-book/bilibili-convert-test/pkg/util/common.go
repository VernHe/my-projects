package util

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
)

// MD5 加密
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

// Get 以GET方法发送HTTP请求
func Get(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("get 出现错误:%s\n", err.Error())
	}
	return resp
}
