package gossdb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
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
	err := c.send(args)
	if err != nil {
		return nil, err
	}
	resp, err := c.recv()
	return resp, err
}

func (c *Client) Del(key string) (interface{}, error) {
	resp, err := c.Do("del", key)
	if err != nil {
		return nil, err
	}

	//response looks like this: [ok 1]
	if len(resp) > 0 && resp[0] == "ok" {
		return true, nil
	}
	return nil, fmt.Errorf("bad response:resp:%v:", resp)
}

func (c *Client) Send(args ...interface{}) error {
	return c.send(args)
}

func (c *Client) send(args []interface{}) error {
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

func (c *Client) Recv() ([]string, error) {
	return c.recv()
}

func (c *Client) recv() ([]string, error) {
	var tmp [8192]byte
	n, err := c.conn.Read(tmp[0:])
	if err != nil {
		return nil, err
	}
	c.recvBuf.Write(tmp[0:n])
	//fmt.Println("---------start------------")
	//fmt.Println(string(tmp[0:n]))
	//fmt.Println("---------end------------")
	//resp := c.parse()
	//resp := c.parse_1()
	//if resp == nil || len(resp) > 0 {
	//	return resp, nil
	//}
	return c.parse(tmp[0:n])
}

//参数解析
//func (c *Client) parse() []string {
//	resp := []string{}
//	buf := c.recvBuf.Bytes()
//	var idx, offset int
//	idx = 0
//	offset = 0
//	for {
//		idx = bytes.IndexByte(buf[offset:], '\n')
//		if idx == -1 {
//			break
//		}
//		p := buf[offset : offset+idx]
//		offset += idx + 1
//		//fmt.Printf("> [%s]\n", p);
//		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
//			if len(resp) == 0 {
//				continue
//			} else {
//				var new_buf bytes.Buffer
//				new_buf.Write(buf[offset:])
//				c.recvBuf = new_buf
//				return resp
//			}
//		}
//		size, err := strconv.Atoi(string(p))
//		if err != nil || size < 0 {
//			return nil
//		}
//		if offset+size >= c.recvBuf.Len() {
//			break
//		}
//
//		v := buf[offset : offset+size]
//		resp = append(resp, string(v))
//		offset += size + 1
//	}
//	//fmt.Printf("buf.size: %d packet not ready...\n", len(buf))
//	return []string{}
//}

//success if error==nil
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

// Close The Client Connection
func (c *Client) Close() error {
	return c.conn.Close()
}
