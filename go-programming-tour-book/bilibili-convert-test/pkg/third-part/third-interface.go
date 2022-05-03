package third_part

import "github.com/bilibili-convert-test/entity"

type ThirdPartInterface interface {
	Download(path string, videos *entity.Videos)
}
