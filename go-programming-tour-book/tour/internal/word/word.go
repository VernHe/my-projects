package word

import (
	"strings"
	"unicode"
)

/**
单词全部转为小写。
单词全部转为大写。
下划线单词转为大写驼峰。	go_home GoHome
下划线单词转为小写驼峰。	go_home goHome
驼峰转为小写驼峰。
*/

// ToUpper 全部变大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower 全部变小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// UnderscoreToUpperCamelCase 下划线转大写驼峰
func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)   // 将 _ 替换 空格，-1表示全部替换
	s = strings.Title(s)                   // 将单词第一个字母变成大写
	return strings.Replace(s, " ", "", -1) // 将空格去除
}

// UnderscoreToLowerCamelCase 下划线转小写驼峰
func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// CamelCaseToUnderscore 驼峰转下划线
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		// 第一个字符转换成小写	（防止出现第一个字母大写从而出现一开始就是下划线的情况）
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		// 遇见大写字母加下划线
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		// 字母通通通小写
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
