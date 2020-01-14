package log

import(
    "testing"
    "strconv"
    "websocket/app/log"
)

func BenchmarkWrite(b *testing.B){
    log.LogInit()
    for i := 0; i < b.N; i++ {
		log.Write("info","这是一个测试数据" + strconv.Itoa(i))
	}
    
}