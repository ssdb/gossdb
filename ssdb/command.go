package ssdb

import "strconv"

func (c *Client) Set(key, val string) (bool, error) {
	resp, err := c.Do("set", key, val)
	if err != nil {
		return false, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return true, nil
	}
	return false, ErrBadResponse
}

func (c *Client) Setx(key, val string, timeout int) (bool, error) {
	resp, err := c.Do("setx", key, val, timeout)
	if err != nil {
		return false, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return true, nil
	}
	return false, ErrBadResponse
}

func (c *Client) Setnx(key, val string) (exists bool, err error) {
	resp, err := c.Do("setnx", key, val)
	if err != nil {
		return false, nil
	}
	if len(resp) == 2 && resp[0] == "ok" {
		if resp[1] == "1" {
			return true, nil
		}
		return false, nil
	}
	return false, ErrBadResponse
}

func (c *Client) Get(key string) (string, error) {
	resp, err := c.Do("get", key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", nil
	}
	return "", ErrBadResponse
}

func (c *Client) Del(key string) (int, error) {
	resp, err := c.Do("del", key)
	if err != nil {
		return 0, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, ErrBadResponse
}

func (c *Client) MultiDel(keys []string) (int, error) {
	resp, err := c.Do("multi_del", keys)
	if err != nil {
		return 0, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, ErrBadResponse
}

func (c *Client) Zset(name, key string, score int) (bool, error) {
	resp, err := c.Do("zset", name, key, score)
	if err != nil {
		return false, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return true, nil
	}
	return false, ErrBadResponse
}

func (c *Client) Zget(name, key string) (string, error) {
	resp, err := c.Do("zget", name, key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", ErrNotFound
	}
	return "", ErrBadResponse
}

// return del score count
func (c *Client) Zdel(name, key string) (int, error) {
	resp, err := c.Do("zdel", name, key)
	if err != nil {
		return 0, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, ErrBadResponse
}

func (c *Client) Zkeys(name, keyStart string, scoreStart, scoreEnd, limit int) ([]string, error) {
	resp, err := c.Do("zkeys", name, keyStart, scoreStart, scoreEnd, limit)
	if err != nil {
		return nil, err
	}
	if resp[0] == "ok" {
		return resp[1:], nil
	}
	return nil, ErrBadResponse
}

// return cleared keys count
func (c *Client) Zclear(name string) (int, error) {
	resp, err := c.Do("zclear", name)
	if err != nil {
		return 0, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	return 0, ErrBadResponse
}
