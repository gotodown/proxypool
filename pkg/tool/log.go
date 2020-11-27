package tool

import (
	"io"
	"log"
	"os"
	"time"
)

func LogToFile() *log.Logger {

	file := "logs/" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	defer logFile.Close()
	if err != nil {
		panic(err)
	}
	// 创建一个Logger: 参数1： 日志写入的目的地，参数2： 自定义前缀，参数3：日志属性
	return log.New(logFile, "proxypool--", log.Lshortfile)
}

func LogToFileAndStdout() *log.Logger {
	file := "./" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	// defer logFile.Close()
	writers := []io.Writer{
		logFile,
		os.Stdout}

	fileAndStdoutWrite := io.MultiWriter(writers...)
	return log.New(fileAndStdoutWrite, "[proxypool]:##", log.Ldate|log.Ltime|log.Lshortfile)
}

var Logger *log.Logger

func init() {
	Logger = LogToFileAndStdout()
}
