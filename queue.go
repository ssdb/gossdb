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

//:qslice Returns a portion of elements from the queue at the specified range [begin, end]. begin and end could be negative.
//false on error, otherwise an array containing items.
func (c *Client) Qslice(key string, start, end int64) ([]string, error) {
	result, err := c.Do("qslice", key, start, end)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:qtrim_back Remove multi elements from the tail of a queue.
//false on error. Return the number of elements removed.
func (c *Client) Qtrimback(key string, count int64) ([]string, error) {
	result, err := c.Do("qtrim_back", key, count)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:qtrim_front Remove multi elements from the head of a queue.
//false on error. Return the number of elements removed.
func (c *Client) Qtrimfront(key string, count int64) ([]string, error) {
	result, err := c.Do("qtrim_front", key, count)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:qlist List list/queue names in range (name_start, name_end]. ("", ""] means no range limit.
//false on error, otherwise an array containing the names.
//example $ssdb->qlist('a', 'z', 10);
func (c *Client) Qlist(key, nameStart, nameEnd string, limit int) ([]string, error) {
	result, err := c.Do("qlist", key, nameStart, nameEnd, limit)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}
