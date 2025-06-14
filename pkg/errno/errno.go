package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode      = 0
	ParamErrCode     = 40000
	NotLoginErrCode  = 40100
	NoAuthErrCode    = 40101
	SystemErrCode    = 50000
	OperationErrCode = 50001
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code - %d, err_msg - %v", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		code,
		msg,
	}
}

var (
	Success      = NewErrNo(SuccessCode, "ok")
	SystemErr    = NewErrNo(SystemErrCode, "系统内部错误")
	ParamErr     = NewErrNo(ParamErrCode, "参数错误")
	NotLoginErr  = NewErrNo(NotLoginErrCode, "用户未登录")
	NoAuthErr    = NewErrNo(NoAuthErrCode, "无权限")
	OperationErr = NewErrNo(OperationErrCode, "操作失败")
)

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

func ConvertErr(e error) ErrNo {
	err := ErrNo{}
	if errors.As(e, &err) {
		return err
	}
	s := SystemErr
	s.ErrMsg = e.Error()
	return s
}
