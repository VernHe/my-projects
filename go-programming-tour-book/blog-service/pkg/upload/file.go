package upload

import (
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const (
	TypeImage FileType = iota + 1
	TypeExcel
	TypeTxt
)

// GetFileName 获取文件名
func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// GetFileExt 获取后缀
func GetFileExt(name string) string {
	return path.Ext(name)
}

// GetSavePath 获取存储的路径
func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

// CheckSavePath 检查存储的路径是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

// CheckContainExt 检查文件后缀是否合法
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

// CheckMaxSize 检查上次文件的大小
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}

// CheckPermission 检查文件权限是否足够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// CreateSavePath 创建要保存的文件
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

// SaveFile 保存文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
