package common

import (
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var Logger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		//AddSource:   true,
		//ReplaceAttr: nil,
	}
	//todo: handle logging os.Stdout to file
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	if Config.IsProduction {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)

	// log to file
	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Logger.Error("Failed to open log file")
	}
	defer f.Close()
	log.SetOutput(f)

	Logger.Info("test")
}

func GetFuncName() string {
	_, file, line, _ := runtime.Caller(1)
	fileStr := strings.Split(file, "/")
	fileName := fileStr[len(fileStr)-1]
	lineNum := strconv.FormatInt(int64(line), 10)

	return fileName + "(" + lineNum + ") - "
}
