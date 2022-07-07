package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logging    *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func NewLogger() error {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = openLogFile(fileName, filePath)
	if err != nil {
		return err
	}

	logging = log.New(F, DefaultPrefix, log.LstdFlags)
	return nil
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	log.Println("[DEBUG]", v)
	logging.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	log.Println("[INFO]", v)
	logging.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	log.Println("[WARN]", v)
	logging.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	log.Println("[ERROR]", v)
	logging.Println(v)
}

func Errorf(format string, v ...interface{}) {
	setPrefix(ERROR)
	log.Println("[ERROR]", v)
	logging.Printf(format, v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	log.Println("[FATAL]", v)
	logging.Fatalln(v)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logging.SetPrefix(logPrefix)
}
