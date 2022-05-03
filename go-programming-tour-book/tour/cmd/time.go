package cmd

import (
	"fmt"
	"github.com/go-programming-tour-book/tour/tour/internal/timer"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"time"
)

var calculateTime string // 要计算的时间
var duration string

// timeCmd time子命令
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

// time命令下的now子命令,示例: time now
var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		// 当前时间 + 自 1970 年1 月1 日UTC以来经过的秒数
		log.Printf("当前时间: %s,时间戳: %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

//
var calculateTimeCmd = &cobra.Command{
	Use:   "calc",
	Short: "计算所需时间",
	Long:  "计算所需时间",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			// 如果没有传入要计算的时间，默认使用当前时间
			currentTimer = timer.GetNowTime()
		} else {
			// 传入了时间
			var err error
			spaceNum := strings.Count(calculateTime, " ") // 空格数
			if spaceNum == 0 {
				layout = "2006-01-02"
			}
			if spaceNum == 1 {
				layout = "2006-01-02 15:04:05"
			}
			currentTimer, err = time.Parse(layout, calculateTime)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}
		targetTime, err := timer.GetCalculateTime(currentTimer, duration)
		if err != nil {
			log.Fatalf("timer.GetCalculateTime err: %v", err)
		}

		log.Printf("输出结果: %s, %d", targetTime.Format(layout), targetTime.Unix())
	},
}

func init() {
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calculateTimeCmd)

	// time calc [--calculate/-c] [--duration/-d]
	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", ` 需要计算的时间，有效单位为时间戳或已格式化后的时间 `)
	calculateTimeCmd.Flags().StringVarP(&duration, "duration", "d", "", ` 持续时间，有效时间单位为"ns", "us" (or "µ s"), "ms", "s", "m", "h"`)
}
