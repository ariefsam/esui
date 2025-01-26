package logger

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Println(v ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	if len(v) > 0 {
		for _, vv := range v {
			fmt.Printf("%s:%s:%d: %v\n", time.Now().Format(time.RFC3339), filename, line, vv)
		}
	}
}
