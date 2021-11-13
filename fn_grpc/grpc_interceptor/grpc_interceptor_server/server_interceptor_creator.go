package grpc_interceptor_server

import (
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var creator = NewServerInterceptorCreator()

type ServerInterceptorCreator struct {
	Unaries grpc.ServerOption
	Streams grpc.ServerOption
}

func NewServerInterceptorCreator() *ServerInterceptorCreator {
	return &ServerInterceptorCreator{
		Unaries: grpc.EmptyServerOption{},
		Streams: grpc.EmptyServerOption{},
	}
}

func CreateInterceptors() (grpc.ServerOption, grpc.ServerOption) {
	return creator.Unaries, creator.Streams
}

func CreateUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	//recoveryHandlerOption := recovery.WithRecoveryHandler(func(p interface{}) (err error) {
	//	debug.PrintStack()
	//	err = fmt.Errorf("panic: %v", p)
	//	return
	//})

	creator.Unaries = grpc.UnaryInterceptor(middleware.ChainUnaryServer(
		//prometheus.UnaryServerInterceptor,
		//recovery.UnaryServerInterceptor(recoveryHandlerOption),
		interceptors...,
	))
}

func CreateStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	//recoveryHandlerOption := recovery.WithRecoveryHandler(func(p interface{}) (err error) {
	//	debug.PrintStack()
	//	err = fmt.Errorf("panic: %v", p)
	//	return
	//})
	//
	creator.Streams = grpc.StreamInterceptor(middleware.ChainStreamServer(
		//prometheus.StreamServerInterceptor,
		//recovery.StreamServerInterceptor(recoveryHandlerOption),
		//StreamLogServerInterceptor(),
		interceptors...,
	))
}
