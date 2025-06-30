package user_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/internal/services"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
	"time"
)

// GetLoginUser 获取登录用户
// params:
//   - c: 请求上下文
//
// returns:
//   - userVo: 脱敏用户信息
//   - error: nil on success, non-nil on failure
func (s *UserService) GetLoginUser(c *app.RequestContext) (*base.UserVo, error) {
	loginUser, err := services.GetLoginUserIdRole(c)
	if err != nil {
		return nil, err
	}
	oldUser, err := db_user.QueryUserById(s.ctx, loginUser.Id)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToVo(oldUser), nil
}

// UserLogout 用户登出
// params:
//   - c: 请求上下文
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *UserService) UserLogout(c *app.RequestContext) error {
	session := sessions.Default(c)
	session.Delete(constants.UserLoginState)
	err := session.Save()
	if err != nil {
		hlog.Errorf("user_services - UserLogout: save session failed, %s\n", err)
		return errno.OperationErr.WithMessage("登出失败")
	}
	return nil
}

// UserEdit 修改用户信息 [用户] - 只允许修改自己的信息
// params:
//   - req: 用户编辑请求体
//     optional: userPassword, userAvatar, userProfile
//   - c: 请求上下文
//
// returns:
//   - userVo: 脱敏用户信息
//   - error: nil on success, non-nil on failure
func (s *UserService) UserEdit(req *user.UserEditReq, c *app.RequestContext) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	if req.UserPassword == nil && req.UserAvatar == nil && req.UserProfile == nil {
		return nil, errno.ParamErr.WithMessage("未更新数据")
	}
	loginUser, err := services.GetLoginUserIdRole(c)
	if err != nil {
		return nil, err
	}
	updates := &db_user.User{
		Id:          loginUser.Id,
		UserAvatar:  req.GetUserAvatar(),
		UserProfile: req.GetUserProfile(),
		EditTime:    time.Now(),
	}
	if req.UserPassword != nil {
		password, err := utils.GeneratePassword(req.GetUserPassword())
		if err != nil {
			hlog.Errorf("user_services - EditUser: generate password failed, %s\n", err)
			return nil, errno.OperationErr.WithMessage("更新失败")
		}
		updates.UserPassword = password
	}
	err = db_user.UpdateUser(s.ctx, updates)
	if err != nil {
		return nil, errno.OperationErr.WithMessage("更新失败")
	}
	oldUser, err := db_user.QueryUserById(s.ctx, loginUser.Id)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToVo(oldUser), nil
}

// UserSearch 分页搜素用户
// params:
//   - req: 用户搜索请求体
//     required: currentPage, pageSize
//     optional: id, userAccount, userProfile
//
// returns:
//   - total: total number of matched users
//   - userVos: 用户脱敏信息列表
//   - error: nil on success, non-nil on failure
func (s *UserService) UserSearch(req *user.UserSearchReq) (int64, []*base.UserVo, error) {
	if req == nil {
		return 0, nil, errno.ParamErr
	}
	currentPage := req.CurrentPage
	pageSize := req.PageSize
	if currentPage < 1 {
		currentPage = constants.CurrentPage
	}
	if pageSize < 1 || pageSize > 30 {
		pageSize = constants.PageSize
	}
	search := &db_user.User{
		Id:          req.GetID(),
		UserAccount: req.GetUserAccount(),
		UserProfile: req.GetUserProfile(),
	}
	total, oldUsers, err := db_user.QueryUser(s.ctx, search, currentPage, pageSize)
	if err != nil {
		return 0, nil, errno.NotFoundErr
	}
	return total, ObjsToVos(oldUsers), nil
}
