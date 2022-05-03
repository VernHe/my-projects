package third_part

import (
	"encoding/json"
	"fmt"
	"github.com/bilibili-convert-test/entity"
	"github.com/bilibili-convert-test/pkg/fixedThreadPool"
	"github.com/bilibili-convert-test/pkg/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	clientId        = "ee8785f893c95acg"
	clientSecretKey = "160e1c0299db6273a023fbb4800bc4bf"
	serviceUrl      = "https://service.iiilab.com/video/download"
)

type IiilabInterface struct{}

// Resolve 解析视频Bilibili视频Url，得到解析后的url
func Resolve(link string) string {

	timestamp := time.Now().Unix()

	sign := util.MD5(fmt.Sprintf("%s%d%s", link, timestamp, clientSecretKey))

	param := make(map[string]string, 4)
	param["link"] = link
	param["timestamp"] = fmt.Sprintf("%d", timestamp)
	param["sign"] = sign
	param["client"] = clientId

	return SendRequest(serviceUrl, param)
}

func SendRequest(requestPostURL string, params map[string]string) string {
	client := http.Client{}

	// 数据
	urlValues := url.Values{}
	for k, v := range params {
		urlValues.Add(k, v)
	}
	reqBody := urlValues.Encode()
	// 构建请求体
	req, err := http.NewRequest(http.MethodPost, requestPostURL, strings.NewReader(reqBody))
	if err != nil {
		log.Fatal(err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// 读取请求
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result entity.Result
	if err = json.Unmarshal(b, &result); err != nil {
		log.Fatal(err)
	}

	return result.Data.Video
}

// Download 逐个的解析并下载视频
func (thirdPartInterface IiilabInterface) Download(path string, videos *entity.Videos) {
	downloadWithSingleThread(path, videos)
	// 测试线程池
	//downloadWithThreadPool(path, videos)
}

func downloadWithSingleThread(path string, videos *entity.Videos) {
	//i := 0
	for _, videoItem := range videos.Items {
		downloadUrl := Resolve(videoItem.Url)
		util.DownloadFile(fmt.Sprintf("%s%s.mp4", path, videoItem.Title), downloadUrl)
	}
}

func downloadWithThreadPool(path string, videos *entity.Videos) {

	log.Println("开始解析每一个Part的url")
	tasks := make([]*fixedThreadPool.Task, len(videos.Items))
	// 初始化任务
	for index, item := range videos.Items {
		log.Printf("正在解析[ %s ]的url\n", item.Title)
		// 防止出现闭包
		curItem := item
		downloadUrl := Resolve(item.Url)
		tasks[index] = fixedThreadPool.NewTask(
			func(progress *fixedThreadPool.TaskProgress) {
				util.CustomDownloadFile(fmt.Sprintf("%s%s.mp4", path, curItem.Title), downloadUrl, fixedThreadPool.NewThreadPoolWithProgressReader(progress))
			})
	}
	log.Println("解析完成，开始下载")

	// 线程池执行任务
	fixedThreadPool.NewFixedThreadPool(tasks, 3, true).Exec()

}

func TestDownloadWithThreadPool(path string, urls []string) {

	log.Println("开始解析每一个Part的url")
	tasks := make([]*fixedThreadPool.Task, len(urls))
	// 初始化任务
	for index, url := range urls {
		// 防止出现闭包问题
		curInde := index
		tasks[index] = fixedThreadPool.NewTask(
			func(progress *fixedThreadPool.TaskProgress) {
				fmt.Println("开始下载图片 ", curInde)
				util.CustomDownloadFile(fmt.Sprintf("%s图片%d.mp4", path, curInde), url, fixedThreadPool.NewThreadPoolWithProgressReader(progress))
			})
	}

	// 线程池执行任务
	fixedThreadPool.NewFixedThreadPool(tasks, 3, false).Exec()

}

func TestDownloadWithSingleThread(path string, urls []string) {
	//i := 0
	for index, url := range urls {
		i := index
		curUrl := url
		util.DownloadFile(fmt.Sprintf("%s图片%d.jpg", path, i), curUrl)
		//i++
		// TODO 测试代码，只下载一次
		//if i == 3 {
		//	return
		//}
	}
}
