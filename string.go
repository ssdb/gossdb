package gossdb

//:set
func (c *Client) Set(key string, val string) error {
	_, err := c.Do("set", key, val)
	return err
}

//:get
func (c *Client) Get(key string) (interface{}, error) {
	result, err := c.Do("get", key)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}
