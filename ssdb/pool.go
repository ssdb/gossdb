package ssdb

import (
	"net"
	"runtime"
	"time"
	"strconv"
)

type Pool struct {
	ctype    string
	clink    string
	ctimeout time.Duration
	conns    chan *Client
	config   Config
}

func NewPool(cfg Config) (*Pool, error) {

	if cfg.MaxConn < 1 {

		cfg.MaxConn = 1

	} else {

		maxconn := runtime.NumCPU() * 2
		if maxconn > 100 {
			maxconn = 100
		}

		if cfg.MaxConn > maxconn {
			cfg.MaxConn = maxconn
		}
	}

	pl := &Pool{
		ctype:    "tcp",
		clink:    cfg.Host + ":" + strconv.Itoa(cfg.Port),
		ctimeout: time.Duration(cfg.Timeout) * time.Second,
		conns:    make(chan *Client, cfg.MaxConn),
		config:   cfg,
	}

	if pl.ctimeout < 1*time.Second {
		pl.ctimeout = 10 * time.Second
	}

	for i := 0; i < cfg.MaxConn; i++ {

		cn, err := dialTimeout(pl.ctype, pl.clink, pl.ctimeout)
		if err != nil {
			return pl, err
		}
		pl.conns <- cn
	}

	return pl, nil
}

func dialTimeout(network, addr string, timeout time.Duration) (*Client, error) {

	raddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP(network, nil, raddr)
	if err != nil {
		return nil, err
	}
	sock.SetReadDeadline(time.Now().Add(timeout))
	sock.SetWriteDeadline(time.Now().Add(timeout))

	return &Client{sock: sock}, nil
}

func (pl *Pool) Cmd(args ...interface{}) *Reply {

	cn, _ := pl.pull()
	defer pl.push(cn)

	return cn.Cmd(args...)
}

func (pl *Pool) Close() {

	for i := 0; i < pl.config.MaxConn; i++ {
		cn, _ := pl.pull()
		cn.Close()
	}
}

func (pl *Pool) push(cn *Client) {
	pl.conns <- cn
}

func (pl *Pool) pull() (cn *Client, err error) {
	return <-pl.conns, nil
}
