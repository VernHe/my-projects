package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type DownloadFileReader interface {
	Read(p []byte) (int, error)
	Init(resp *http.Response, file *os.File)
}

type ShowProgressReader struct {
	io.Reader
	FileName        string
	TotalSize       int64
	CurrentSize     int64
	CurrentProgress int
}

func (reader *ShowProgressReader) Read(p []byte) (int, error) {
	n, err := reader.Reader.Read(p)
	if err == io.EOF {
		log.Println("下载完成")
		return n, err
	}
	if err != nil {
		log.Fatal(err)
	}
	reader.updateProgress(n)

	//log.Printf("当前下载进度: [%f]", )
	return n, err
}

func (reader *ShowProgressReader) updateProgress(n int) {
	reader.CurrentSize += int64(n)
	reader.CurrentProgress = int(reader.CurrentSize * 100 / reader.TotalSize)
	var processBar string
	if reader.CurrentProgress < 100 {
		processBar = fmt.Sprintf("[%s%s%s]", strings.Repeat("=", reader.CurrentProgress>>1), ">", strings.Repeat(" ", 49-reader.CurrentProgress>>1))
		fmt.Printf("\r%s: %s %d%% ", reader.FileName, processBar, reader.CurrentProgress)
	} else {
		processBar = fmt.Sprintf("[%s]", strings.Repeat("=", reader.CurrentProgress>>1))
		fmt.Printf("\r%s: %s %d%% ", reader.FileName, processBar, reader.CurrentProgress)
	}
}

func (s *ShowProgressReader) Init(resp *http.Response, file *os.File) {
	s.Reader = resp.Body
	s.TotalSize = resp.ContentLength
	// TODO 待测试
	s.FileName = file.Name()
}

func NewShowProgressReader(r io.Reader, totalSize int64, fileName string) *ShowProgressReader {
	return &ShowProgressReader{
		Reader:    r,
		TotalSize: totalSize,
		FileName:  fileName,
	}
}

func DownloadFile(filepath string, url string) {
	log.Printf("正在准备下载文件到: [%s]\n", filepath)
	response := Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Fatalf("downloadFile 出现错误:%s\n", "链接错误")
	}

	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
		}
	}(file)

	reader := NewShowProgressReader(response.Body, response.ContentLength, filepath[2:])

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
	}
}

func CustomDownloadFile(filepath string, url string, reader DownloadFileReader) {

	response := Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Fatalf("downloadFile 出现错误,状态码:%s\n", response.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
		}
	}(file)

	reader.Init(response, file)

	//reader := NewShowProgressReader(response.Body, response.ContentLength, filepath[2:])

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatalf("downloadFile 出现错误:%s\n", err.Error())
	}
}
