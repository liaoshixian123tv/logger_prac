package main

import (
	"errors"
	"os"
	"time"

	go_logger "github.com/phachon/go-logger"
)

var logPath string
var globalLogger *go_logger.Logger

func init() {
	var setLoggerErr error
	logPath = "leoFile"
	globalLogger, setLoggerErr = SetLogger(logPath)
	if setLoggerErr != nil {
		panic(setLoggerErr)
	}
}

func newLogFile(path string) {
	if _, err := os.Stat("./" + logPath); os.IsNotExist(err) {
		if err := os.Mkdir("./"+logPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

// SetLogger 設置logger
func SetLogger(path string) (logger *go_logger.Logger, err error) {

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	newLogFile(path)
	logger = go_logger.NewLogger()

	detachErr := logger.Detach("console")
	if detachErr != nil {

		panic(detachErr)
	}
	consoleConfig := &go_logger.ConsoleConfig{
		Color:  true,
		Format: "%timestamp_format% %body%",
	}
	attachConsoleErr := logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)
	if attachConsoleErr != nil {
		panic(attachConsoleErr)
	}
	fileConfig := &go_logger.FileConfig{
		Filename:   "./" + logPath + "/server.log",
		MaxSize:    10 * 1024, // File maximum (KB), default 0 is not limited
		MaxLine:    100,       // The maximum number of lines in the file, the default 0 is not limited
		DateSlice:  "d",       // Cut the document by date, support "Y" (year), "m" (month), "d" (day), "H" (hour), default "no"
		JsonFormat: false,     // Whether the file data is written to JSON formatting
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	attachFileErr := logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	if attachFileErr != nil {
		panic(attachFileErr)
	}
	logger.SetAsync()
	return
}

func main() {
	for i := 0; i < 100; i++ {
		globalLogger.Infof("info info info %d info info info", 100)
		globalLogger.Criticalf("Critical Critical Critical %s Critical Critical", "123456")
		time.Sleep(time.Millisecond * 200)
	}
}
