package fn_utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type (
	//ProxyHandler 代理 Handler run执行handler    faultTole 是否容错  timeout超时时间
	ProxyHandler func() (run RunHandler, faultTole bool, timeout int64)
	//RunHandler  执行函数
	RunHandler func() error
)

func MutilRunV2(handlers ...ProxyHandler) error {
	//执行的数量
	runNum := len(handlers)
	// 不容错的数量
	unfaultToleNum := 0
	runNumChn := make(chan int)
	unfaultToleNumChn := make(chan int)
	errorChn := make(chan error)
	tout := int64(0)
	defer func() {
		close(runNumChn)
		close(unfaultToleNumChn)
		close(errorChn)
	}()

	for _, handler := range handlers {
		runHandler, faultTole, timeout := handler()
		if !faultTole {
			unfaultToleNum++
		}
		if tout < timeout {
			tout = timeout
		}
		go func(run RunHandler, faultTole bool) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("recover  %s\n", err)
				}
			}()
			err := run()
			if err == nil && !faultTole {
				//不容错的执行
				unfaultToleNumChn <- 1
			}
			if err != nil && !faultTole {
				errorChn <- err
			}
			runNumChn <- 1
		}(runHandler, faultTole)
	}
	timeoutChn := time.After(time.Duration(tout) * time.Millisecond)
	needFinishSign := 0
	finishSign := 0
	for {
		select {
		case key := <-runNumChn:
			log.Printf("finishChn %d", key)
			finishSign++
			if finishSign == runNum {
				return nil
			}
		case key := <-unfaultToleNumChn:
			log.Printf("unfaultToleNumChn %d", key)
			needFinishSign++
		case err := <-errorChn:
			log.Printf("errorChn")
			return err
		case _ = <-timeoutChn:
			if needFinishSign != unfaultToleNum {
				return errors.New("timeout")
			}
			return nil
		}
	}
}

// ------------------------ 下面是例子
type (
	Test1Param struct {
		name string
	}

	Test1Resp struct {
		name string
	}
	Test2Param struct {
		age int
	}

	Test2Resp struct {
		age string
	}
)

func getTest1ProxyHandler(param Test1Param, resp **Test1Resp) ProxyHandler {
	return func() (run RunHandler, faultTole bool, timeout int64) {
		return func() error {
			return test1(param, resp)
		}, true, int64(1000 * time.Millisecond)
	}
}

func getTest2ProxyHandler(param Test2Param, resp **Test2Resp) ProxyHandler {
	return func() (run RunHandler, faultTole bool, timeout int64) {
		return func() error {
			return test2(param, resp)
		}, false, int64(1000 * time.Millisecond)
	}
}

func getTest3ProxyHandler(resp *bool) ProxyHandler {
	return func() (run RunHandler, faultTole bool, timeout int64) {
		return func() error {
			return test3(resp)
		}, false, int64(1000 * time.Millisecond)
	}
}

func test1(param Test1Param, resp **Test1Resp) error {
	ctx, canFun := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer canFun()
	chn := make(chan int)
	go func() {
		//time.Sleep(1000 * time.Millisecond)
		chn <- 1
	}()
	select {
	case _ = <-chn:
		*resp = &Test1Resp{
			fmt.Sprintf("result = %s", param.name),
		}
		return nil
	case <-ctx.Done():
		return errors.New("test1 timeout")
	}
}

func test2(param Test2Param, resp **Test2Resp) error {
	ctx, canFun := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer canFun()
	chn := make(chan int)
	go func() {
		//time.Sleep(1000 * time.Millisecond)
		chn <- 1
	}()
	select {
	case _ = <-chn:
		*resp = &Test2Resp{
			fmt.Sprintf("result = %d", param.age),
		}
		return nil
	case <-ctx.Done():
		return errors.New("test2 timeout")
	}
}

func test3(resp *bool) error {
	*resp = true
	return nil
}

func ConcurrentV2Test() {
	var resp1 *Test1Resp
	test1ProxyHandler := getTest1ProxyHandler(Test1Param{
		name: "this is test1",
	}, &resp1)

	var resp2 *Test2Resp
	test2ProxyHandler := getTest2ProxyHandler(Test2Param{
		age: 20,
	}, &resp2)

	var resp3 bool
	test3ProxyHandler := getTest3ProxyHandler(&resp3)

	err := MutilRunV2(test1ProxyHandler, test2ProxyHandler, test3ProxyHandler)
	if err != nil {
		log.Printf("ConcurrentV2Test err %v", err)
		return
	}

	log.Printf("test1 = %v", resp1)
	log.Printf("test2 = %v", resp2)
	log.Printf("test3 = %v", resp3)
}
