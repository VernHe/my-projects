package main

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var (
	redisConnPool     *redis.Pool
	connectionTimeout = 5 * time.Second
	readTimeout       = 5 * time.Second
	writeTimeout      = 5 * time.Second
)

func init() {
	redisConnPool = newPool()
}

func main() {
	conn := getConnFromPool()
	defer conn.Close()

	res, err := conn.Do("SET", "key1", "value1")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res) // OK
	// bytes 的类型是 []uint8,
	bytes, err := conn.Do("GET", "key1")
	if err != nil {
		log.Fatalln(err)
	}
	v := bytes.([]byte)
	log.Println(string(v)) // value1

}

// 通过拨号获取连接
func getConnByDial() redis.Conn {
	ctx := context.Background()
	conn, err := redis.DialContext(ctx, "tcp", "192.168.56.10:6379",
		//redis.DialUsername("root"),
		//redis.DialPassword("123456"),
		redis.DialConnectTimeout(connectionTimeout),
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout),
	)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

// 通过连接池获取连接
func getConnFromPool() redis.Conn {
	return redisConnPool.Get()
}

// 创建redis连接池
func newPool() *redis.Pool {
	return &redis.Pool{
		Dial:            dialRedis,       // 拨号建立连接的函数
		DialContext:     nil,             // 获取context
		TestOnBorrow:    testOnBorrow,    // 每次获取连接之前都会执行此方法，对连接进行检查
		MaxIdle:         3,               // 连接池最大空闲的连接数，0表示不限制
		MaxActive:       0,               // 连接池最大活跃的连接数，0表示不限制
		IdleTimeout:     0,               // 每个连接的最大空闲时间（应小于超时时长），超过这个时间会被收回，0表示不限制
		Wait:            false,           // 从连接池获取连接，若没有可用连接时是否要等待
		MaxConnLifetime: 5 * time.Minute, // 连接的最大时长，如果连接超过此时长将会被关闭
	}
}

// 连接池建立连接的方法
func dialRedis() (redis.Conn, error) {
	ctx := context.Background()
	return redis.DialContext(ctx, "tcp", "192.168.56.10:6379",
		// 下面是用户名和密码
		//redis.DialUsername("root"),
		//redis.DialPassword("123456"),
		redis.DialConnectTimeout(connectionTimeout),
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout),
	)
}

// 参数t是连接回到pool的时间
func testOnBorrow(c redis.Conn, t time.Time) error {
	// Since返回从t开始经过的时间
	if time.Since(t) < time.Minute {
		// 空闲时间小于1分钟
		return nil
	}
	// 对于空闲1分钟以上的连接进行测试
	_, err := c.Do("PING")
	return err
}
