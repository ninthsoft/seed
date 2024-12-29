package autoprocs

import (
	"io"
	"log"
	"os"
	"runtime"
)

const _EnvNameGoMaxProcs = "GOMAXPROCS"

// AutoSet 自动设置 GOMAXPROCS 的值
func AutoSet(logWriter io.Writer) {
	var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	var last = runtime.GOMAXPROCS(-1)
	defer func() {
		var current = runtime.GOMAXPROCS(-1)
		logger.Printf("[INFO] automaxprocs current GOMAXPROCS=%v, last GOMAXPROCS=%v\n", current, last)
	}()
}

// Set 设置 GOMAXPROCS 值
// 若环境变量 GOMAXPROCS 有值将不进行设置（go自己会使用该值）
// cpuNum: 之前的 GOMAXPROC S值
// ok: 若环境变量 GOMAXPROCS 有值或者num<1，返回 false，其他情况为 true
func Set(num int) (cpuNum int, ok bool) {
	if os.Getenv(_EnvNameGoMaxProcs) != "" {
		return runtime.GOMAXPROCS(-1), false
	}
	return runtime.GOMAXPROCS(num), num >= 1
}

func init() {
	AutoSet(os.Stderr)
}
