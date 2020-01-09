package lib

import(
    "websocket/app/log"
    "strings"
    "fmt"
    "io/ioutil"
)

var appConfig map[string]string


func Config(key string) string{
    if v,ok := appConfig[key]; ok{
        return v
    }
    appConfig = make(map[string]string)
    bytes,err := ioutil.ReadFile("app/app.conf")
    if  err != nil {
        log.Write("error",fmt.Sprintf("读取文件(app.conf)失败：%s",err))
        fmt.Println("读取文件(app.conf)失败：",err)
        return ""
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
                    appConfig[ck] = strings.Trim(strings.Trim(cv,`"`),`'`)
                }
            }
             
        }
        if v,ok := appConfig[key]; ok{
            return v
        }else{
            return ""
        }
        
    }else{
        return ""
    }
}