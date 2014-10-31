package main

import (
	"./ssdb"
	"fmt"
)

func main() {

	db, err := ssdb.NewPool(ssdb.Config{
		Host:    "127.0.0.1",
		Port:    8888,
		Timeout: 3,  // timeout in second, default to 10
		MaxConn: 10, // max connection number, default to 1
	})
	if err != nil {
		fmt.Println("Connect Error:", err)
		return
	}
	defer db.Close()

	// API::Bool() bool
	if db.Cmd("set", "aa", "val-aaaaaaaaaaaaaa").Bool() {
		fmt.Println("set OK")
	}
	// API::String() string
	if rs := db.Cmd("get", "aa"); rs.State == "ok" {
		fmt.Println("get OK\n\t", rs.String())
	}
	// API::Hash() map[string]string
	db.Cmd("set", "bb", "val-bbbbbbbbbbbbbbbbbb")
	if rs := db.Cmd("multi_get", "aa", "bb"); rs.State == "ok" {
		fmt.Println("multi_get OK")
		for k, v := range rs.Hash() {
			fmt.Println("\t", k, v)
		}
	}

	db.Cmd("zset", "z", "a", 3)
	db.Cmd("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	if rs := db.Cmd("zrscan", "z", "", "", "", 10); rs.State == "ok" {
		fmt.Println("zrscan OK")
		for k, v := range rs.Hash() {
			fmt.Println("\t", k, v)
		}
	}

	db.Cmd("set", "key", 10)
	if rs := db.Cmd("incr", "key", 1).Int(); rs > 0 {
		fmt.Println("incr OK\n\t", rs)
	}

	// API::Int() int
	// API::Int64() int64
	db.Cmd("setx", "key", 123456, 300)
	if rs := db.Cmd("ttl", "key").Int(); rs > 0 {
		fmt.Println("ttl OK\n\t", rs)
	}

	if rs := db.Cmd("multi_hset", "zone", "c1", "v-01", "c2", "v-02"); rs.State == "ok" {
		fmt.Println("multi_hset OK")
	}
	if rs := db.Cmd("multi_hget", "zone", "c1", "c2"); rs.State == "ok" {
		fmt.Println("multi_hget OK")
		for k, v := range rs.Hash() {
			fmt.Println("\t", k, v)
		}
	}

	// API::List() []string
	db.Cmd("qpush", "queue", "q-1111111111111")
	db.Cmd("qpush", "queue", "q-2222222222222")
	if rs := db.Cmd("qpop", "queue", 10); rs.State == "ok" {
		fmt.Println("qpop OK")
		for k, v := range rs.List() {
			fmt.Println("\t", k, v)
		}
	}
}
