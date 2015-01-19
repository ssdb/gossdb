package main

import (
	"fmt"
	"os"
	"strconv"
	"./ssdb"
)

func main() {
	ip := "127.0.0.1"
	port := 8888
	db, err := ssdb.Connect(ip, port)
	if err != nil {
		fmt.Errorf("ssdb.Connect:err:%v:\n", err)
		os.Exit(1)
	}

	defer db.Close()
	var val interface{}

	keys := []string{}
	keys = append(keys, "c");
	keys = append(keys, "d");
	val, err = db.Do("multi_get", "a", "b", keys);
	fmt.Printf(":%s:\n", val);
	if err != nil {
		os.Exit(1)
	}

	val, err = db.Do("info");
	fmt.Printf("called info:%s:\n", val);
	if err != nil {
		os.Exit(1)
	}

	val, err = db.Info();
	fmt.Printf("called info again:%s:\n", val);
	if err != nil {
		os.Exit(1)
	}

	val, err = db.Do("keys", "", "", 100);
	fmt.Printf("called keys:%s:\n", val);
	if err != nil {
		os.Exit(1)
	}

	val, err = db.Do("scan", "", "", 100);
	fmt.Printf("called scan:%s:\n", val);
	if err != nil {
		os.Exit(1)
	}

	db.Set("a", "xxx")
	val, err = db.Get("a")
	fmt.Printf("Got:%v:\n", val)
	_, err = db.Del("a")
	fmt.Printf("deleted it:err:%v:\n", err)
	val, err = db.Get("a")
	fmt.Printf("got it again:val:%v:err:%v:\n", val, err)

	fmt.Printf("----\n");

	db.Do("zset", "z", "a", 3)
	db.Do("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	resp, err := db.Do("zrscan", "z", "", "", "", 10)
	if err != nil {
		os.Exit(1)
	}
	if len(resp)%2 != 1 {
		fmt.Printf("bad response")
		os.Exit(1)
	}

	fmt.Printf("Status: %s\n", resp[0])
	for i := 1; i < len(resp); i += 2 {
		fmt.Printf("  %s : %3s\n", resp[i], resp[i+1])
	}

	fmt.Printf("call:init pool\n")
	cpm, err := ssdb.InitPool(ip, port, 11)
	fmt.Printf("done call:init pool\n")

	if err != nil {
		fmt.Printf("failed to init pool:%v:", err)
		os.Exit(1)
	}

	// cpm.Put
	fmt.Printf("cpm.Set:a:xxx\n")
	cpm.Set("a", "xxx")

	fmt.Printf("cpm.Get:a:\n")
	val, err = cpm.Get("a")
	fmt.Printf("cpm.Get:return:val:%s:err:%v:\n", val, err)

	cpm.Close()


	//_ = db.Send("dump", "", "", "-1");
	_ = db.Send("sync140");
	// receive multi responses on one request
	for{
		resp, _ := db.Recv()
		fmt.Printf("%s\n", strconv.Quote(fmt.Sprintf("%s", resp)));
	}

	return
}
