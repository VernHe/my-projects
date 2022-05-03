package entity

import (
	"fmt"
	"github.com/bilibili-convert-test/global"
)

// VideoItem Bilibili视频选集的条目
type VideoItem struct {
	Title string
	Url   string
}

type Videos struct {
	Items []*VideoItem
}

// NewVideos 创建Videos
func NewVideos(items [][]string) *Videos {
	videos := Videos{
		Items: make([]*VideoItem, len(items)),
	}
	videos.Init(items)
	return &videos
}

func (v *Videos) Init(items [][]string) {
	if len(items) > 0 {
		num := global.BeginPart
		for index, title := range items {
			v.Items[index] = &VideoItem{
				Title: title[1],
				Url:   fmt.Sprintf("%s?p=%d", global.InputUrlPrefix, num),
			}
			num++
		}
	}
}

// PrintInfo 打印条目信息
func (item *VideoItem) PrintInfo() {
	fmt.Printf("选集标题: %s, 选集链接: %s\n", item.Title, item.Url)
}
