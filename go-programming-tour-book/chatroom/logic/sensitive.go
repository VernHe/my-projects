package logic

import (
	"github.com/goprogramming-tour-book/chatroom/global"
	"strings"
)

//FilterSensitive 过滤敏感词汇，将其替换成*
func FilterSensitive(msg string) string {
	for _, sensitiveWord := range global.SensitiveWords {
		msg = strings.ReplaceAll(msg, sensitiveWord, "**")
	}
	return msg
}
