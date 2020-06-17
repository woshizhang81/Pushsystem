package main
import (
	"fmt"
	red "github.com/gomodule/redigo/redis"
	"time"
)

type Redis struct {
	pool     *red.Pool
}

var redis *Redis

func initRedis() {
	redis = new(Redis)
	redis.pool = &red.Pool{
		MaxIdle:     256,
		MaxActive:   0,
		IdleTimeout: time.Duration(120),
		Dial: func() (red.Conn, error) {
			return red.Dial(
				"tcp",
				"127.0.0.1:6379",
				red.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				red.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				red.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				red.DialDatabase(0),
				//red.DialPassword(""),
			)
		},
	}
}

func Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := redis.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func main() {
	initRedis()

	Exec("set","hello","world")
	fmt.Print(2)
	result,err := Exec("get","hello")
	if err != nil {
		fmt.Print(err.Error())
	}
	str,_:=red.String(result,err)
	fmt.Print(str)
}