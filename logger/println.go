package logger

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var SourceFiles embed.FS

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

}

func Println(v ...interface{}) {

	// Capture the call stack
	callers := make([]uintptr, 10)
	n := runtime.Callers(2, callers[:])
	frames := runtime.CallersFrames(callers[:n])

	// Log the provided message
	for _, vv := range v {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintf(os.Stderr, "%s:%s:%d %v\n", time.Now().Format(time.RFC3339), file, line, vv)
		} else {
			fmt.Fprintf(os.Stderr, "%s %v\n", time.Now().Format(time.RFC3339), vv)
		}
	}

	// Iterate through the call stack
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// Read the source code of the caller
		sourceLine := getSourceLine(frame.File, frame.Line)

		// Log the details
		fmt.Fprintf(os.Stderr, "%s:%s:%d: %s\n", time.Now().Format(time.RFC3339), frame.File, frame.Line, frame.Function)
		if sourceLine != "" {
			fmt.Fprintf(os.Stderr, "\n%s\n", sourceLine)
		}
	}

}

func getSourceLine(filename string, line int) string {
	// Extract the file name from the full path
	parts := strings.Split(filename, "/")
	shortFile := parts[len(parts)-1]

	var data []byte
	var err error
	// Read the embedded source file
	data, err = SourceFiles.ReadFile(shortFile)
	if err != nil {
		data, err = os.ReadFile(filename)
		if err != nil {
			return ""
		}
	}

	lines := strings.Split(string(data), "\n")
	var sourceLines []string
	for i := line - 6; i <= line+4; i++ {
		if i >= 0 && i < len(lines) {
			lineText := lines[i]
			if i == line-1 {
				lineText = fmt.Sprintf("\033[31m%d: %s\033[0m", i+1, lineText) // Red color
			} else {
				lineText = fmt.Sprintf("%d: %s", i+1, lineText)
			}
			sourceLines = append(sourceLines, lineText)
		}
	}
	return strings.Join(sourceLines, "\n")
}
