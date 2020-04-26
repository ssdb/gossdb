package gossdb

import (
	"strconv"
)

//:set Set the value of the key.
//false on error, other values indicate OK.
func (c *Client) Set(key string, val interface{}) error {
	_, err := c.Do("set", key, val)
	return err
}

//:setnx Set the string value in argument as value of the key if and only if the key doesn't exist.
//false on error, 1: value is set, 0: key already exists.
func (c *Client) Setnx(key string, val interface{}) (bool, error) {
	result, err := c.Do("setnx", key, val)
	if err != nil {
		return false, err
	}
	ok := false
	if result[0] == "1" {
		ok = true
	}
	return ok, nil
}

//:setx Set the value of the key, with a time to live.
//false on error, other values indicate OK.
func (c *Client) Setx(key string, val interface{}, ttl int64) error {
	_, err := c.Do("setx", key, val, ttl)
	return err
}

//:get
func (c *Client) Get(key string) (string, error) {
	result, err := c.Do("get", key)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:incr Increment the number stored at key by num. The num argument could be a negative integer. The old number is first converted to an integer before increment, assuming it was stored as literal integer.
//false on error, other values the new value
func (c *Client) Incr(key string, incrment int64) (int64, error) {
	result, err := c.Do("incr", key, incrment)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:getset Sets a value and returns the previous entry at that key.
//Returns null if key not found, false on error, therwise, the value related to that key is returned.
func (c *Client) Getset(key string, val interface{}) (string, error) {
	result, err := c.Do("getset", key, val)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:ttl Returns the time left to live in seconds, only for keys of KV type.
//false on error, or time to live of the key, in seconds, -1 if there is no associated expire to the key.
func (c *Client) TTL(key string) (int64, error) {
	result, err := c.Do("ttl", key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:expire Set the time left to live in seconds, only for keys of KV type.
//false on error, if key exists and ttl is set, return 1, if key not exists, return 0.
func (c *Client) Expire(key string) (bool, error) {
	result, err := c.Do("expire", key)
	if err != nil {
		return false, err
	}
	ok := false
	if result[0] == "1" {
		ok = true
	}
	return ok, nil
}

//:del Delete specified key.
//false on error, other values indicate OK. You can not determine whether the key exists or not.
func (c *Client) Del(key string, val interface{}) error {
	_, err := c.Do("Del", key, val)
	return err
}
