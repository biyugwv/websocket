package app

import(
	//"log"
	"time"
	"os"
	"sync"
	"fmt"
	"strconv"
)

var allLog map[string] *appLog
const(
	INFO string = "info"        // 消息
	WARNING string = "warning"  // 警告
	ERROR string = "error"      // 错误
	MESSAGE string = "message"   //  websocket通讯记录
)


type appLog struct{
	updateTime int64  //  更新时间
	file *os.File
	logType  string       // 日志类型
	logBuffer []byte      // 日志写入buffer
	copyBuffer []byte
	lockBuffer *sync.Mutex
	lockFile  *sync.RWMutex
	maxSize int64
	maxTime int64
}


func getMonthString(month string) string {
    var months string
    switch month {
		case "January":
		    months = "01"
		case "February":
		    months = "02"
		case "March":
		    months = "03"
		case "April":
		    months = "04"
		case "May":
		    months = "05"
		case "June":
		    months = "06"
		case "July":
		    months = "07"
		case "August":
		    months = "08"
		case "September":
		    months = "09"
		case "October":
		    months = "10"
		case "November":
		    months = "11"
		case "December":
		    months = "12"
	}
	return months
}



func LogInit(){
	allLog = make(map[string]*appLog)
	var logs = []string{INFO,WARNING,ERROR,MESSAGE}

	for _,v := range logs {
		l := &appLog{}
		l.logType = v
		l.updateTime = time.Now().Unix()
		l.maxSize = 1 * 1024 * 1024
		l.maxTime =  10
		// 创建文件
		year:=strconv.Itoa(time.Now().Year())
		month:= getMonthString(time.Now().Month().String())
		
		day:=strconv.Itoa(time.Now().Day())
		filePath := "../websocket/logs/"+year+month+day+"/"
		filename := v+".log"
		f ,err := os.Open(filePath + filename)
		if err != nil {
			err = os.MkdirAll(filePath,os.ModePerm)
		}
		if  err != nil  {
			fmt.Printf("%s日志文件夹创建失败：%s\n",v,err)
			continue
		}else{
    		f ,err = os.Create(filePath + filename)
		}
		if err != nil {
			fmt.Printf("%s日志文件打开失败：%s\n",v,err)
			continue
		}
		
		l.file = f
		allLog[v] = l
	}
}