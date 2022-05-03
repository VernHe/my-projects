package global

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var (
	once    = new(sync.Once)
	RootDir string
)

func Init() {
	once.Do(func() {
		inferRootDir()
		log.Printf("项目根目录: %v\n", RootDir)
		initConfig()
		log.Println("配置文件读取成功")
	})
}

func inferRootDir() {
	// 返回对应于当前目录的根路径名
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		// 确保项目根目录下存在template目录
		if exists(d + "/template") {
			return d
		}

		// 递归的向上查找，知道找到存在 /template 的路径，然后返回
		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
