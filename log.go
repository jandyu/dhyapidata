//日志处理
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogLevel string

const (
	LOGLEVEL_DEBUG LogLevel = "DEBUG"
	LOGLEVEL_LOG   LogLevel = "LOG"
	LOGLEVEL_INFO  LogLevel = "INFO"
	LOGLEVEL_WARN  LogLevel = "WARN"
	LOGLEVEL_ERROR LogLevel = "ERROR"
)

//输出日志的级别
var LOG_LEVEL = LogLevel(strings.ToUpper(os.Getenv("LOG_LEVEL")))

//输出日志对应的输出定义，false代表不输出
var LOGLEVEL_FILTER = map[LogLevel]map[LogLevel]bool{
	LOGLEVEL_DEBUG: map[LogLevel]bool{
		LOGLEVEL_ERROR: true, LOGLEVEL_WARN: true, LOGLEVEL_INFO: true, LOGLEVEL_LOG: true, LOGLEVEL_DEBUG: true,
	},
	LOGLEVEL_LOG: map[LogLevel]bool{
		LOGLEVEL_ERROR: true, LOGLEVEL_WARN: true, LOGLEVEL_INFO: true, LOGLEVEL_LOG: true, LOGLEVEL_DEBUG: false,
	},
	LOGLEVEL_INFO: map[LogLevel]bool{
		LOGLEVEL_ERROR: true, LOGLEVEL_WARN: true, LOGLEVEL_INFO: true, LOGLEVEL_LOG: false, LOGLEVEL_DEBUG: false,
	},
	LOGLEVEL_WARN: map[LogLevel]bool{
		LOGLEVEL_ERROR: true, LOGLEVEL_WARN: true, LOGLEVEL_INFO: false, LOGLEVEL_LOG: false, LOGLEVEL_DEBUG: false,
	},
	LOGLEVEL_ERROR: map[LogLevel]bool{
		LOGLEVEL_ERROR: true, LOGLEVEL_WARN: false, LOGLEVEL_INFO: false, LOGLEVEL_LOG: false, LOGLEVEL_DEBUG: false,
	},
}

var LOGLEVEL_LOG_COLOR = map[LogLevel]int{
	LOGLEVEL_DEBUG: 37,
	LOGLEVEL_LOG:   36,
	LOGLEVEL_INFO:  32,
	LOGLEVEL_WARN:  33,
	LOGLEVEL_ERROR: 31,
}

//WriteLevelLog , 写日志
func WriteLevelLog(level LogLevel, txt ...interface{}) {

	if LOGLEVEL_FILTER[LOG_LEVEL][level] {
		//前端显示，背景，前景

		//前端
		//  0  终端默认设置
		//  1  高亮显示
		//  4  使用下划线
		//  5  闪烁
		//  7  反白显示
		//  8  不可见
		// 前景，背景
		//    // 30  40  黑色
		//    // 31  41  红色
		//    // 32  42  绿色
		//    // 33  43  黄色
		//    // 34  44  蓝色
		//    // 35  45  紫红色
		//    // 36  46  青蓝色
		//    // 37  47  白色
		//f := LOGLEVEL_LOG_COLOR[level]
		//fmt.Printf("\n%c[%d;%d;%dm[%s]%s%s",0x1B,1,40,f,level,time.Now().Local().Format("2006-01-02 15:04:05"), "::")
		//fmt.Printf(time.Now().Local().Format("2006-01-02 15:04:05"))

		s := time.Now().Local().Format("2006-01-02 15:04:05") + " " + fmt.Sprintln(txt...)
		fmt.Print(s)
		f := time.Now().Local().Format("20060102") + ".log"
		//fmt.Printf("%c[0m",0x1B)

		fileObj, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to open the file", err.Error())

		}
		if _, err := io.WriteString(fileObj, s); err != nil {
			fmt.Println("error write log", err.Error())
		}

		fileObj.Close()
	}
}

//InfoLog ,记录info日志，一般为系统基本信息或状态【包::func::,info,info1...】
func InfoLog(txt ...interface{}) {
	WriteLevelLog(LOGLEVEL_INFO, txt...)
}

//LLog ,记录log日志，一般为系统流程或调用信息【包::func::,info,info1...】
func LLog(txt ...interface{}) {
	WriteLevelLog(LOGLEVEL_LOG, txt...)
}

//ErrorLog ,记录error日志，一般为系统错误【包::func::,info,info1...】
func ErrorLog(txt ...interface{}) {
	WriteLevelLog(LOGLEVEL_ERROR, txt...)
}

//WarnLog ,记录Wran日志，一般为系统提示警告【包::func::,info,info1...】
func WarnLog(txt ...interface{}) {
	WriteLevelLog(LOGLEVEL_WARN, txt...)
}

//DebugLog ,记录debug日志，一般为开发测试用【包::func::,info,info1...】
func DebugLog(txt ...interface{}) {
	WriteLevelLog(LOGLEVEL_DEBUG, txt...)
}

func init() {
	//输出相关环境变量
	if LOG_LEVEL == "" {
		//LOG_LEVEL = LOGLEVEL_WARN
		LOG_LEVEL = LOGLEVEL_DEBUG
	}
	fmt.Println("env log init:", LOG_LEVEL)
}
