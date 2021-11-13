package grpc_interceptor_server

import (
	"context"
	"github.com/liyanze888/funny-core/fn_exception"
	"github.com/liyanze888/funny-core/fn_log"
	"github.com/liyanze888/funny-core/fn_utils"
	"google.golang.org/grpc"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var callId uint64

func StreamLogServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		callId := atomic.AddUint64(&callId, 1)
		start := time.Now()
		fn_log.Printf("callId %d  callPath%v stream\n", callId, info.FullMethod)
		defer func() {
			fn_log.Printf("StreamServerInterceptor defer")
			fn_utils.Clear()
		}()
		defer func() {
			fn_log.Printf("call %d finished %v", callId, time.Since(start))
			if err := recover(); err != nil {
				fn_log.Printf("call %d %v stream\n  error = %v  ,  stack =%v", callId, info.FullMethod, err, string(debug.Stack()))
			}
		}()
		defer func() {
			if err := recover(); err != nil {
				fn_log.Printf("this is      %v", err)
			}
		}()
		fn_utils.Set("CallId", callId)
		return handler(srv, NewGrpcServerStreamWrapper(ss, callId))
	}
}

func UnaryLogServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("GoRoutineId = %v", fn_utils.GetGoRoutineId())
		callId := atomic.AddUint64(&callId, 1)
		start := time.Now()
		fn_log.Printf("call %d %v {%v}\n", callId, info.FullMethod, req)
		fn_utils.Set("CallId", callId)
		defer func() {
			fn_log.Printf("UnaryServerInterceptor defer")
			fn_utils.Clear()
		}()
		defer func() {
			fn_log.Printf("call %d %v finished %v", callId, info.FullMethod, time.Since(start))
			if err := recover(); err != nil {
				fn_log.Printf("error call %d %v stream\n  error = %v  ,  stack =%v", callId, info.FullMethod, err, string(debug.Stack()))
			}
		}()
		defer func() {
			//捕捉的是panic中传递的数据
			if err := recover(); err != nil {
				//log.Printf("%v", reflect.TypeOf(err).Elem().Name())
				if tErr, ok := err.(fn_exception.BizError); ok {
					fn_log.Printf("this is  tErr    %v", tErr)
				} else {
					fn_log.Printf("this is  other    %v", tErr)
				}
			}
		}()
		return handler(ctx, req)
	}
}
