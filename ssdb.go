package gossdb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	CODE_SUCCESS = "2"
	INFO_SUCCESS = "ok"

	Status_Default      int32 = iota
	Status_Networkclose       //网络断开
)

type Client struct {
	conn     net.Conn
	recvBuf  bytes.Buffer
	tryCount int8 //断线重连尝试次数
	option   Options
	//lock sync.Mutex
	cond   *sync.Cond
	status int32
}

func NewClient(opt Options) (*Client, error) {
	conn, err := net.DialTimeout("tcp", opt.Addr, opt.DialTimeout)
	return &Client{conn: conn, option: opt, cond: sync.NewCond(new(sync.Mutex))}, err
}

func (c *Client) Do(args ...interface{}) ([]string, error) {
	err := c.send(args...)
	if err != nil {
		return nil, err
	}
	resp, err := c.recv()
	return resp, err
}

func (c *Client) send(args ...interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = fmt.Sprintf("%d", arg)
		case int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')

	if atomic.LoadInt32(&c.status) == Status_Networkclose {
		c.cond.L.Lock()
		c.cond.Wait()
		c.cond.L.Unlock()
	}

	_, err := c.conn.Write(buf.Bytes())
	//可能是网络断开
	if err != nil {
		atomic.AddInt32(&c.status, Status_Networkclose)
		//atomic.LoadInt32()
		c.conn.Close()
		//断线重连机制
		for {
			c.tryCount++
			fmt.Println("ssdb reconnect times=", c.tryCount, " time:", time.Now().String(), " at gossdb/ssdb.go:89")
			c.conn, err = net.DialTimeout("tcp", c.option.Addr, time.Millisecond*50)
			if err != nil && c.tryCount <= c.option.ReconnectCount {
				time.Sleep(c.option.ReconnectDuration)
				continue
			}
			atomic.AddInt32(&c.status, -Status_Networkclose)
			c.tryCount = 0
			c.cond.Broadcast()
			//重新发送一次
			_, err = c.conn.Write(buf.Bytes())
			break
		}
	}
	return err
}

func (c *Client) recv() ([]string, error) {
	var data [1024 * 8]byte //8M
	n, err := c.conn.Read(data[0:])
	if err != nil {
		return []string{}, err
	}
	_, err = c.recvBuf.Write(data[0:n])
	if err != nil {
		return []string{}, err
	}
	return c.parse(data[0:n])
}

//success if nil error
func (c *Client) parse(data []byte) ([]string, error) {
	contentList := strings.Split(strings.TrimRight(string(data), "\n\n"), "\n")
	size := len(contentList)
	if size < 2 {
		return nil, Error_Null_Reply
	}
	//fail cmd
	if contentList[0] != CODE_SUCCESS || contentList[1] != INFO_SUCCESS {
		tips := contentList[1]
		if size >= 4 {
			tips = contentList[1] + ":" + contentList[3]
		}
		return nil, errors.New(tips)
	}
	reply := make([]string, 0, size)
	for i := 2; i < len(contentList); i++ {
		if i%2 != 0 {
			reply = append(reply, contentList[i])
		}
	}
	return reply, nil
}

//:dbsize
func (c *Client) DBsize() (int64, error) {
	result, err := c.Do("dbsize")
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(result[0], 10, 64)
}

//:flushdb
func (c *Client) Flushdb() error {
	_, err := c.Do("flushdb")
	return err
}

//:keys List keys in range (key_start, key_end].("", ""] means no range limit.
//false on error, otherwise an array containing the keys.
//limit -1 is not limit
func (c *Client) Keys(limit int64) ([]string, error) {
	result, err := c.Do("keys", "", "", limit)
	if err != nil {
		return []string{}, err
	}
	return result, err
} //:keys List keys in range (key_start, key_end].("", ""] means no range limit.

//Verify if the specified key exists.
//If the key exists, return true, otherwise return false.
func (c *Client) Exists(key string) (bool, error) {
	result, err := c.Do("exists", key)
	if err != nil {
		return false, err
	}
	exist := false
	if result[0] == "1" {
		exist = true
	}
	return exist, err
}

// Close The Client Connection
func (c *Client) Close() error {
	return c.conn.Close()
}
