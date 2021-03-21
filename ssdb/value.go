package ssdb

/**
code copy from github.com/seefan/gossdb/client  v1.1.3
*/
import (
	"encoding/json"
	"strconv"
	"time"
)

//Value string
//
//扩展值，原始类型为 string
type Value string

//String 返回 string 的值
func (v Value) String() string {
	return string(v)
}

//Int64 返回 int64 的值
func (v Value) Int64() int64 {
	if f, e := strconv.ParseInt(string(v), 10, 64); e == nil {
		return f
	}
	return 0
}

//Int32 返回 int32 的值
func (v Value) Int32() int32 {
	return int32(v.Int64())
}

//Int16 返回 int16 的值
func (v Value) Int16() int16 {
	return int16(v.Int64())
}

//Int8 返回 int8 的值
func (v Value) Int8() int8 {
	return int8(v.Int64())
}

//Int 返回 int 的值
func (v Value) Int() int {
	return int(v.Int64())
}

//UInt64 返回 uint64 的值
func (v Value) UInt64() uint64 {
	if f, e := strconv.ParseUint(string(v), 10, 64); e == nil {
		return f
	}
	return 0
}

//UInt32 返回 uint32 类型的值
func (v Value) UInt32() uint32 {
	return uint32(v.UInt64())
}

//UInt16 返回 uint16 类型的值
func (v Value) UInt16() uint16 {
	return uint16(v.UInt64())
}

//UInt8 返回 uint8 类型的值
func (v Value) UInt8() uint8 {
	return uint8(v.UInt64())
}

//Byte 返回 byte 类型的值
func (v Value) Byte() byte {
	return v.UInt8()
}

//UInt 返回 uint 类型的值
func (v Value) UInt() uint {
	return uint(v.UInt64())
}

//Float64 返回 float64 类型的值
func (v Value) Float64() float64 {
	if f, e := strconv.ParseFloat(string(v), 64); e == nil {
		return f
	}
	return 0
}

//Float32 返回 float32 类型的值
func (v Value) Float32() float32 {
	return float32(v.Float64())
}

//Bool 返回 bool 类型的值
func (v Value) Bool() bool {
	return v == "1"
}

//Time 返回 time.Time 类型的值
func (v Value) Time() time.Time {
	return time.Unix(v.Int64(), 0)
}

//Duration 返回 time.Duration 类型的值
func (v Value) Duration() time.Duration {
	return time.Duration(v.Int64())
}

//Bytes 返回 []byte 类型的值
func (v Value) Bytes() []byte {
	return []byte(v)
}

//IsEmpty 判断是否为空
func (v Value) IsEmpty() bool {
	return v == ""
}

//As 按json 转换指定类型
//
//  value 传入的指针
//
//示例
//  var abc time.Time
//  v.As(&abc)
func (v Value) As(value interface{}) (err error) {
	err = json.Unmarshal(v.Bytes(), value)
	return
}
