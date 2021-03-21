package ssdb

/**
code copy from github.com/seefan/gossdb/client  v1.1.3
*/

import (
	"github.com/pkg/errors"
)

//HSet 设置 hashmap 中指定 key 对应的值内容.
//
//  setName hashmap 的名字
//  key hashmap 的 key
//  value key 的值
//  返回 err，执行的错误
func (c *Client) HSet(setName, key string, value interface{}) (err error) {
	resp, err := c.Do("hset", setName, key, value)
	if err != nil {
		return errors.WithMessagef(err, "Hset %s/%s error ", setName, key)
	}

	if len(resp) > 0 && resp[0] == oK {
		return nil
	}
	return makeError(resp, setName, key)
}

//HGet 获取 hashmap 中指定 key 的值内容.
//
//  setName hashmap 的名字
//  key hashmap 的 key
//  返回 value key 的值
//  返回 err，执行的错误
func (c *Client) HGet(setName, key string) (value Value, err error) {
	resp, err := c.Do("hget", setName, key)
	if err != nil {
		return "", errors.WithMessagef(err, "Hget %s/%s error", setName, key)
	}
	if len(resp) == 2 && resp[0] == oK {
		return Value(resp[1]), nil
	}
	return "", makeError(resp, setName, key)
}

//HDel 删除 hashmap 中的指定 key，不能通过返回值来判断被删除的 key 是否存在.
//
//  setName hashmap 的名字
//  key hashmap 的 key
//  返回 err，执行的错误
func (c *Client) HDel(setName, key string) (err error) {
	resp, err := c.Do("hdel", setName, key)
	if err != nil {
		return errors.WithMessagef(err, "Hdel %s/%s error", setName, key)
	}
	if len(resp) > 0 && resp[0] == oK {
		return nil
	}
	return makeError(resp, setName, key)
}

//HExists 判断指定的 key 是否存在于 hashmap 中.
//
//  setName hashmap 的名字
//  key hashmap 的 key
//  返回 re，如果当前 key 不存在返回 false
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HExists(setName, key string) (re bool, err error) {
	resp, err := c.Do("hexists", setName, key)
	if err != nil {
		return false, errors.WithMessagef(err, "Hexists %s/%s error", setName, key)
	}

	if len(resp) == 2 && resp[0] == oK {
		return resp[1] == "1", nil
	}
	return false, makeError(resp, setName, key)
}

//HClear 删除 hashmap 中的所有 key
//
//  setName hashmap 的名字
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HClear(setName string) (err error) {
	resp, err := c.Do("hclear", setName)
	if err != nil {
		return errors.WithMessagef(err, "Hclear %s error", setName)
	}

	if len(resp) > 0 && resp[0] == oK {
		return nil
	}
	return makeError(resp, setName)
}

//HScan 列出 hashmap 中处于区间 (key_start, key_end] 的 key-value 列表. ("", ""] 表示整个区间.
//
//  setName - hashmap 的名字.
//  keyStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  keyEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回包含 key-value 的关联字典.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HScan(setName string, keyStart, keyEnd string, limit int64, reverse ...bool) (map[string]Value, error) {
	cmd := "hscan"
	if len(reverse) > 0 && reverse[0] {
		cmd = "hrscan"
	}

	resp, err := c.Do(cmd, setName, keyStart, keyEnd, limit)

	if err != nil {
		return nil, errors.WithMessagef(err, "%s %s %s %s %v error", cmd, setName, keyStart, keyEnd, limit)
	}

	if len(resp) > 0 && resp[0] == oK {
		re := make(map[string]Value)
		size := len(resp)
		for i := 1; i < size-1; i += 2 {
			re[resp[i]] = Value(resp[i+1])
		}
		return re, nil
	}
	return nil, makeError(resp, setName, keyStart, keyEnd, limit)
}

//HScanArray 列出 hashmap 中处于区间 (key_start, key_end] 的 key,value 列表. ("", ""] 表示整个区间.
//
//  setName - hashmap 的名字.
//  keyStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  keyEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回包含 key-value 的关联字典.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HScanArray(setName string, keyStart, keyEnd string, limit int64, reverse ...bool) ([]string, []Value, error) {
	cmd := "hscan"
	if len(reverse) > 0 && reverse[0] {
		cmd = "hrscan"
	}
	resp, err := c.Do(cmd, setName, keyStart, keyEnd, limit)

	if err != nil {
		return nil, nil, errors.WithMessagef(err, "%s %s %s %s %v error", cmd, setName, keyStart, keyEnd, limit)
	}

	if len(resp) > 0 && resp[0] == oK {
		keys := []string{}
		values := []Value{}
		size := len(resp)
		for i := 1; i < size-1; i += 2 {
			keys = append(keys, resp[i])
			values = append(values, Value(resp[i+1]))
		}
		return keys, values, nil
	}
	return nil, nil, makeError(resp, setName, keyStart, keyEnd, limit)
}

//HRScanArray 列出 hashmap 中处于区间 (key_start, key_end] 的 key,value 列表. ("", ""] 表示整个区间.
//
//  setName - hashmap 的名字.
//  keyStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  keyEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回包含 key-value 的关联字典.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HRScanArray(setName string, keyStart, keyEnd string, limit int64, reverse ...bool) ([]string, []Value, error) {
	return c.HScanArray(setName, keyStart, keyEnd, limit, true)
}

//HRScan 列出 hashmap 中处于区间 (key_start, key_end] 的 key-value 列表. ("", ""] 表示整个区间.
//
//  setName - hashmap 的名字.
//  keyStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  keyEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回包含 key-value 的关联字典.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HRScan(setName string, keyStart, keyEnd string, limit int64) (map[string]Value, error) {
	return c.HScan(setName, keyStart, keyEnd, limit, true)
}

//MultiHSet 批量设置 hashmap 中的 key-value.
//
//  setName - hashmap 的名字.
//  kvs - 包含 key-value 的关联数组 .
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHSet(setName string, kvs map[string]interface{}) (err error) {

	args := []interface{}{"multi_hset", setName}
	for k, v := range kvs {
		args = append(args, k)
		args = append(args, v)
	}
	resp, err := c.Do(args...)

	if err != nil {
		return errors.WithMessagef(err, "MultiHset %s %s error", setName, kvs)
	}

	if len(resp) > 0 && resp[0] == oK {
		return nil
	}
	return makeError(resp, setName, kvs)
}

//MultiHGet 批量获取 hashmap 中多个 key 对应的权重值.
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
//  返回 包含 key-value 的关联数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGet(setName string, key ...string) (val map[string]Value, err error) {
	if len(key) == 0 {
		return make(map[string]Value), nil
	}

	args := make([]interface{}, len(key)+2)
	args[0] = "multi_hget"
	args[1] = setName

	for i, v := range key {
		args[i+2] = v
	}

	resp, err := c.Do(args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "MultiHget %s %s error", setName, key)
	}
	size := len(resp)
	if size > 0 && resp[0] == oK {
		val = make(map[string]Value)
		for i := 1; i < size && i+1 < size; i += 2 {
			val[resp[i]] = Value(resp[i+1])
		}
		return val, nil
	}
	return nil, makeError(resp, key)
}

//MultiHGetSlice 批量获取 hashmap 中多个 key 对应的权重值.
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
//  返回 包含 key和value 的有序数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGetSlice(setName string, key ...string) (keys []string, values []Value, err error) {
	if len(key) == 0 {
		return []string{}, []Value{}, nil
	}
	args := []interface{}{"multi_hget", setName}
	for _, v := range key {
		args = append(args, v)
	}
	resp, err := c.Do(args...)

	if err != nil {
		return nil, nil, errors.WithMessagef(err, "MultiHgetSlice %s %s error", setName, key)
	}
	if len(resp) > 0 && resp[0] == oK {
		size := len(resp)
		keys := make([]string, 0, (size-1)/2)
		values := make([]Value, 0, (size-1)/2)

		for i := 1; i < size && i+1 < size; i += 2 {
			keys = append(keys, resp[i])
			values = append(values, Value(resp[i+1]))
		}
		return keys, values, nil
	}
	return nil, nil, makeError(resp, key)
}

//MultiHGetArray 批量获取 hashmap 中多个 key 对应的权重值.（输入分片）
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
//  返回 包含 key-value 的关联数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGetArray(setName string, key []string) (val map[string]Value, err error) {
	return c.MultiHGet(setName, key...)
}

//MultiHGetSliceArray 批量获取 hashmap 中多个 key 对应的权重值.（输入分片）
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
//  返回 包含 key和value 的有序数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGetSliceArray(setName string, key []string) (keys []string, values []Value, err error) {
	return c.MultiHGetSlice(setName, key...)
}

//MultiHGetAll 批量获取 hashmap 中全部 对应的权重值.
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组 .
//  返回 包含 key-value 的关联数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGetAll(setName string) (val map[string]Value, err error) {

	resp, err := c.Do("hgetall", setName)

	if err != nil {
		return nil, errors.WithMessagef(err, "MultiHgetAll %s error", setName)
	}
	size := len(resp)
	if size > 0 && resp[0] == oK {
		val = make(map[string]Value)
		for i := 1; i < size && i+1 < size; i += 2 {
			val[resp[i]] = Value(resp[i+1])
		}
		return val, nil
	}
	return nil, makeError(resp)
}

//MultiHGetAllSlice 批量获取 hashmap 中全部 对应的权重值.
//
//  setName - hashmap 的名字.
//  返回 包含 key和value 的有序数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHGetAllSlice(setName string) (keys []string, values []Value, err error) {

	resp, err := c.Do("hgetall", setName)

	if err != nil {
		return nil, nil, errors.WithMessagef(err, "MultiHgetAllSlice %s error", setName)
	}
	if len(resp) > 0 && resp[0] == oK {
		size := len(resp)
		keys := make([]string, 0, (size-1)/2)
		values := make([]Value, 0, (size-1)/2)

		for i := 1; i < size && i+1 < size; i += 2 {
			keys = append(keys, resp[i])
			values = append(values, Value(resp[i+1]))
		}
		return keys, values, nil
	}
	return nil, nil, makeError(resp)
}

//MultiHDel 批量删除 hashmap 中的 key.
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHDel(setName string, key ...string) (err error) {
	if len(key) == 0 {
		return nil
	}
	args := []interface{}{"multi_hdel", setName}
	for _, v := range key {
		args = append(args, v)
	}
	resp, err := c.Do(args...)
	if err != nil {
		return errors.WithMessagef(err, "MultiHdel %s %s error", setName, key)
	}

	if len(resp) > 0 && resp[0] == oK {
		return nil
	}
	return makeError(resp, key)
}

//MultiHDelArray 批量删除 hashmap 中的 key.（输入分片）
//
//  setName - hashmap 的名字.
//  keys - 包含 key 的数组.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) MultiHDelArray(setName string, key []string) (err error) {
	return c.MultiHDel(setName, key...)
}

//HList 列出名字处于区间 (name_start, name_end] 的 hashmap. ("", ""] 表示整个区间.
//
//  nameStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  nameEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回 包含名字的数组
//  返回 err，执行的错D，操作成功返回 nil
func (c *Client) HList(nameStart, nameEnd string, limit int64) ([]string, error) {
	resp, err := c.Do("hlist", nameStart, nameEnd, limit)
	if err != nil {
		return nil, errors.WithMessagef(err, "Hlist %s %s %v error", nameStart, nameEnd, limit)
	}

	if len(resp) > 0 && resp[0] == oK {
		size := len(resp)
		keyList := make([]string, 0, size-1)

		for i := 1; i < size; i++ {
			keyList = append(keyList, resp[i])
		}
		return keyList, nil
	}
	return nil, makeError(resp, nameStart, nameEnd, limit)
}

//HIncr 设置 hashmap 中指定 key 对应的值增加 num. 参数 num 可以为负数.
//
//  setName - hashmap 的名字.
//  key 键值
//  num 增加的值
//  返回 val，整数，增加 num 后的新值
//  返回 err，可能的错误，操作成功返回 nil
func (c *Client) HIncr(setName, key string, num int64) (val int64, err error) {

	resp, err := c.Do("hincr", setName, key, num)

	if err != nil {
		return -1, errors.WithMessagef(err, "Hincr %s error", key)
	}
	if len(resp) == 2 && resp[0] == oK {
		return Value(resp[1]).Int64(), nil
	}
	return -1, makeError(resp, key)
}

//HSize 返回 hashmap 中的元素个数.
//
//  setName - hashmap 的名字.
//  返回 val，整数，增加 num 后的新值
//  返回 err，可能的错误，操作成功返回 nil
func (c *Client) HSize(setName string) (val int64, err error) {

	resp, err := c.Do("hsize", setName)

	if err != nil {
		return -1, errors.WithMessagef(err, "Hsize %s error", setName)
	}
	if len(resp) == 2 && resp[0] == oK {
		return Value(resp[1]).Int64(), nil
	}
	return -1, makeError(resp, setName)
}

//HKeys 列出 hashmap 中处于区间 (keyStart, keyEnd] 的 key 列表.
//
//  name - hashmap 的名字.
//  keyStart - 返回的起始 key(不包含), 空字符串表示 -inf.
//  keyEnd - 返回的结束 key(包含), 空字符串表示 +inf.
//  limit - 最多返回这么多个元素.
//  返回 包含名字的数组
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HKeys(setName, keyStart, keyEnd string, limit int64) ([]string, error) {
	resp, err := c.Do("hkeys", setName, keyStart, keyEnd, limit)
	if err != nil {
		return nil, errors.WithMessagef(err, "Hkeys %s %s %s %v error", setName, keyStart, keyEnd, limit)
	}

	if len(resp) > 0 && resp[0] == oK {
		return resp[1:], nil
	}
	return nil, makeError(resp, keyStart, keyEnd, limit)
}

//HGetAll 批量获取 hashmap 中全部 对应的权重值.
//
//  setName - hashmap 的名字.
//  返回 包含 key-value 的关联数组, 如果某个 key 不存在, 则它不会出现在返回数组中.
//  返回 err，执行的错误，操作成功返回 nil
func (c *Client) HGetAll(setName string) (val map[string]Value, err error) {
	return c.MultiHGetAll(setName)
}
