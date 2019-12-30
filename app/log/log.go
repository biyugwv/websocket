package log

import(
	//"log"
	"time"
	"os"
	"sync"
	"fmt"

)
const(
	INFO string = "info"        // 消息
	WARNING string = "warning"  // 警告
	ERROR string = "error"      // 错误
	MESSAGE string = "message"   //  websocket通讯记录
)

type appLog struct{
	updateTime int64  //  更新时间
	file  *os.File
	logType  string       // 日志类型
	logBuffer string      // 日志写入buffer
	chBuffer chan string
	lockFile  *sync.RWMutex
	maxLen int
	maxTime int64
}
var allLog = make( map[string] *appLog )
var  logtypes = []string{INFO,WARNING,ERROR,MESSAGE}


// 将msg写入buffer 

func write(logtype string , msg string){
	var inarr bool = false
	for _,v := range logtypes{
		if  v == logtype {
			inarr = true
			break
		}
	}
	if !inarr {
		fmt.Printf("logtype 不存在\n")
		return
	}
	allLog[logtype].chBuffer<-msg
}
func writeLog(){
	for{
		now := time.Now().Unix()
		for _,v := range allLog{
			if len(v.logBuffer) >= v.maxLen || now - v.updateTime >= v.maxTime{
				// 将buff写入log文件
				if len(v.logBuffer) > 1{
					f , err := getFile(v.logType) 
					if err !=nil {
						fmt.Printf("文件创建或读取失败 \n")
						
					}
					v.file = f
					fmt.Println("写入文件")
					v.logBuffer = ""
					v.updateTime = now
				}
			}
			select{
				case msg:= <- v.chBuffer : 
					v.logBuffer = v.logBuffer +msg  + "\n"
				default :
					break
					
			}
		}

		
	}
}


func getFile(logtype string)(f *os.File ,err error){
	
	filePath := "../websocket/logs/"+time.Now().Format("20160102")+"/"
	filename := logtype+".log"
	f ,err = os.Open(filePath + filename)
	if err != nil {
		err = os.MkdirAll(filePath,os.ModePerm)
	}else{
		return 
	}
	if  err != nil  {
		fmt.Printf("%s日志文件夹创建失败：%s\n",logtype,err)
		return
	}else{
		f ,err = os.Create(filePath + filename)
	}
	if err != nil {
		fmt.Printf("%s日志文件打开失败：%s\n",logtype,err)
		return
	}
	f.Close()
	return
}


func LogInit(){
	
	for _,v := range logtypes {
		l := &appLog{}
		l.logType = v
		l.updateTime = time.Now().Unix()
		l.maxLen =  1000
		l.maxTime =  10
		l.chBuffer = make(chan string,10)
		
		allLog[v] = l
	}

	go writeLog()
}