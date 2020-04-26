package main

import (
	"fmt"
	"github.com/Quantumoffices/gossdb"
	"strconv"
	"time"
)

func main() {
	db, err := gossdb.NewClient(gossdb.Options{
		Addr:              "192.168.0.220:6389",
		OnConnect:         nil,
		ReadTimeout:       0,
		WriteTimeout:      0,
		PoolSize:          0,
		MinIdleConns:      0,
		MaxConnAge:        0,
		PoolTimeout:       0,
		IdleTimeout:       0,
		DialTimeout:       time.Second * 3,
		ReconnectCount:    15,
		ReconnectDuration: time.Second * 5,
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//hash(db)
	//Db(db)
	kv(db)
}
func Db(ssdb *gossdb.Client) {
	keys, err := ssdb.Keys(-1)
	if err != nil {
		panic(err)
	}
	fmt.Println(keys)
}

//k-v example
func kv(ssdb *gossdb.Client) {
	for {
		time.Sleep(time.Second)
		go func() {
			//err := ssdb.Set("a", "bbb")
			//if err != nil {
			//	panic(err)
			//}
			v, err := ssdb.Get("a")
			if err != nil {
				fmt.Println(err)
			}
			if v == "1" {
				fmt.Println()
			}
			fmt.Println("ssdb.a =", v)
		}()
	}
}

//hash example
func hash(ssdb *gossdb.Client) {
	key := "a"
	v1, err := ssdb.Hincr(key, "b", -100)
	if err != nil {
		panic(err)
	}
	fmt.Println("ssdb.Hincr = ", v1)

	err = ssdb.Hset(key, "c", "www")
	if err != nil {
		panic(err)
	}
	v2, err := ssdb.Hget(key, "c")
	if err != nil {
		panic(err)
	}
	fmt.Println("ssdb.Hget = ", v2)
}

func zset(ssdb *gossdb.Client) {
	for i := 0; i < 10; i++ {
		count, err := ssdb.Zset("scores", int64(i), strconv.Itoa(i))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("zset count ", count)
	}
	//err := ssdb.Flushdb()
	scores, err := ssdb.Zrange("scores", 0, -1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("scores = ", scores)
}
