package gossdb

//:qpush_front
func (c *Client) QpushFront(key string, val interface{}) error {
	_, err := c.Do("qpush_front", key, val)
	return err
}

//:qset
func (c *Client) Qset(key string, index int64, val interface{}) error {
	_, err := c.Do("qset", key, index, val)
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

//:qget
func (c *Client) Qget(key string, index int64) (string, error) {
	result, err := c.Do("qget", key, index)
	if err != nil {
		return "", err
	}
	return result[0], nil
}
