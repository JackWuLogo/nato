package debug

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"micro-libs/utils/color"
	"runtime"
	"strings"
	"time"
)

var (
	units = []string{" bytes", "KB", "MB", "GB", "TB", "PB"}
)

// 性能分析
type Prof struct {
	label string
	log   *logger.Helper
	t     time.Time
}

func (p *Prof) Result() {
	var mstat runtime.MemStats
	runtime.ReadMemStats(&mstat)

	p.log.Debug(color.Question.Text(
		"[Prof]%s: Uptime: %s, Threads: %d, Memory: ( Total: %s, Alloc: %s, Frees: %s ), GC: ( Num: %d, Pause: %d )",
		p.label, time.Since(p.t), runtime.NumGoroutine(),
		formatBytes(mstat.TotalAlloc), formatBytes(mstat.Alloc), formatBytes(mstat.Frees),
		mstat.NumGC, mstat.PauseTotalNs,
	))
}

func NewProf(log *logger.Helper, label ...string) *Prof {
	p := &Prof{
		log: log,
		t:   time.Now(),
	}
	if len(label) > 0 {
		p.label = label[0]
	} else {
		if _, file, line, ok := runtime.Caller(1); ok {
			p.label = fmt.Sprintf("[%s:%d]", logCallerFilePath(file), line)
		}
	}
	return p
}

func formatBytes(val uint64) string {
	var i int
	var target uint64
	for i = range units {
		target = 1 << uint(10*(i+1))
		if val < target {
			break
		}
	}
	if i > 0 {
		return fmt.Sprintf("%0.2f%s (%d bytes)", float64(val)/(float64(target)/1024), units[i], val)
	}
	return fmt.Sprintf("%d bytes", val)
}

// logCallerFilePath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func logCallerFilePath(loggingFilePath string) string {
	// To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	idx := strings.LastIndexByte(loggingFilePath, '/')
	if idx == -1 {
		return loggingFilePath
	}
	idx = strings.LastIndexByte(loggingFilePath[:idx], '/')
	if idx == -1 {
		return loggingFilePath
	}
	return loggingFilePath[idx+1:]
}
