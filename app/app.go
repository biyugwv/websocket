package app

import(
    "websocket/app/im"
    "websocket/app/log"
    "github.com/gorilla/websocket"
    "net/http"
    "fmt"
    "time"
    "strings"
    "io/ioutil"
    
)

var Config map[string]string

func loadConfig(){
    Config = make(map[string]string)
    bytes,err := ioutil.ReadFile("app/app.conf")
    if  err != nil {
        log.Write("error",fmt.Sprintf("读取文件(app.conf)失败：%s",err))
        fmt.Println("读取文件(app.conf)失败：",err)
        Config["stat"] = "fail"
        return
    }
    s := string(bytes)
    if len(s) > 0 {
        b := strings.Split(s,"\n")
        for  _,v := range b{
            if v =  strings.TrimSpace(v) ; v !="" {
                c := strings.Split(v,"=")
                ck := strings.TrimSpace(c[0])
                cv := strings.TrimSpace(c[1])
                if len(c) == 2 && ck!="" && cv!="" {
                    Config[ck] = strings.Trim(cv,`"`)
                     
                }
            }
             
        }

        if _,ok := Config["stat"] ; !ok{
            Config["stat"] = "success"
        }
    }else{
        Config["stat"] = "none"
    }
    
    
}

func recordRequest(r *http.Request ){
    var requestInfo string
    requestInfo = "userAgent-" + r.UserAgent() +"|" +"referer-" + r.Referer() +"|" +"remoteAddr-" + r.RemoteAddr +"|"+ "host-" + r.Host +"|" + "requestURI-" + r.RequestURI
    log.Write("info",requestInfo)
}

func ws(w http.ResponseWriter, r *http.Request){
    var SSID string
    recordRequest(r)
    r.ParseForm()
    if len(r.FormValue("SSID")) > 0 {
        SSID = r.FormValue("SSID")
    }else{
        fmt.Println("need SSID")
        log.Write("info","need SSID")
        return 
    }
    conn, err := (&websocket.Upgrader{
		// 允许跨域
		CheckOrigin : func(r *http.Request) bool{
			return true
		},

	}).Upgrade(w, r, nil)
	
	if _, ok := err.(websocket.HandshakeError); ok {
        fmt.Println("HandshakeError")
        return
    } else if err != nil {
        fmt.Println(err)
        return
    }
    client := &im.Client{}
    client.SSID = SSID
    client.Conn = conn
    client.ConnectTime = time.Now().Unix()
    // 如果用户列表中没有该用户
    if _,ok := im.Clients[SSID] ; !ok  {
        im.Join <- client
    }else{
        fmt.Println("该用户已在线")
    }
}


func Run(){
    log.LogInit()     // 日志初始化
    loadConfig()  // 加载配置文件
    im.Run()     //  im初始化 等待conn连接进来
    go FileListener()  // 监听文件变化 暂时没想好干嘛用
    
    http.HandleFunc("/", ws)
    http.ListenAndServe(":"+Config["port"], nil)
    fmt.Println("端口监听：",Config["port"])
}