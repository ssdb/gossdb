package ssdb

import (
	"math/rand"
	"testing"
)

func makeValue(size int) []byte {
	val := make([]byte, size)
	for i := 0; i < size; i++ {
		str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		val[i] = str[rand.Int31n(int32(62))]
	}
	return val
}

func benchmarkSSDBGoRecv(valSize, batchSize int, b *testing.B) {
	c, e := Connect("127.0.0.1", 8888)
	if e != nil {
		b.Fatal(e)
	}
	defer c.Close()
	testKey1 := "ssdb.benchmark.key1"
	defer c.Do("hclear", testKey1)
	for n := 0; n < batchSize; n++ {
		c.Do("hset", testKey1, n, makeValue(valSize))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, e := c.Do("hscan", testKey1, "", "", batchSize)
		if e != nil {
			b.Error(e)
		}
	}
}

// go test -bench=.
// go test -bench=BenchmarkSSDBGoRecv_1k_10 -benchtime=20s
// go test -run=none -bench=BenchmarkSSDBGoRecv_10k_50 -cpuprofile=cprof
func BenchmarkSSDBGoRecv_64b_10(b *testing.B) {
	benchmarkSSDBGoRecv(64, 10, b)
}
func BenchmarkSSDBGoRecv_64b_50(b *testing.B) {
	benchmarkSSDBGoRecv(64, 50, b)
}

func BenchmarkSSDBGoRecv_1k_10(b *testing.B) {
	benchmarkSSDBGoRecv(1*1024, 10, b)
}
func BenchmarkSSDBGoRecv_1k_50(b *testing.B) {
	benchmarkSSDBGoRecv(1*1024, 50, b)
}

func BenchmarkSSDBGoRecv_10k_10(b *testing.B) {
	benchmarkSSDBGoRecv(10*1024, 10, b)
}

func BenchmarkSSDBGoRecv_10k_50(b *testing.B) {
	benchmarkSSDBGoRecv(10*1024, 50, b)
}
