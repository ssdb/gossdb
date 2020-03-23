package main

import (
	"fmt"
	"github.com/Quantumoffices/gossdb"
	"strconv"
)

func main() {
	db, err := gossdb.Connect("116.62.245.150:6389")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//hash(db)
	zset(db)
}

//k-v example
func kv(ssdb *gossdb.Client) {
	v, err := ssdb.Hincr("a", "b", 100)
	if err != nil {
		panic(err)
	}
	fmt.Println("ssdb.Hincr = ", v)

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
	size, err := ssdb.Zremrangebyrank("scores", 0, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("szie = ", size)
	m := map[string]int64{
		"a": 1200,
		"c": 2300,
	}

	v, err := ssdb.Zmultiset("scores", m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
	//
	scores, err := ssdb.Zrange("scores", 0, -1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(scores)
}
