package app

import(
    //"os/exec"
    "fmt"
    "websocket/app/log"
    "github.com/fsnotify/fsnotify"
    "io/ioutil"
    "path/filepath"
)

const rootpath = "../websocket"

func listenerDir(path string ,watch *fsnotify.Watcher){
    path,_ = filepath.Abs(path)
    watch.Add(path)
    dir, err := ioutil.ReadDir(path)
    if err != nil {
        log.Write("error",fmt.Sprintf("%s not a dir",path))
        return
    }
    for _, dirchild := range dir {
        if dirchild.IsDir() { // 目录, 递归遍历
            if string(dirchild.Name()[0]) !="."{
                pathchild := path+"/"+dirchild.Name()
                listenerDir(pathchild ,watch)
            }
        } 
    }
}


func FileListener() {
    watch, err := fsnotify.NewWatcher();
    if err != nil {
        log.Write("error",fmt.Sprintf("文件监听失败，错误如下：%s",err))
    }
    defer watch.Close();
    
    listenerDir(rootpath,watch)
    
    if err != nil {
        log.Write("error",fmt.Sprintf("文件监听失败，错误如下：%s",err))
    }
   
    for {
        select {
        case ev := <-watch.Events:
            {
               
                if ev.Op&fsnotify.Create == fsnotify.Create {
                    log.Write("info",fmt.Sprintf("创建文件 : ", ev.Name));
                }
                if ev.Op&fsnotify.Write == fsnotify.Write {
                    log.Write("info",fmt.Sprintf("写入文件 : ", ev.Name));
                    
                }
                if ev.Op&fsnotify.Remove == fsnotify.Remove {
                    log.Write("info",fmt.Sprintf("删除文件 : ", ev.Name));
                    
                }
                if ev.Op&fsnotify.Rename == fsnotify.Rename {
                    log.Write("info",fmt.Sprintf("重命名文件 : ", ev.Name));
                    
                }
                if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
                    log.Write("info",fmt.Sprintf("修改权限 : ", ev.Name));
                }
            }
        case err := <-watch.Errors:
            {
                log.Write("error", fmt.Sprintf("%s",err));
                return;
            }
        }
    }
    
    
}

