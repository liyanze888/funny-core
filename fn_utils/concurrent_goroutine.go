package fn_utils

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type (
	Run              func(BaseRoutineParam) (interface{}, error)
	BaseRoutineParam struct {
		Timeout int64
		Param   interface{}
	}
	RoutineContext struct {
		Timeout int64
		// 是否容错
		Tolerate      map[string]bool
		NeedFinishNum int
		Params        map[string]BaseRoutineParam
		Handlers      map[string]Run
	}

	dataPayload struct {
		key  string
		data interface{}
	}

	RoutineResult struct {
		// 税局
		Datas map[string]interface{}
		// 错误
		Err error
	}
)

func MutilRun(ctx RoutineContext) RoutineResult {
	untolerateChn := make(chan string)
	finishChn := make(chan string)
	dataChn := make(chan dataPayload)
	errorChn := make(chan error)
	defer close(finishChn)
	defer close(errorChn)
	total := len(ctx.Handlers)
	result := RoutineResult{
		Datas: make(map[string]interface{}),
	}
	for key, handler := range ctx.Handlers {
		go func(key string, run Run) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("recover %s %s\n", key, err)
				}
			}()
			param := ctx.Params[key]
			param.Timeout = ctx.Timeout - 10
			if resp, err := run(param); err != nil {
				if !ctx.Tolerate[key] {
					errorChn <- err
				}
			} else {
				dataChn <- dataPayload{
					key, resp,
				}
				if !ctx.Tolerate[key] {
					untolerateChn <- key
				}
			}
			finishChn <- key
		}(key, handler)
	}

	timeoutChn := time.After(time.Duration(ctx.Timeout) * time.Millisecond)
	needFinishSign := 0
	finishSign := 0
	for {
		select {
		case key := <-finishChn:
			log.Printf("finishChn %s", key)
			finishSign++
			if finishSign == total {
				return result
			}
		case key := <-untolerateChn:
			log.Printf("tolerateChn %s", key)
			needFinishSign++
		case payload := <-dataChn:
			result.Datas[payload.key] = payload.data
		case err := <-errorChn:
			log.Printf("errorChn")
			result.Err = err
			return result
		case _ = <-timeoutChn:
			if needFinishSign != ctx.NeedFinishNum {
				result.Err = errors.New("timeout")
			}
			return result
		}
	}
}
