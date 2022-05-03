package validate

import (
	"log"
	"regexp"
)

var bilibiliUrlPattern = "^((http|https)://)?www.bilibili.com/video/(.{12})\\?.+"

// IsBilibiliUrl 判断是否为bilibili视频url
func IsBilibiliUrl(url string) bool {
	// 参数校验
	if len(url) == 0 {
		return false
	}
	matched, err := regexp.MatchString(bilibiliUrlPattern, url)
	if err != nil {
		log.Fatal(err)
	}
	return matched
}
