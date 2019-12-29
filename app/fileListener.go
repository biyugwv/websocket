package app

import(
    //"os/exec"
    "github.com/fsnotify/fsnotify"
    "log"
    "io/ioutil"
    "path/filepath"
)

const rootpath = "../websocket"

func listenerDir(path string ,watch *fsnotify.Watcher){
    path,_ = filepath.Abs(path)
    watch.Add(path)
    dir, err := ioutil.ReadDir(path)
    if err != nil {
        log.Printf("%s not a dir\n",path)
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
        log.Fatal(err);
    }
    defer watch.Close();
    
    listenerDir(rootpath,watch)
    
    if err != nil {
        log.Println("文件监听失败，错误如下：");
        log.Println(err);
    }
   
    for {
        select {
        case ev := <-watch.Events:
            {
               
                if ev.Op&fsnotify.Create == fsnotify.Create {
                    log.Println("创建文件 : ", ev.Name);
                }
                if ev.Op&fsnotify.Write == fsnotify.Write {
                    log.Println("写入文件 : ", ev.Name);
                }
                if ev.Op&fsnotify.Remove == fsnotify.Remove {
                    log.Println("删除文件 : ", ev.Name);
                }
                if ev.Op&fsnotify.Rename == fsnotify.Rename {
                    log.Println("重命名文件 : ", ev.Name);
                }
                if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
                    log.Println("修改权限 : ", ev.Name);
                }
            }
        case err := <-watch.Errors:
            {
                log.Println("error : ", err);
                return;
            }
        }
    }
    
    
}

