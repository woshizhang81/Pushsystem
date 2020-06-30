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
/*
	key-value 操作
*/
func SetKeyValue() {
	c, err := red.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	defer c.Close()
	_, err = c.Do("Set", "abc", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := red.Int(c.Do("Get", "abc"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}
	fmt.Println(r)
}

/*
	has操作
*/
func HashOperate(){
	c, err := red.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	defer c.Close()
	_, err = c.Do("HSet", "books", "abc", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := red.Int(c.Do("HGet", "books", "abc"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}
	fmt.Println(r)
}


func MSet(){
		c, err := red.Dial("tcp", "localhost:6379")
		if err != nil {
			fmt.Println("conn redis failed,", err)
			return
		}
		defer c.Close()
		_, err = c.Do("MSet", "abc", 100, "efg", 300)
		if err != nil {
			fmt.Println(err)
			return
		}
		r, err := red.Ints(c.Do("MGet", "abc", "efg"))
		if err != nil {
			fmt.Println("get abc failed,", err)
			return
		}
		for _, v := range r {
			fmt.Println(v)
		}
}

func SetExpireTime(){
		c, err := red.Dial("tcp", "localhost:6379")
		if err != nil {
			fmt.Println("conn redis failed,", err)
			return
		}
		defer c.Close()
		_, err = c.Do("expire", "abc", 10)
		if err != nil {
			fmt.Println(err)
			return
		}
}

func QueueOperate(){
		c, err := red.Dial("tcp", "localhost:6379")
		if err != nil {
			fmt.Println("conn redis failed,", err)
			return
		}
		defer c.Close()
		_, err = c.Do("lpush", "book_list", "abc", "ceg", 300)
		if err != nil {
			fmt.Println(err)
			return
		}
		r, err := red.String(c.Do("lpop", "book_list"))
		if err != nil {
			fmt.Println("get abc failed,", err)
			return
		}
		fmt.Println(r)
}