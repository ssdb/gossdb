package main

import (
	"fmt"
	"github.com/ssdb/gossdb/ssdb"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestData(t *testing.T) {
	must := require.New(t)
	ip := "127.0.0.1"
	port := 8888
	db, err := ssdb.Connect(ip, port)
	must.NoError(err)

	defer db.Close()

	//keys := []string{}
	//keys = append(keys, "c")
	//keys = append(keys, "d")
	//val, err = db.Do("multi_get", "a", "b", keys)
	mval, err := db.MultiGet("a", "b", "c", "d")
	must.NoError(err)
	fmt.Printf("%s\n", mval)

	err = db.Set("a", "xxx")
	must.NoError(err)

	val, err := db.Get("a")
	must.NoError(err)
	//fmt.Printf("%s\n", val)
	must.Equal(val.String(), "xxx")

	err = db.Del("a")
	must.NoError(err)

	val, err = db.Get("a")
	must.NoError(err)
	//fmt.Printf("%s\n", val)
	must.Equal(val.String(), "")

	fmt.Printf("----\n")

	//db.Do("zset", "z", "a", 3)
	err = db.ZSet("z", "a", 3)
	must.NoError(err)

	//db.Do("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	err = db.MultiZSet("z", map[string]int64{"b": -2, "c": 5, "d": 3})
	must.NoError(err)

	resp, err := db.Do("zrscan", "z", "", "", "", 10)
	must.NoError(err)
	must.False(len(resp)%2 != 1)

	fmt.Printf("Status: %s\n", resp[0])
	for i := 1; i < len(resp); i += 2 {
		fmt.Printf("  %s : %3s\n", resp[i], resp[i+1])
	}

	fmt.Println("==zrscan.int64")
	respKey, respVal, err := db.ZRScan("z", "", "", "", 10)
	for i, n := 0, len(respVal); i < n; i++ {
		fmt.Printf("  %s : %d\n", respKey[i], respVal[i])
	}

	//_ = db.Send("dump", "", "", "-1");
	_ = db.Send("sync140")
	// receive multi responses on one request
	for {
		resp, _ := db.Recv()
		fmt.Printf("%s\n", strconv.Quote(fmt.Sprintf("%s", resp)))
	}

	return
}
