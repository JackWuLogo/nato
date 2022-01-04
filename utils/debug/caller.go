package debug

import (
	"runtime"
)

// 获取调用信息 (单条)
func GetCaller(skip int) runtime.Frame {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip, rpc)
	if n < 1 {
		return runtime.Frame{}
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame
}

// 获取调用信息 (多条)
func GetCallers(skip int, limit int) []runtime.Frame {
	rpc := make([]uintptr, limit)
	n := runtime.Callers(skip, rpc)
	if n < 1 {
		return nil
	}
	frames := runtime.CallersFrames(rpc[:n])
	var lines []runtime.Frame
	for {
		frame, more := frames.Next()
		lines = append(lines, frame)
		if !more {
			break
		}
	}
	return lines
}
