package logger

import (
	"fmt"
	"os"
	"time"
	"tpcs/global"
	"tpcs/pkg/file"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", global.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		global.AppSetting.LogSaveName,
		time.Now().Format(global.AppSetting.TimeFormat),
		global.AppSetting.LogFileExt,
	)
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	return nil, fmt.Errorf("os.Getwd err: %v\n", err)
	//}
	//src := currentDir + "/" + filePath

	src := filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s\n", src)
	}

	err := file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v\n", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile :%v\n", err)
	}

	return f, nil
}
