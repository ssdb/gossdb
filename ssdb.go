package gossdb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	CODE_SUCCESS = "2"
	INFO_SUCCESS = "ok"
)

type Client struct {
	conn    net.Conn
	recvBuf bytes.Buffer
}

func Connect(addr string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
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
	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) recv() ([]string, error) {
	var data [8192]byte //8M
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
func (c *Client) Keys() ([]string, error) {
	result, err := c.Do("keys")
	if err != nil {
		return []string{}, err
	}
	return result, err
}

// Close The Client Connection
func (c *Client) Close() error {
	return c.conn.Close()
}
