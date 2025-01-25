package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var PrintV = func(v interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s:%s:%d: %+v\n", time.Now().Format(time.RFC3339), filename, line, v)
	fmt.Printf("%s:%s:%d: %+v\n", time.Now().Format(time.RFC3339), filename, line, v)

}

func Println(v ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	if len(v) > 0 {
		for _, vv := range v {
			fmt.Printf("%s:%s:%d: %v\n", time.Now().Format(time.RFC3339), filename, line, vv)
		}
	}
}
