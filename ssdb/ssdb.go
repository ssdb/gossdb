package ssdb

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	sock chan *net.TCPConn
	recv_buf bytes.Buffer
	_sock  *net.TCPConn
}

type ConnectionPoolWrapper struct {
	size int
	conn chan *Client
}

func InitPool(ip string, port int, size int) (*ConnectionPoolWrapper, error) {

    cpm := new(ConnectionPoolWrapper)

	cpm.conn = make(chan *Client, size)
	for x := 0; x < size; x++ {
		conn, err := Connect(ip, port)
		if err != nil {
			return cpm, err
		}
 
		// If the init function succeeded, add the connection to the channel
		cpm.conn <- conn
	}
	cpm.size = size
	return cpm, nil

}

func (p *ConnectionPoolWrapper) GetConnection() *Client {
	return <-p.conn
}
 
func (p *ConnectionPoolWrapper) ReleaseConnection(conn *Client) {
	p.conn <- conn
}

func Connect(ip string, port int) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	var c Client

	c.sock = make(chan *net.TCPConn, 1)
	// fmt.Printf("Connect:putting socket:\n")
	c.sock <- sock
	// fmt.Printf("Connect:done putting socket:\n")
	return &c, nil
}

func (c *Client) Do(args ...interface{}) ([]string, error) {

	// fmt.Printf("Do:pulling socket\n")
	c._sock = <- c.sock
	// fmt.Printf("Do:done pulling socket\n")
	defer func () { 
		// fmt.Printf("Do:putting socket\n")
		c.sock <- c._sock 
		// fmt.Printf("Do:done putting socket\n")

	}()

	return c.do(args...)
}

func (c *Client) do(args ...interface{}) ([]string, error) {
	err := c.send(args)
	if err != nil {
		return nil, err
	}
	resp, err := c.recv()
	return resp, err
}

func (c *Client) Set(key string, val string) (interface{}, error) {
	resp, err := c.Do("set", key, val)
	if err != nil {
		return nil, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return true, nil
	}
	return nil, fmt.Errorf("bad response")
}


func (c *ConnectionPoolWrapper) Set(key string, val string) (interface{}, error) {

	db := c.GetConnection()
	defer c.ReleaseConnection(db)

	return db.Set(key, val)
}

// TODO: Will somebody write addition semantic methods?
func (c *Client) Get(key string) (interface{}, error) {
	resp, err := c.Do("get", key)
	if err != nil {
		return nil, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return nil, nil
	}
	return nil, fmt.Errorf("bad response")
}

func (c *ConnectionPoolWrapper) Get(key string) (interface{}, error) {

	db := c.GetConnection()
	defer c.ReleaseConnection(db)

	return db.Get(key)
}

func (c *Client) Del(key string) (interface{}, error) {
	resp, err := c.Do("del", key)
	if err != nil {
		return nil, err
	}

	//response looks like this: [ok 1]
	if len(resp) == 2 && resp[0] == "ok" {
		return true, nil
	}
	return nil, fmt.Errorf("bad response:resp:%v:", resp)
}

func (c *ConnectionPoolWrapper) Del(key string) (interface{}, error) {

	db := c.GetConnection()
	defer c.ReleaseConnection(db)

	return db.Del(key)
}

func (c *Client) send(args []interface{}) error {

	var sock = c._sock
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

	// fmt.Printf("send:pulling socket\n")
	// sock := <- c.sock
	// fmt.Printf("send:done pulling socket\n")
	// defer func () { 
		// fmt.Printf("send:putting socket\n")
		// c.sock <- sock 
		// fmt.Printf("send:done putting socket\n")

	// }()

	_, err := sock.Write(buf.Bytes())

	return err
}

func (c *Client) recv() ([]string, error) {

	var sock = c._sock
	var tmp [8192]byte

	// fmt.Printf("recv:pulling socket\n")
	// sock := <- c.sock
	// fmt.Printf("recv:done pulling socket\n")

	// defer func () { 
		// fmt.Printf("recv:putting socket\n")
		// c.sock <- sock 
		// fmt.Printf("recv:done putting socket\n")

	// }()

	for {
		n, err := sock.Read(tmp[0:])
		if err != nil {

			return nil, err
		}
		c.recv_buf.Write(tmp[0:n])
		resp := c.parse()
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
	}
}

func (c *Client) parse() []string {
	resp := []string{}
	buf := c.recv_buf.Bytes()
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
		//fmt.Printf("> [%s]\n", p);
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				c.recv_buf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= c.recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	return []string{}
}

// Close The Client Connection
func (c *Client) Close() error {

	// fmt.Printf("Close:pulling socket\n")
	sock := <- c.sock
	// fmt.Printf("Close:done pulling socket\n")
	// defer func () { c.sock <- sock }()
		// fmt.Printf("recv:done pulling socket\n")

	defer func () { 
		// fmt.Printf("Close:putting socket\n")
		c.sock <- sock 
		// fmt.Printf("Close:done putting socket\n")

	}()

	return sock.Close()

}


func (cpm *ConnectionPoolWrapper) Close() error {

	for {

		select {
			case db := <- cpm.conn:
				// fmt.Printf("ConnectionPoolWrapper.Close:closing one\n")
				db.Close()
			default:
				// fmt.Printf("ConnectionPoolWrapper.Close:skipping\n")
				return nil
		}

	}

	return nil
}
