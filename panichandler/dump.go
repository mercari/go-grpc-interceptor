package panichandler

import (
	"fmt"
	"os"
	"runtime/debug"

	"golang.org/x/net/context"
)

var _ PanicHandler = LogPanicDump

// LogPanicDump is a PanicHandler which dumps stack trace.
func LogPanicDump(ctx context.Context, r interface{}) {
	fmt.Fprintf(os.Stderr, string(debug.Stack()))
}
