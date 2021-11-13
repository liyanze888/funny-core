package fn_grpc

import (
	"context"
	"fmt"
	"github.com/liyanze888/funny-core/fn_factory"
	"github.com/liyanze888/funny-core/fn_grpc/grpc_interceptor/grpc_interceptor_server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func init() {
	fn_factory.BeanFactory.RegisterBean(NewInitalGrpcWorker())
}

type InitialGrpcWorker struct {
	grpcServer *grpc.Server
}

func (i *InitialGrpcWorker) InitBean(b fn_factory.BeanContext) {
	a, c := grpc_interceptor_server.CreateInterceptors()
	i.CreateGrpcServer(a, c)
}

func (i *InitialGrpcWorker) CreateGrpcServer(options ...grpc.ServerOption) {
	i.grpcServer = grpc.NewServer(
		options...,
	)
	fn_factory.BeanFactory.RegisterBean(i.grpcServer)
}

func (i *InitialGrpcWorker) Start(port int) {
	reflection.Register(i.grpcServer)
	isRunning := true

	httpMux := http.NewServeMux()

	// 健康检查，返回200表示正常运行
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if isRunning {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not running"))
		}
	})

	var wg sync.WaitGroup
	h2Handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			i.grpcServer.ServeHTTP(w, r)
		} else {
			httpMux.ServeHTTP(w, r)
		}
	}), &http2.Server{})
	server := createServer(h2Handler, port)

	// 收到SIGTERM退出流程：
	// 1. 先将isRunning置为false让健康检查失败
	// 2. 等待10秒确保被k8s健康检查调用 & 摘掉此实例 & envoy dns刷新 & 不会有新请求过来
	// 3. 等待现有请求都结束，然后退出
	signalHandler(syscall.SIGTERM, func() {
		isRunning = false
		for i := 10; i > 0; i-- {
			log.Printf("shut down in %d seconds...\n", i)
			time.Sleep(1 * time.Second)
		}
		log.Printf("shutting down...\n")
		_ = server.Shutdown(context.Background())
		//grpcServer.GracefulStop() // Drain() is not implemented
		log.Printf("goodbye.\n")
	})
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	wg.Wait()
}

func NewInitalGrpcWorker() *InitialGrpcWorker {
	return &InitialGrpcWorker{}
}

func createServer(handler http.Handler, port int) *http.Server {
	if port == 0 {
		port = 50053
	}
	if value := os.Getenv("PORT"); value != "" {
		var err error
		port, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
	}
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at %s\n", addr)

	return &http.Server{Addr: addr, Handler: handler}
}

func signalHandler(sig os.Signal, handler func()) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, sig)
	go func() {
		<-sigs
		handler()
	}()
}
