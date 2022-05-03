package timer

import (
	"time"
)

// GetNowTime 获取当前时间
func GetNowTime() time.Time {
	return time.Now()
}

func GetCalculateTime(currentTimer time.Time, d string) (time.Time, error) {
	// 解析持续时间
	// 其支持的有效单位有"ns", "us" (or "µs"), "ms", "s", "m", "h"
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return currentTimer.Add(duration), nil
}
