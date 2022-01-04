package tool

import (
	"fmt"
	"micro-libs/utils/log"
)

// 安全协程
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Error("[SafeGo] Error: %+v", err)
			}
		}()

		fn()
	}()
}
