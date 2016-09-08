package glog

import (
	"fmt"
	"runtime"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func LogPanicStackMultiLine(ctx context.Context, r interface{}) {
	_, file, line, ok := runtime.Caller(0)
	if ok {
		glog.Errorf("Recovered from panic: %v in %s(%d)", r, file, line)
	}

	callers := []string{}
	for i := 0; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		callers = append(callers, fmt.Sprintf("%d: %s(%d): %s", i, file, line, fn.Name()))
	}
	glog.Warningf("StackTrace:")
	for i := 0; len(callers) > i; i++ {
		glog.Warningf("  %s", callers[i])
	}
}
