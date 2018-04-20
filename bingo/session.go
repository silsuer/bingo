package bingo


// 操作session


type SessionProvider struct {
	ID string   // 每一个客户端的连接，都会分配一个唯一的ID
	Session Session  // 对于每个连接，都分配一个Session结构体，用来对这个唯一的连接进行Session操作
}


type Session struct {
   Set func(key string,value interface{}) error  // 设置session
   Get func(key string) (interface{},error)  // 获取session
   Destroy func(key string) error // 销毁session中的值
   GC func(maxTime int64)  // 根据session的有效时间，自动销毁session
}

var sessionMap map[string]interface{}  // 当使用map存储session的时候，使用这个变量

// 当使用file存储的时候

// 引入bolt 持久化key-value数据