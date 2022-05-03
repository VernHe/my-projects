package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/util"
	"time"
)

type Claims struct {
	AppKey             string `json:"app_key"`
	AppSecret          string `json:"app_secret"`
	jwt.StandardClaims        // jwt标准字段
}

// GetJWTSecret 获取密钥
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

// GenerateToken 通过 appKey appSecret 生成JWT
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire) // 过期时间
	// 初始化claims
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}

	// 指定加密方法加密claims
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

// ParseToken 解析Token并获取Claims
func ParseToken(token string) (*Claims, error) {
	// 校验并解析claims
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		// 判断是否合法、有效
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
