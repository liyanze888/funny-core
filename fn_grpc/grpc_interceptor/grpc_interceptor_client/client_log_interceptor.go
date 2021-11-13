package grpc_interceptor_client

import (
	"context"
	"github.com/liyanze888/funny-core/fn_log"
	"google.golang.org/grpc"
	"log"
)

type ClientLogInterceptor struct {
}

func NewClientLogInterceptor() *ClientLogInterceptor {
	return &ClientLogInterceptor{}
}

func (c *ClientLogInterceptor) UnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fn_log.Printf("method = %v  req =%v , reply = %v", method, req, reply)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func (c *ClientLogInterceptor) ClientStreamerInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Printf("method = %v  desc = %v", method, *desc)
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return clientStream, err
	}
}
