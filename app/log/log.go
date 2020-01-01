package log

import(
	//"log"
	"time"
	"os"
	"sync"
	"fmt"
	"runtime"
	"strconv"

)
const(
	INFO string = "info"        // 消息
	WARNING string = "warning"  // 警告
	ERROR string = "error"      // 错误
	MESSAGE string = "message"   //  websocket通讯记录
)

type appLog struct{
	updateTime int64  //  更新时间
	logType  string       // 日志类型
	logBuffer string      // 日志写入buffer
	chBuffer chan string
	lockFile  *sync.RWMutex
	maxLen int
	maxTime int64
}
var allLog = make( map[string] *appLog )
var  logtypes = []string{INFO,WARNING,ERROR,MESSAGE}

func Write(logtype string , msg string){
	var callfun string
	_,file,line,ok := runtime.Caller(1)
	if ok {
		callfun = file +":" + strconv.Itoa(line) + " time:" + time.Now().Format("2006-01-02 15:04:05")
	}else{
		callfun = "undefended func" + " time:" + time.Now().Format("2006-01-02 15:04:05")
	}
	msg = callfun + "::" + msg +"\n"
	go write(logtype,msg)
}
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
					if err != nil {
						fmt.Printf("文件创建或读取失败 \n")
						
					}
					v.lockFile.Lock() 
					n,err:=f.WriteString(v.logBuffer)
					v.lockFile.Unlock()
					if n != len(v.logBuffer) ||  err !=nil {
						fmt.Printf("文件写入失败 %s  %d %d\n",err,n,len(v.logBuffer))
					}
					v.logBuffer = ""
					v.updateTime = now
					
					f.Close()
				}
			}
			select{
				case msg:= <- v.chBuffer : 
					v.logBuffer = v.logBuffer + msg 
				default :
					//fmt.Printf("没有要写入的buffer\n")
					
			}
		}
		//time.Sleep ( 2000 * time.Millisecond )
	}
}


func getFile(logtype string)(f *os.File ,err error){
	
	filePath := "../websocket/logs/"+time.Now().Format("20060102")+"/"
	filename := logtype+".log"
	f ,err = os.OpenFile(filePath + filename ,os.O_APPEND|os.O_WRONLY, os.ModeAppend)
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
		l.lockFile = new(sync.RWMutex)
		
		allLog[v] = l
	}

	go writeLog()
}