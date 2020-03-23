package gossdb

import "strconv"

//:hclear
func (c *Client) Hclear(key string) error {
	_, err := c.Do("hclear", key)
	return err
}

//:hget
func (c *Client) Hget(key string, filed string) (string, error) {
	result, err := c.Do("hget", key, filed)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:multi_hget
func (c *Client) Hmget(key string, m map[string]interface{}) ([]string, error) {
	args := []interface{}{}
	for filed, val := range m {
		args = append(args, filed)
		args = append(args, val)
	}
	result, err := c.Do("multi_hget", key, args)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:hset
func (c *Client) Hset(key string, filed string, val interface{}) error {
	_, err := c.Do("hset", key, filed, val)
	return err
}

//:multi_hset
func (c *Client) HMset(key string, filed string, val interface{}) error {
	_, err := c.Do("multi_hset", key, filed, val)
	return err
}

//:hdel
func (c *Client) Hdel(key string, filed string) error {
	_, err := c.Do("hdel", key, filed)
	return err
}

//:multi_hdel
func (c *Client) Hmdel(key string, filed ...string) error {
	_, err := c.Do("multi_hdel", key, filed)
	return err
}

//:hincr
func (c *Client) Hincr(key string, filed string, increment int64) (int64, error) {
	result, err := c.Do("hincr", key, filed, increment)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:hsize
func (c *Client) Hsize(key string) (int64, error) {
	result, err := c.Do("hsize", key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:hexists
func (c *Client) Hexists(key string, filed string) (bool, error) {
	result, err := c.Do("hexists", key, filed)
	if err != nil {
		return false, err
	}
	exit := true
	if result[0] == "0" {
		exit = false
	}
	return exit, nil
}

//:hdecr
func (c *Client) Hdecr(key string, filed string, decrement int64) (int64, error) {
	result, err := c.Do("hdecr", key, filed, decrement)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:hkeys
func (c *Client) Hkeys(key string) ([]string, error) {
	result, err := c.Do("hkeys", key)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:hscan
func (c *Client) Hscan(key string) ([]string, error) {
	result, err := c.Do("hscan", key)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}
