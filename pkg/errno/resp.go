package errno

import (
	"github.com/Alf-Grindel/clide/internal/model/base"
)

func BuildBaseResp(err error) *base.BaseResp {
	e := ConvertErr(err)

	resp := &base.BaseResp{
		Code: e.ErrCode,
		Msg:  e.ErrMsg,
	}
	return resp
}
