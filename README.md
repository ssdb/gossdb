SSDB Go API Documentation {#mainpage}
============

@author: [ideawu](http://www.ideawu.com/)

## About

All SSDB operations go with ```ssdb.Client.Do()```, it accepts variable arguments. The first argument of Do() is the SSDB command, for example "get", "set", etc. The rest arguments(maybe none) are the arguments of that command.

The Do() method will return an array of string if no error. The first element in that array is the response code, ```"ok"``` means the following elements in that array(maybe none) are valid results. The response code may be ```"not_found"``` if you are calling "get" on an non-exist key.

Refer to the [PHP documentation](http://www.ideawu.com/ssdb/docs/php/) to checkout a complete list of all avilable commands and corresponding responses.

## gossdb is not thread-safe(goroutine-safe)

Never use one connection(returned by ssdb.Connect()) through multi goroutines, because the connection is not thread-safe.

## Example
```go
        db, err := gossdb.Connect("116.62.245.150:6389")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
    		count, err := ssdb.Zset("scores", int64(i), strconv.Itoa(i))
    		if err != nil {
    			fmt.Println(err.Error())
    		}
    		fmt.Println("zset count ", count)
    	}
    	//err := ssdb.Flushdb()
    	scores, err := ssdb.Zrange("scores", 0, -1)
    	if err != nil {
    		fmt.Println(err)
    	}
    	fmt.Println("scores = ", scores)
```
	
#限制
 * 最大 Key 长度	200 字节  
 * 最大 Value 长度	31MB  
 * 最大请求或响应长度	31MB  
 * 单个 HASH 中的元素数量	9,223,372,036,854,775,807  
 * 单个 ZSET 中的元素数量	9,223,372,036,854,775,807  
 * 单个 QUEUE 中的元素数量	9,223,372,036,854,775,807  
 * 命令最多参数个数	所有参数加起来体积不超过 31MB 大小  

