package utilaio

import (
	"fmt"
	"runtime/debug"
)

func RecoverOnPanic(f func(), opt ...chan error) {
	defer func() {
		if r := recover(); r != nil {
			if len(opt) > 0 {
				ch := opt[0]
				ch <- fmt.Errorf("retrying function execution, panic [%v] triggered at\n%v", r, string(debug.Stack()))
				RecoverOnPanic(f, ch)
			} else {
				RecoverOnPanic(f)
			}
		}
	}()
	f()
}

func ExecuteSafe(f func(), opt ...chan error) {
	defer func() {
		if r := recover(); r != nil {
			if len(opt) > 0 {
				ch := opt[0]
				ch <- fmt.Errorf("function execution stopped, panic [%v] triggered at\n%v", r, string(debug.Stack()))
			}
		}
	}()
	f()
}

/*
always call HandlePanic as a deferred function above the function which is susceptible to panicking
*/
func HandlePanic(opt ...chan error) {
	if r := recover(); r != nil {
		if len(opt) > 0 {
			ch := opt[0]
			ch <- fmt.Errorf("function execution stopped, panic [%v] triggered at\n%v", r, string(debug.Stack()))
		}
	}
}
