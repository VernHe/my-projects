package bilibili

import (
	"github.com/bilibili-convert-test/entity"
	"github.com/bilibili-convert-test/global"
	"github.com/bilibili-convert-test/pkg/util"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// GetPrefix 获取视频url前缀
func GetPrefix(url string) string {
	prefix := url[:strings.Index(url, "video/")+18]
	return prefix
}

// GetVideos 获取视频的所有选集的名称和对应的url
func GetVideos(videoUrl string) *entity.Videos {
	response := util.Get(videoUrl)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("getVideos:%s\n", err)
	}
	strData := string(resp)

	// 获取视频url的前缀
	//urlPrefix := GetPrefix(videoUrl)

	// + 为贪婪，?为非贪婪
	// 匹配课程列表的名称
	reg := regexp.MustCompile("\"part\":\"(.*?)\",\"duration\"")
	itemTitles := reg.FindAllStringSubmatch(strData, -1)
	if len(itemTitles) == 0 {
		panic("请确认是否有视频选集")
	}

	if global.BeginPart > len(itemTitles) || global.BeginPart < 1 {
		log.Fatal("请输入正确的BeginPart")
	}

	itemTitles = itemTitles[global.BeginPart-1:]

	// 所有视频条目
	videos := entity.NewVideos(itemTitles)
	// 用于返回的结果
	//videos := make(map[string]string)
	//if num > 0 {
	//	i := 1
	//	for _, title := range itemTitles {
	//		videos[title[1]] = fmt.Sprintf("%s?p=%d", urlPrefix, i)
	//		fmt.Printf("选集名称：%s, url: %s\n", title[1], videos[title[1]])
	//		i++
	//	}
	//} else {
	//	panic("请确认是否有视频选集")
	//}
	return videos
}
