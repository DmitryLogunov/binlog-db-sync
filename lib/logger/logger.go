package logger

import (
	"encoding/json"
	"fmt"
	"time"
)


func printLog(level string, info string, params interface{}) (n int, err error) {
	currentTime := time.Now().Format("2006/01/02 15:04:05")

	if params == nil {
		return fmt.Printf("[%s] [%s] %s\n", currentTime, level, info)		
	}

	jsonParams, _ := json.Marshal(params)	
	return fmt.Printf("[%s] [%s] %s %s\n", currentTime, level, info, jsonParams)
}

func Info(info string, params interface{}) {
	printLog("info", info, params)
}

func Debug(info string, params interface{}) {
	printLog("debug", info, params)
}

func Warning(info string, params interface{}) {
	printLog("warning", info, params)
}

func Error(info string, params interface{}) {
	printLog("error", info, params)
}