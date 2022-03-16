package botUtil

import (
	"fmt"
	"github.com/MultiMx/CqHttpClient/state"
	"runtime"
)

func Recover() {
	if p := recover(); p != nil {
		fmt.Println(p)
		var buf [4096]byte
		fmt.Printf(string(buf[:runtime.Stack(buf[:], false)]))
		state.ErrPanicCounter++
	}
}
