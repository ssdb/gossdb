package gossdb

import "errors"

var (
	Error_Null_Reply    = errors.New("null reply")
	Error_Not_Found     = errors.New("not_found")
	Error_Multi_Values  = errors.New("multi_values")
	Error_Bad_Arguments = errors.New("bad arguments")
	Error_Fail          = errors.New("fail")
)
