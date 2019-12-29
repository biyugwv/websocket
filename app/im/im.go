package im
 
import (
		"github.com/gorilla/websocket"
		"fmt"
		"encoding/json"
		"time"
		"strconv"
	)

const  (
    pingMaxTime = 12
)
type Client struct {
    Conn *websocket.Conn    // 用户websocket连接
    SSID string             // 登陆后的验证Token
    PingTime int64          // 心跳监测 客户端发起  5秒收不到自动断开
    PongTime int64          // 收到客户端ping消息后发送到客户端
    Userid string           // 登录用户名
    ConnectTime int64
    lastMsgTime int64
    stat int
}

type Message struct {
    EventType int8  `json:"type"`       //  -2 表示pong； -1 表ping ； 0表示用户发布消息；1表示用户进入；2表示用户退出；
    Name string     `json:"name"`       // 用户名称
    Message string  `json:"msg"`    // 消息内容
}

// 接收到的消息
type receiveMsg  struct {
    MsgType int8  `json:"type"`       
    Msg string     `json:"msg"`       
    SSID string  `json:"SSID"`    
}

var (
    Clients map[string] Client
    Join  chan Client
    Leave chan Client
    Msg   chan Message
)



func channelLoop(){  // 三个channel为空 读取操作会造成阻塞
    for{
        
        select {
            case msg := <- Msg:
                for _,client := range Clients {
                    
                    data, err := json.Marshal(msg)
                    if err != nil {
                        fmt.Println("json.Marshal 错误")
                        return
                    }
                    
                    if err := client.Conn.WriteMessage(websocket.TextMessage, data) ; err != nil {  // 转换成字符串类型便于查看
                        fmt.Println(err)
                    }
                    fmt.Println(string(data))
                    fmt.Printf("当前在线人数：%d\n",len(Clients))
                }
            case client := <- Join:
                Clients[client.SSID] = client
                go readMsg(client)  // 监听connection发送消息
                msg := Message{1,"system",fmt.Sprintf("%s进入了房间",client.SSID)}
                Msg<-msg
                
            case client := <- Leave:
                msg := Message{2,"system",fmt.Sprintf("%s离开了房间",client.SSID)}
                Msg<-msg
                delete(Clients, client.SSID)
               
        }
        
    }
}


//  
func readMsg(client Client){
    for{
        if _,ok :=  Clients[client.SSID] ; !ok {  
            fmt.Println("one client remove")
            client.Conn.Close()
            break
        }
        now  := time.Now().Unix()
        // 超时没有收到客户端的心跳监测消息 将处理掉该链接
        if client.PingTime == 0 {
            client.PingTime = now
        } 
        if now - client.PingTime>= pingMaxTime{  
            Leave <- client
            return 
        }
        // 接收数据
        var (
            data []byte
            err error
        )
        if _,data,err = client.Conn.ReadMessage() ; err!=nil {
            Leave <- client
            return 
        }
        msgget := receiveMsg{}
        json.Unmarshal(data,&msgget)
        
        if msgget.MsgType == -1 { // 心跳数据接收
            client.PingTime = now
            pong(client)  // 向客户端发送监测数据
        }else if  msgget.MsgType == 0 {      // 此处会报interface{} is  nil  , not int  处理办法是先进行nil判断
            msg := Message{0,client.SSID,msgget.Msg}
            Msg<-msg
        }else {
            fmt.Println("invilid data：%s",string(data))
        }
    }
}


func pong(client Client){
    now  := time.Now().Unix()
    client.PongTime = now
    msg := Message{-2,"system",strconv.FormatInt(now,10)}
    data, err := json.Marshal(msg)
    if err != nil {
        fmt.Println("pong err - 1")
        return
    }
    if err := client.Conn.WriteMessage(websocket.TextMessage, data) ;err != nil {
        
    }
}


func Run(){
    // 初始化
    Clients = make(map [string] Client)      
    Join = make(chan Client ,100)
    Leave = make(chan Client ,100)
    Msg = make(chan Message ,1000)

    go channelLoop()

}