package utils

import (
	"runtime"
	"strconv"
	"strings"
)

func GetFuncName() string {
	_, file, line, _ := runtime.Caller(1)
	fileStr := strings.Split(file, "/")
	fileName := fileStr[len(fileStr)-1]
	lineNum := strconv.FormatInt(int64(line), 10)

	return fileName + "(" + lineNum + ")"
}
