package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
	"time"
)

// AccessLog 日志拦截器
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestLog := "access request log: method: %s, begin_time: %d, request: %v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err = handler(ctx, req)

	responseLog := "access response log: method: %s, begin_time: %d, end_time: %d, response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return
}

// ErrorLog 错误日志
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	// 对错误进行处理
	if err != nil {
		errLog := "error log: method: %s, code: %v, message: %v, details: %v"
		s, _ := status.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return
}

// Recpvery 异常处理拦截器
func Recpvery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			recoveryLog := "recovery log: method: %s, message: %v, stack: %s"
			log.Printf(recoveryLog, info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()
	return handler(ctx, req)
}
