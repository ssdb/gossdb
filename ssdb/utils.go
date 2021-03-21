package ssdb

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	notFound = "not_found"
	oK       = "ok"

	// public
	NotFound = notFound
	Ok       = oK
)

//生成错误信息，已经确定是有错误
func makeError(resp []string, errKey ...interface{}) error {
	if len(resp) < 1 {
		return errors.New("ssdb response error")
	}
	//正常返回的不存在不报错，如果要捕捉这个问题请使用exists
	if resp[0] == notFound {
		return nil
	}
	if len(errKey) > 0 {
		return fmt.Errorf("access ssdb error, code is %v, parameter is %v", resp, errKey)
	}
	return fmt.Errorf("access ssdb error, code is %v", resp)
}
