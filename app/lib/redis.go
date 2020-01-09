package lib

import(
    
    //"websocket/app/log"
    "github.com/go-redis/redis/v7"
    "time"
)


func GetRedis()(client  *redis.Client) {
    client = redis.NewClient(&redis.Options{
		Addr:     Config("redisurl"),
		Password: Config("redispwd"),
		DB:       10,
	})
	if  _, err := client.Ping().Result() ; err !=nil{
    	panic(err)
	}
	return 
}


func RedisGet(key string) string {
		
        client := GetRedis()
        
        val , err := client.Get(key).Result() 
        if  err !=nil {
            panic(err)
        }
        return  val
}

func RedisPut(key string,val string,duration int64) bool{
        client := GetRedis()
        timeoutDuration := time.Duration(duration)  * time.Second
        rs,err := client.SetNX(key,val,timeoutDuration).Result()
        if err!=nil{
            return false
        }
        return  rs
}