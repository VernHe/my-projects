package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-programming-tour-book/tag-service/pkg/errcode"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
)

const (
	APP_KEY    = "eddycjy"
	APP_SECRET = "go-programming-tour-book"
)

type API struct {
	URL string
}

type AccessToken struct {
	Token string `json:"token"`
}

// NewAPI 创建API对象
func NewAPI(url string) *API {
	return &API{URL: url}
}

// 发送HTTP Get请求
func (a API) httpGet(ctx context.Context, path string) ([]byte, error) {
	// 没有传入上下文参数，使用的是默认的Background，所以无法实现超时控制
	//resp, err := http.Get(fmt.Sprintf("%s/%s", a.URL, path))
	// 使用下面方式
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, fmt.Sprintf("%s/%s", a.URL, path))

	if err != nil {
		// 将自定义Error转换成gRPC的格式的error
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// 获取Token
func (a *API) getAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?app_key=%s&app_secret=%s", "auth", APP_KEY, APP_SECRET))
	if err != nil {
		return "", err
	}
	var accessToken AccessToken
	// 提取相应的token字段
	_ = json.Unmarshal(body, &accessToken)
	return accessToken.Token, nil
}

// GetTagList 获取标签列表
func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	// 获取token
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTokenFail)
	}

	// 发送getTagList请求
	body, err := a.httpGet(ctx, fmt.Sprintf("api/v1/tags?token=%s&name=%s", token, name))
	if err != nil {
		// 将自定义Error转换成gRPC的格式的error
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}
	return body, nil
}
