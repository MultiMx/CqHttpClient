package botUtil

import (
	"fmt"
	"runtime"
)

func Recover() interface{} {
	if e := recover(); e != nil {
		fmt.Println(e)
		var buf [4096]byte
		fmt.Printf(string(buf[:runtime.Stack(buf[:], false)]))
		return e
	}
	return nil
}
