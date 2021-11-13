package fn_exception

import (
	"fmt"
	"runtime"
)

type BizError struct {
	err string
}

func (b BizError) Error() string {
	return b.err
}

func NewBizError(resource error) error {
	return BizError{
		err: getOutPrintErr(resource),
	}
}

func getOutPrintErr(resource error) string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("file = %s:%v \n BizError =%v", file, line, resource)
	}
	return fmt.Sprintf("BizError =%v", resource)
}
