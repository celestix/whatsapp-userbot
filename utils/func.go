package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(v any) string {
	str := strings.Split(
		runtime.FuncForPC(
			reflect.ValueOf(v).Pointer(),
		).Name(),
		".",
	)
	return str[len(str)-1]
}
