package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	test1()
	test2()
	test3()
	test4()
}

func test1() {
	str1 := "hello,world"
	for i := 0; i < len(str1); i++ {
		fmt.Printf("[%x]", str1[i])
	}
	fmt.Println()
	str2 := "hi,世界"
	for i := 0; i < len(str2); i++ {
		fmt.Printf("[%x]", str2[i])
	}
	fmt.Println()
}

func test2() {
	// 使用UTF-8编码，英文占2个字节，中文占3个字节
	str1 := "GoGOGO"
	str2 := "Go语言"
	fmt.Println(len(str1)) // 6
	fmt.Println(len(str2)) // 8

	for index, c := range str2 {
		fmt.Printf("%#U starts at byte position %d\n", c, index)
	}
	// 输出:
	// U+0047 'G' position is 0
	// U+006F 'o' position is 1
	// U+8BED '语' position is 2
	// U+8A00 '言' position is 5

	for i, w := 0, 0; i < len(str2); i += w {
		runeVal, width := utf8.DecodeRuneInString(str2[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeVal, i)
		w = width
	}
}

func test3() {
	// 1个字节
	var rune1 rune = 'a'
	str1 := "a"
	// 2个字节
	var rune2 rune = 'ä'
	str2 := "ä"
	fmt.Printf("rune1: Unicode %#U, size: %d, value: %d\n", rune1, len(str1), rune1)
	fmt.Printf("rune2: Unicode %#U, size: %d, value: %d\n", rune2, len(str2), rune2)
	// 输出:
	// rune1: Unicode U+0061 'a', size: 1
	// rune2: Unicode U+00E4 'ä', size: 2
}

func test4() {
	rowStr := `\123\t\321`
	normalStr := "\t\123"
	fmt.Println(rowStr)    // 输出: \123\t\321
	fmt.Println(normalStr) // 输出:          S
}

func test5() {
	str1 := "hello " + "world" // 在编译时完成
	str2 := str1 + "a"         // 在运行时完成
	fmt.Println(str2)
}
