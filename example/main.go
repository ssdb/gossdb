package main

import (
	"fmt"
	"github.com/Quantumoffices/gossdb"
)

func main() {
	db, err := gossdb.Connect("116.62.245.150:6389")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	hash(db)
}

//k-v example
func string(ssdb *gossdb.Client) {
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
