package grpc_interceptor_client

import (
	"context"
	"google.golang.org/grpc"
)

type ClientInterceptorCreator struct {
	ClientLogInterceptor *ClientLogInterceptor `autowire:""`
}

func NewClientInterceptorCreator() *ClientInterceptorCreator {
	return &ClientInterceptorCreator{}
}

func (i *ClientInterceptorCreator) CrearteClientInterceptors() (grpc.DialOption, grpc.DialOption) {
	return i.createUnaryInterceptors(), i.createStreamInterceptors()
}

func (i *ClientInterceptorCreator) createUnaryInterceptors() grpc.DialOption {
	client := ChainUnaryClient(
		i.ClientLogInterceptor.UnaryClientInterceptor(),
	)
	return grpc.WithUnaryInterceptor(client)
}

func (i *ClientInterceptorCreator) createStreamInterceptors() grpc.DialOption {
	client := ChainStreamInterceptors(
		i.ClientLogInterceptor.ClientStreamerInterceptor(),
	)
	return grpc.WithStreamInterceptor(client)
}

func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//主要伪装Invoker
		chainer := func(currentInter grpc.UnaryClientInterceptor, currentInvoker grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(currentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return currentInter(currentCtx, method, req, reply, cc, currentInvoker, opts...)
			}
		}
		chainInvoker := invoker
		for i := n - 1; i >= 0; i-- {
			//每次invoker 不一样
			chainInvoker = chainer(interceptors[i], chainInvoker)
		}
		return chainInvoker(ctx, method, req, reply, cc, opts...)
	}
}

func ChainStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	n := len(interceptors)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		chainer := func(currentInter grpc.StreamClientInterceptor, currentStreamer grpc.Streamer) grpc.Streamer {
			return func(currentCtx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return currentInter(currentCtx, desc, cc, method, currentStreamer, opts...)
			}
		}

		chainedHandler := streamer
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, desc, cc, method, opts...)
	}
}
