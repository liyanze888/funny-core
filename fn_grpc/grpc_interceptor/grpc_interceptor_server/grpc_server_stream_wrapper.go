package grpc_interceptor_server

import (
	"context"
	"github.com/liyanze888/funny-core/fn_utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"sync/atomic"
)

type GrpcServerStreamWrapper struct {
	resource     grpc.ServerStream
	callId       uint64
	streamCallId uint64
}

func (g *GrpcServerStreamWrapper) SetHeader(md metadata.MD) error {
	return g.resource.SetHeader(md)
}

// SendHeader sends the header metadata.
// The provided md and headers set by SetHeader() will be sent.
// It fails if called multiple times.
func (g *GrpcServerStreamWrapper) SendHeader(md metadata.MD) error {
	return g.resource.SendHeader(md)
}

// SetTrailer sets the trailer metadata which will be sent with the RPC status.
// When called more than once, all the provided metadata will be merged.
func (g *GrpcServerStreamWrapper) SetTrailer(md metadata.MD) {
	g.resource.SetTrailer(md)
}

// Context returns the context for this stream.
func (g *GrpcServerStreamWrapper) Context() context.Context {
	return g.resource.Context()
}

// SendMsg sends a message. On error, SendMsg aborts the stream and the
// error is returned directly.
//
// SendMsg blocks until:
//   - There is sufficient flow control to schedule m with the transport, or
//   - The stream is done, or
//   - The stream breaks.
//
// SendMsg does not wait until the message is received by the client. An
// untimely stream closure may result in lost messages.
//
// It is safe to have a goroutine calling SendMsg and another goroutine
// calling RecvMsg on the same stream at the same time, but it is not safe
// to call SendMsg on the same stream in different goroutines.
func (g *GrpcServerStreamWrapper) SendMsg(m interface{}) error {
	return g.resource.SendMsg(m)
}

// RecvMsg blocks until it receives a message into m or the stream is
// done. It returns io.EOF when the client has performed a CloseSend. On
// any non-EOF error, the stream is aborted and the error contains the
// RPC status.
//
// It is safe to have a goroutine calling SendMsg and another goroutine
// calling RecvMsg on the same stream at the same time, but it is not
// safe to call RecvMsg on the same stream in different goroutines.
func (g *GrpcServerStreamWrapper) RecvMsg(m interface{}) error {
	err := g.resource.RecvMsg(m) //阻塞的
	callId := atomic.AddUint64(&g.streamCallId, 1)
	fn_utils.Set("StreamCallId", callId)
	return err
}

func NewGrpcServerStreamWrapper(resource grpc.ServerStream, callId uint64) grpc.ServerStream {
	return &GrpcServerStreamWrapper{
		resource: resource,
		callId:   callId,
	}
}
