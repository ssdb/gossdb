package ssdb

import "testing"

var (
	addr = "127.0.0.1"
	port = 8888
)

func TestCommand(t *testing.T) {
	client := Client{Addr: addr, Port: port}
	ok, err := client.Set("TestString", "abc")
	if err != nil || ok == false {
		t.Fatal("set command error")
	}
	value, err := client.Get("TestString")
	if err != nil {
		t.Fatal("get command error")
	}
	if value != "abc" {
		t.Fatal("get result not match")
	}
	num, err := client.Del("TestString")
	if err != nil {
		t.Fatal("del command error")
	}
	if num != 1 {
		t.Fatal("del result error")
	}

	client.Set("TestString1", "abc")
	client.Set("TestString2", "abc")
	num, err = client.MultiDel([]string{"TestString1", "TestString2"})
	if err != nil {
		t.Fatal("multi_del command error")
	}
	if num != 2 {
		t.Fatal("multi_del result error")
	}

	ok, err = client.Zset("TestZset", "abc", 10)
	if err != nil || ok != true {
		t.Fatal("zset command error")
	}
	value, err = client.Zget("TestZset", "abc")
	if err != nil {
		t.Fatal("zget command error")
	}
	if value != "10" {
		t.Fatal("zget value not match")
	}
	// key not exist
	value, err = client.Zget("TestZset0", "abc")
	if err == nil || err != ErrNotFound {
		t.Fatal("zget when key not exist error")
	}

	_, err = client.Zdel("TestZset", "abc")
	if err != nil {
		t.Fatal("zdel command error")
	}
	// key had been deleted
	num, err = client.Zdel("TestZset", "abc")
	if err != nil || num != 0 {
		t.Fatal("zdel when key has been deleted command error")
	}

	client.Zset("TestZset", "abc", 10)
	client.Zset("TestZset", "abcd", 20)
	keys, err := client.Zkeys("TestZset", "0", 0, 100, -1)
	if err != nil {
		t.Fatal("zkeys command error")
	}
	if keys[0] != "abc" || keys[1] != "abcd" {
		t.Fatal("zkeys value not match")
	}
	num, err = client.Zclear("TestZset")
	if err != nil {
		t.Fatal("zclear command error")
	}
	if num != 2 {
		t.Fatal("zclear result error")
	}
}
