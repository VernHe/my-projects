package server

import (
	"context"
	"encoding/json"
	"github.com/go-programming-tour-book/tag-service/pkg/bapi"
	"github.com/go-programming-tour-book/tag-service/pkg/errcode"
	"github.com/go-programming-tour-book/tag-service/proto"
	"google.golang.org/grpc/metadata"
	"log"
)

// TagServer grpc服务器对象
type TagServer struct {
	auth *Auth
}

type Auth struct {
}

func (a *Auth) GetAppKey() string {
	return "go-programming-tour-book"
}

func (a *Auth) GetAppSecret() string {
	return "eddycjy"
}

func (a *Auth) CheckContext(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

func NewTagServer() *TagServer { return &TagServer{} }

func (s *TagServer) GetList(ctx context.Context, r *proto.GetTagListRequest) (*proto.GetTagListReply, error) {
	//panic("测试抛异常")

	if err := s.auth.CheckContext(ctx); err != nil {
		return nil, err
	}

	// 获取metadata
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("incoming md: %+v\n", md)

	// 调用http服务器的接口
	api := bapi.NewAPI("http://127.0.0.1:8888")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, err
	}

	tagList := proto.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}

	return &tagList, nil
}
