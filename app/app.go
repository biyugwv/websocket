package app

import(
    "websocket/app/im"
    "websocket/app/log"
    "websocket/app/lib"
    "github.com/gorilla/websocket"
    "net/http"
    "fmt"
    "time"
    
    
)



func recordRequest(r *http.Request ){
    var requestInfo string
    requestInfo = "userAgent-" + r.UserAgent() +"|" +"referer-" + r.Referer() +"|" +"remoteAddr-" + r.RemoteAddr +"|"+ "host-" + r.Host +"|" + "requestURI-" + r.RequestURI
    log.Write("info",requestInfo)
}

// 登录获取SSID

func login(w http.ResponseWriter, r *http.Request){
    recordRequest(r)
    r.ParseForm()
    w.Header().Set("Access-Control-Allow-Origin","*")
    var username,password string
    if len(r.FormValue("username")) > 0 {
        username = r.FormValue("username")
    }else{
        w.Write([]byte("请填写姓名"))
    }
    if len(r.FormValue("password")) > 0 {
        password = r.FormValue("password")
    }else{
        w.Write([]byte("请填写密码"))
    }
    log.Write("info",fmt.Sprintf("登录：username(%s),password(%p)",username,password))
    if (username == "admin" && password == "admin123") || (username == "admin" && password == "admin456") {
        SSID := lib.AESEncodeStr( username ,"login_to_get_sid")
        lib.RedisPut(SSID,username,3600)
        w.Write([]byte(SSID))
    }else{
        w.Write([]byte("账号密码错误"))
    }
}

func ws(w http.ResponseWriter, r *http.Request){
    recordRequest(r)
    var SSID string
    r.ParseForm()
    if len(r.FormValue("SSID")) > 0 {
        SSID = r.FormValue("SSID")
        if uname := lib.RedisGet(SSID); uname =="" {
            fmt.Printf("SSID验证失败：%s\n",SSID)
            log.Write("info",fmt.Sprintf("SSID验证失败：%s",SSID))
            return
        }
    }else{
        fmt.Printf("need SSID")
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
    }else{  //挤掉
        im.Leave <- im.Clients[SSID]
        im.Join <- client
    }
}


func Run(){
    log.LogInit()     // 日志初始化
    im.Run()     //  im初始化 等待conn连接进来
    go FileListener()  // 监听文件变化 暂时没想好干嘛用
    
    http.HandleFunc("/", ws)
    http.HandleFunc("/login", login)
    log.Write("info",fmt.Sprintf("端口监听：%s",lib.Config("port")))
    http.ListenAndServe(":"+ lib.Config("port"), nil)
}