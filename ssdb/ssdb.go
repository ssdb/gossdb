package ssdb

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

const (
	defaulPoolSize = 10
)

var ErrBadResponse = errors.New("bad response")
var ErrNotFound = errors.New("not_found")

type Client struct {
	Addr        string
	Port        int
	MaxPoolSize int
	pool        chan net.Conn
}

func (self *Client) Do(args ...interface{}) (data []string, err error) {
	buf, err := self.formatCommand(args)
	if err != nil {
		return nil, err
	}
	conn, err := self.popConn()
	if err != nil {
		goto End
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	data, err = recv(bufio.NewReader(conn))

End:
	self.pushConn(conn)
	return data, err
}

func (self *Client) formatCommand(args []interface{}) (bytes.Buffer, error) {
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
				buf.WriteString(strconv.Itoa(len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = strconv.Itoa(arg)
		case int64:
			s = strconv.FormatInt(arg, 10)
		case float64:
			s = strconv.FormatFloat(arg, 'f', 6, 64)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return buf, ErrBadResponse
		}
		buf.WriteString(strconv.Itoa(len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	return buf, nil
}

func recv(reader *bufio.Reader) ([]string, error) {
	var tmp [8192]byte
	var recv_buf bytes.Buffer
	for {
		n, err := reader.Read(tmp[0:])
		if err != nil {
			return nil, err
		}
		recv_buf.Write(tmp[0:n])
		resp := parse(recv_buf)
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
	}
}

func parse(recv_buf bytes.Buffer) []string {
	resp := []string{}
	buf := recv_buf.Bytes()
	var idx, offset int
	idx = 0
	offset = 0

	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				recv_buf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	return []string{}
}

func (self *Client) popConn() (net.Conn, error) {
	if self.pool == nil {
		poolSize := self.MaxPoolSize
		if poolSize == 0 {
			poolSize = defaulPoolSize
		}
		self.pool = make(chan net.Conn, poolSize)
		for i := 0; i < poolSize; i++ {
			self.pool <- nil
		}
	}

	c := <-self.pool
	if c == nil {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", self.Addr, self.Port))
		if err != nil {
			return nil, err
		}
		return conn, err
	}
	return c, nil
}

func (self *Client) pushConn(c net.Conn) {
	self.pool <- c
}
