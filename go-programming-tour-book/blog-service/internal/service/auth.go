package service

import "errors"

// AuthRequest 封装请求参数
type AuthRequest struct {
	AppKey    string `form:"app_key" binding:"required"`
	AppSecret string `form:"app_secret" binding:"required"`
}

// CheckAuth 检查权限
func (svc *Service) CheckAuth(param *AuthRequest) error {
	auth, err := svc.dao.GetAuth(param.AppKey, param.AppSecret)
	if err != nil {
		return err
	}
	// 正确
	if auth.ID > 0 {
		return nil
	}
	// 不存在
	return errors.New("auth info dose not exist")
}
