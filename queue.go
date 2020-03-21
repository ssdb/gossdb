package gossdb

import (
	"strconv"
)

//:qpush_front
func (c *Client) QpushFront(key string, val interface{}) error {
	_, err := c.Do("qpush_front", key, val)
	return err
}

//:qpop_front
func (c *Client) QpopFront(key string) (string, error) {
	result, err := c.Do("qpop_front", key)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:qpush_back
func (c *Client) Qpushback(key string, val interface{}) error {
	_, err := c.Do("qpush_back", key, val)
	return err
}

//:qpop_back
func (c *Client) Qpopback(key string) (string, error) {
	result, err := c.Do("qpop_back", key)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:qset
func (c *Client) Qset(key string, index int64, val interface{}) error {
	_, err := c.Do("qset", key, index, val)
	return err
}

//:qget
func (c *Client) Qget(key string, index int64) (string, error) {
	result, err := c.Do("qget", key, index)
	if err != nil {
		return "", err
	}
	return result[0], nil
}

//:qsize
func (c *Client) Qsize(key string) (int64, error) {
	result, err := c.Do("qsize", key)
	if err != nil {
		return 0, err
	}
	size, err := strconv.ParseInt(result[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

//:qclear
func (c *Client) Qclear(key string) error {
	_, err := c.Do("qclear", key)
	return err
}

//:qslice
func (c *Client) Qslice(key string, start, end int64) ([]string, error) {
	result, err := c.Do("qslice", key, start, end)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:qlist (for queue/list type only)
//func (c *Client) Qlist() ([]string, error) {
//	result, err := c.Do("qlist")
//	if err != nil {
//		return []string{}, err
//	}
//	return result, nil
//}
