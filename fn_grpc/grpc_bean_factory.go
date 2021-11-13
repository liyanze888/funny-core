package fn_grpc

import (
	"github.com/liyanze888/funny-core/fn_factory"
	"github.com/liyanze888/funny-core/fn_grpc/fn_grpc_config"
	"github.com/liyanze888/funny-core/fn_log"
	"log"
	"reflect"
)

func init() {
	fn_log.Printf("GrpcBeanContext")
}

var GrpcBeanFactory = newGrpcBeanFactory()

func init() {
	fn_factory.BeanFactory.RegisterBean(GrpcBeanFactory)
}

type GrpcBeanContext interface {
	Register(registerMethod interface{})
	StartUp()
}

type grpcBeanFactory struct {
	grpcBeans  map[string]interface{}
	GrpcWorker *InitialGrpcWorker           `autowire:""`
	Config     *fn_grpc_config.FnGrpcConfig `autowire:""`
}

func (g *grpcBeanFactory) StartUp() {
	for _, method := range g.grpcBeans {
		mType := reflect.TypeOf(method)
		mValue := reflect.ValueOf(method)
		numParam := mType.NumIn()
		params := make([]reflect.Value, 0, numParam)
		for i := 0; i < numParam; i++ {
			log.Printf("%v", mType.In(i))
			log.Printf("%v", fn_factory.BeanFactory.FindBeanDefinitionsByType(mType.In(i)))
			params = append(params, fn_factory.BeanFactory.FindBeanDefinitionsByType(mType.In(i))[0].Value)
		}
		mValue.Call(params)
	}
	g.GrpcWorker.Start(g.Config.Port)
}

func (g *grpcBeanFactory) Register(registerMethod interface{}) {
	name := reflect.TypeOf(registerMethod).Name()
	g.grpcBeans[name] = registerMethod
}

func newGrpcBeanFactory() GrpcBeanContext {
	return &grpcBeanFactory{
		grpcBeans: make(map[string]interface{}),
	}
}
