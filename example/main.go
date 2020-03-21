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
	//if err := db.QpushFront("a", "1"); err != nil {
	//	panic(err)
	//}
	for i := 0; i < 10; i++ {
		err := db.QpushFront("a", i)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
	//for i := 0; i < 10; i++ {
	//	err = db.Qset("a", int64(i), 99)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		panic(err)
	//	}
	//}
	for i := 0; i < 10; i++ {
		v, err := db.Qget("a", int64(i))
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		fmt.Println("qpop_fornt:", v)
	}
	//k-v
	//_, err = db.Do("set", "a", "www")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//reply, err := db.Do("get", "a")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(len(reply))
	//fmt.Println(reply)
	//return
	//HSET key field value
	//db.Do("hset", "a", "1", "2020")
	//db.Do("hget", "a", "1")
	//
	//db.Do("zset", "a", "1", "aa")
	//db.Do("zset", "a", "2", "bb")
	//db.Do("zsize", "a")
	//
	//db.Do("qpush_front", "a", "fff", "eee")
	//db.Do("qslice", "a", "0", "-1")
	//
	//var val interface{}
	//keys := []string{}
	//keys = append(keys, "c")
	//keys = append(keys, "d")
	//val, err = db.Do("multi_get", "a", "b", keys)
	//fmt.Printf("%s\n", val)
	//db.Set("a", "xxx")
	//val, err = db.Get("a")
	//fmt.Printf("%s\n", val)
	////db.Del("a")
	//val, err = db.Get("aaa")
	//fmt.Printf("%v\n", val)
	////return
	//fmt.Printf("----\n")
	//
	//reply, err = db.Do("zset", "z", "a", 3)
	//if err != nil {
	//	fmt.Println(reply)
	//}
	//db.Do("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	//resp, err := db.Do("zrscan", "z", "", "", "", 10)
	//if err != nil {
	//	os.Exit(1)
	//}
	//if len(resp)%2 != 1 {
	//	fmt.Printf("bad response")
	//	os.Exit(1)
	//}
	//
	//fmt.Printf("Status: %s\n", resp[0])
	//for i := 1; i < len(resp); i += 2 {
	//	fmt.Printf("  %s : %3s\n", resp[i], resp[i+1])
	//}
	//
	//_ = db.Send("dump", "", "", "-1")
	//_ = db.Send("sync140")
	//// receive multi responses on one request
	//for {
	//	resp, _ := db.Recv()
	//	fmt.Printf("%s\n", strconv.Quote(fmt.Sprintf("%s", resp)))
	//}
	//return
}
