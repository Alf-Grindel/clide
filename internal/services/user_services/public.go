package user_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

// UserRegister 注册
// params:
//   - req: 用户注册请求体，
//     required: userAccount, userPassword
//
// returns:
//   - userID
//   - error: nil on success, non-nil on failure
func (s *UserService) UserRegister(req *user.UserRegisterReq) (int64, error) {
	if req == nil {
		return 0, errno.ParamErr
	}
	account := req.UserAccount
	_, err := db_user.QueryUserByAccount(s.ctx, account)
	if err == nil {
		return 0, errno.ParamErr.WithMessage("账号已存在")
	}
	password, err := utils.GeneratePassword(req.UserPassword)
	if err != nil {
		hlog.Errorf("user_services - UserRegister: generate password failed, %s\n", err)
		return 0, errno.OperationErr.WithMessage("注册失败")
	}
	id, err := db_user.CreateUser(s.ctx, account, password)
	if err != nil {
		return 0, errno.SystemErr.WithMessage("注册失败")
	}
	return id, nil
}

// UserLogin 登录
// params:
//   - req: 用户登录请求体
//     required: userAccount, userPassword
//   - c: 请求上下文
//
// returns:
//   - userVo: 已脱敏的用户信息
//   - error: nil on success, non-nil on failure
func (s *UserService) UserLogin(req *user.UserLoginReq, c *app.RequestContext) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	account := req.UserAccount
	oldUser, err := db_user.QueryUserByAccount(s.ctx, account)
	if err != nil {
		return nil, errno.ParamErr.WithMessage("用户不存在或密码错误")
	}
	password := req.UserPassword
	if !utils.ComparePassword(oldUser.UserPassword, password) {
		return nil, errno.ParamErr.WithMessage("用户不存在或密码错误")
	}
	oldUserByte, err := sonic.Marshal(&oldUser)
	if err != nil {
		hlog.Errorf("user_services - UserLogin: marshal user failed, %s\n", err)
		return nil, errno.OperationErr.WithMessage("登录失败")
	}
	session := sessions.Default(c)
	session.Set(constants.UserLoginState, oldUserByte)
	err = session.Save()
	if err != nil {
		hlog.Errorf("user_services - UserLogin: save session failed, %s\n", err)
		return nil, errno.OperationErr.WithMessage("登录失败")
	}
	return ObjToVo(oldUser), nil
}
