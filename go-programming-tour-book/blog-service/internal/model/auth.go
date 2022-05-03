package model

import (
	"github.com/jinzhu/gorm"
)

type Auth struct {
	*Model
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (a Auth) TableName() string {
	return "blog_auth"
}

/**
真正的Dao层方法
*/

func (a Auth) Get(db *gorm.DB) (Auth, error) {
	var auth Auth
	// 构建SQL语句
	db = db.Where("app_key = ? AND app_secret = ? AND is_del = ?", a.AppKey, a.AppSecret, 0)
	// 查询一条记录,保存至auth中
	err := db.First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}

	return auth, nil
}
