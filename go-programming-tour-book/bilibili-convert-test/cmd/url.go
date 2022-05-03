package cmd

import (
	"github.com/bilibili-convert-test/entity"
	"github.com/bilibili-convert-test/global"
	"github.com/bilibili-convert-test/pkg/bilibili"
	thirdPart "github.com/bilibili-convert-test/pkg/third-part"
	"github.com/bilibili-convert-test/pkg/validate"
	"github.com/spf13/cobra"
	"log"
)

var urlCmd = &cobra.Command{
	Use:   "download",
	Short: "对应B站视频的url",
	Long:  "对应B站视频的url",
	Run: func(cmd *cobra.Command, args []string) {
		if validate.IsBilibiliUrl(global.InputUrl) {
			log.Println("开始解析Bilibili视频url")
			global.InputUrlPrefix = bilibili.GetPrefix(global.InputUrl)
			videos := bilibili.GetVideos(global.InputUrl)
			for _, videoItem := range videos.Items {
				videoItem.PrintInfo()
			}
			resolveAndDownload(&thirdPart.IiilabInterface{}, videos)
		} else {
			log.Println("Bilibili视频url格式有误，请检查后重试")
		}
	},
}

func init() {
	// download --url/u = "url"
	urlCmd.Flags().StringVarP(&global.InputUrl, "url", "u", "", "下载指定url对应的视频(仅限Bilibili视频)")
	urlCmd.Flags().IntVarP(&global.BeginPart, "beginPart", "b", 1, "起始Part(默认是1，就是从第一Part开始)")
}

func resolveAndDownload(thirdPartInterface thirdPart.ThirdPartInterface, videos *entity.Videos) {
	thirdPartInterface.Download("./", videos)
}
