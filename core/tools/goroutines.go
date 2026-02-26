package tools

import (
	"runtime/debug"
	"turtle/core/lgr"
)

func SafeGoRoutine(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic with stack trace
			lgr.Error("PANIC RECOVERED in goroutine: %v\n", r)
			lgr.Error("Stack trace:\n%s\n", debug.Stack())

			// Optional: Send alert, metric, or notification here
			// alerting.SendPanicAlert(fmt.Sprintf("Goroutine panic: %v", r))
		}
	}()

	fn()
}
