package ssdb

import (
	"strconv"
)

const (
	ReplyOK          string = "ok"
	ReplyNotFound    string = "not_found"
	ReplyError       string = "error"
	ReplyFail        string = "fail"
	ReplyClientError string = "client_error"
)

type Reply struct {
	State string
	Data  []string
}

func (c *Client) Cmd(args ...interface{}) *Reply {

	r := &Reply{
		State: ReplyError,
		Data:  []string{},
	}

	if err := c.send(args); err != nil {
		return r
	}

	resp, err := c.recv()
	if err != nil || len(resp) < 1 {
		return r
	}

	switch resp[0] {
	case ReplyOK, ReplyNotFound, ReplyError, ReplyFail, ReplyClientError:
		r.State = resp[0]
	}

	if r.State == ReplyOK {
		for k, v := range resp {

			if k == 0 {
				continue
			}

			r.Data = append(r.Data, v)
		}
	}

	return r
}

func (r *Reply) String() string {

	if len(r.Data) > 0 {
		return r.Data[0]
	}

	return ""
}

func (r *Reply) Int() int {
	return int(r.Int64())
}

func (r *Reply) Int64() int64 {

	if len(r.Data) < 1 {
		return 0
	}

	i64, err := strconv.ParseInt(r.Data[0], 10, 64)
	if err == nil {
		return i64
	}

	return 0
}

func (r *Reply) Bool() bool {
	return r.Int64() == 1
}

func (r *Reply) List() []string {

	if len(r.Data) < 1 {
		return []string{}
	}

	return r.Data
}

func (r *Reply) Hash() map[string]string {

	hs := map[string]string{}

	if len(r.Data) < 2 {
		return hs
	}

	for i := 0; i < (len(r.Data) - 1); i += 2 {
		hs[r.Data[i]] = r.Data[i+1]
	}

	return hs
}
