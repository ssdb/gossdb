package gossdb

import "strconv"

//:zclear Delete all keys in a zset.
//false on error, or the number of keys deleted
func (c *Client) Zclear(key string) error {
	_, err := c.Do("zclear", key)
	return err
}

//:zget Get the score related to the specified key of a zset
//Returns null if key not found, false on error, otherwise, the score related to this key is returned.
func (c *Client) Zget(key string, member string) (int64, error) {
	result, err := c.Do("zget", key, member)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:multi_zset Set the score of the key of a zset.
//false on error, other values indicate OK.
//func (c *Client) Zmultiset(key string, m map[string]int64) (int64, error) {
//	args := []interface{}{}
//	for k, v := range m {
//		args = append(args, k)
//		args = append(args, v)
//	}
//	result, err := c.Do("multi_zset", args)
//	if err != nil {
//		return 0, err
//	}
//	return strconv.ParseInt(result[0], 10, 64)
//}

//:zset Set multiple key-score pairs(kvs) of a zset in one method call.
//false on error, other values indicate OK.
func (c *Client) Zset(key string, score int64, member string) (int64, error) {
	result, err := c.Do("zset", key, score, member)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zdel Delete specified key of a zset.
//false on error, other values indicate OK. You can not determine whether the key exists or not.
func (c *Client) Zdel(key string, member ...string) (int64, error) {
	result, err := c.Do("zdel", key, member)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zcount Returns the number of elements of the sorted set stored at the specified key which have scores in the range [start,end].
//false on error, or the number of keys in specified range.
func (c *Client) zcount(key string) (int64, error) {
	result, err := c.Do("zcount", key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zrange Returns a range of key-score pairs by index range [offset, offset + limit). Zrrange iterates in reverse order.
//false on error, otherwise an array containing key-score pairs.
func (c *Client) Zrange(key string, start, end int64) ([]string, error) {
	result, err := c.Do("zrange", key, start, end)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:zrrange Returns a range of key-score pairs by index range [offset, offset + limit). Zrrange iterates in reverse order.
//false on error, otherwise an array containing key-score pairs.
func (c *Client) Zrrange(key string, start, end int64) ([]string, error) {
	result, err := c.Do("zrrange", key, start, end)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:zscan [min,max] List key-score pairs in a zset, where key-score in range (key_start+score_start, score_end]. If key_start is empty, keys with a score greater than or equal to score_start will be returned. If key_start is not empty, keys with score larger than score_start, and keys larger than key_start also with score equal to score_start will be returned.
//That is: return keys in (key.score == score_start && key > key_start || key.score > score_start), and key.score <= score_end. The score_start, score_end is of higher priority than key_start.
//("", ""] means no range limit.
//false on error, otherwise an associative array containing the key-score pairs.
func (c *Client) Zscan(key string, min, max int64) ([]string, error) {
	result, err := c.Do("zscan", key, "", min, max, -1)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:zrscan [max,min]
func (c *Client) Zrscan(key string, min, max int64) ([]string, error) {
	result, err := c.Do("zrscan", key, "", max, min, -1)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:zsize Return the number of pairs of a zset.
//false on error, otherwise an integer, 0 if the zset does not exist.
func (c *Client) Zsize(key string) (int64, error) {
	result, err := c.Do("zsize", key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zincr Increment the number stored at key in a zset by num. The num argument could be a negative integer. The old number is first converted to an integer before increment, assuming it was stored as literal integer.
//false on error, other values the new value.
func (c *Client) Zincr(key string, member string, increment int64) (int64, error) {
	result, err := c.Do("zincr", key, member, increment)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zexists Verify if the specified key exists in a zset.
//If the key exists, return true, otherwise return false.
func (c *Client) Zexists(key string, member string) (bool, error) {
	result, err := c.Do("zexists", key, member)
	if err != nil {
		return false, err
	}
	var exit bool
	if result[0] == "1" {
		exit = true
	}
	return exit, nil
}

//:zrank Returns the rank(index) of a given key in the specified sorted set, starting at 0 for the item with the smallest score. zrrank starts at 0 for the item with the largest score.
//false on error, otherwise the rank(index) of a specified key, start at 0. null if not found.
func (c *Client) Zrank(key string, member string) (int64, error) {
	result, err := c.Do("zrank", key, member)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zrrank Returns the rank(index) of a given key in the specified sorted set, starting at 0 for the item with the smallest score. zrrank starts at 0 for the item with the largest score.
//false on error, otherwise the rank(index) of a specified key, start at 0. null if not found.
func (c *Client) Zrrank(key string, member string) (int64, error) {
	result, err := c.Do("zrrank", key, member)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zcount  Returns the number of elements of the sorted set stored at the specified key which have scores in the range [start,end].
//false on error, or the number of keys in specified range.
func (c *Client) Zcount(key string, start, end int64) (int64, error) {
	result, err := c.Do("zcount", key, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zsum Returns the sum of elements of the sorted set stored at the specified key which have scores in the range [start,end].
//false on error, or the sum of keys in specified range.
func (c *Client) Zsum(key string, start, end int64) (int64, error) {
	result, err := c.Do("zsum", key, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zavg Returns the average of elements of the sorted set stored at the specified key which have scores in the range [start,end].
//false on error, or the average of keys in specified range.
func (c *Client) Zavg(key string, start, end int64) (int64, error) {
	result, err := c.Do("zavg", key, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zremrangebyrank Delete the elements of the zset which have rank in the range [start,end].
//Return Value:false on error, or the number of deleted elements.
func (c *Client) Zremrangebyrank(key string, start, end int64) (int64, error) {
	result, err := c.Do("zremrangebyrank", key, start, end)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zremrangebyscore Delete the elements of the zset which have score in the range [start,end].
//Return Value:false on error, or the number of deleted elements.
func (c *Client) Zremrangebyscore(key string, min, max int64) (int64, error) {
	result, err := c.Do("zremrangebyscore", key, min, max)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:zremrangebyscore Delete and return `limit` element(s) from back of the zset.
//Return false on error, otherwise an array containing key-score pairs.
func (c *Client) Zpopback(key string, limit int64) ([]string, error) {
	result, err := c.Do("zrscan", key, limit)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

//:zpop_front Delete and return `limit` element(s) from front of the zset.
//Return false on error, otherwise an array containing key-score pairs.
func (c *Client) Zpopfront(key string, limit int64) ([]string, error) {
	result, err := c.Do("zpop_front", key, limit)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}
