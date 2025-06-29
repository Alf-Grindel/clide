package user_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

// AddUser 添加用户
// params:
//   - req: 添加用户请求体
//     required: userAccount
//     optional: userAvatar, userProfile, userRole
//
// returns:
//   - userID
//   - error: nil on success, non-nil on failure
func (s *UserService) AddUser(req *user.AddUserReq) (int64, error) {
	if req == nil {
		return 0, errno.ParamErr
	}
	password, err := utils.GeneratePassword(constants.DefaultPassword)
	if err != nil {
		hlog.Errorf("user_services - AddUser: generate password failed, %s\n", err)
		return 0, errno.OperationErr
	}
	current := &db_user.User{
		UserAccount:  req.UserAccount,
		UserPassword: password,
		UserAvatar:   req.GetUserAvatar(),
		UserProfile:  req.GetUserProfile(),
		UserRole:     req.GetUserRole(),
	}

	id, err := db_user.AddUser(s.ctx, current)
	if err != nil {
		return 0, errno.OperationErr.WithMessage("添加用户失败")
	}
	return id, nil
}

// DeleteUser 删除用户
// params:
//   - req: 删除用户请求体
//     required: id
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *UserService) DeleteUser(req *user.DeleteUserReq) error {
	if req == nil {
		return errno.ParamErr
	}
	if err := db_user.DeleteUser(s.ctx, req.ID); err != nil {
		return errno.OperationErr.WithMessage("删除用户失败")
	}
	return nil
}

// UpdateUser 更新用户
// params:
//   - req: 更新用户请求体
//     required: id
//     optional: userPassword, userAvatar, userProfile, userRole
//
// returns:
//   - userVo: 脱敏用户信息
//   - error: nil on success, non-nil on failure
func (s *UserService) UpdateUser(req *user.UpdateUserReq) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	if req.UserPassword == nil && req.UserAvatar == nil && req.UserProfile == nil && req.UserRole == nil {
		return nil, errno.ParamErr.WithMessage("未有更新数据")
	}
	updates := &db_user.User{
		Id:          req.ID,
		UserAvatar:  req.GetUserAvatar(),
		UserProfile: req.GetUserProfile(),
		UserRole:    req.GetUserRole(),
		EditTime:    time.Now(),
	}
	if req.UserPassword != nil {
		password, err := utils.GeneratePassword(req.GetUserPassword())
		if err != nil {
			hlog.Errorf("user_services - UpdateUser: generate password failed, %s\n", err)
			return nil, errno.OperationErr
		}
		updates.UserPassword = password
	}
	err := db_user.UpdateUser(s.ctx, updates)
	if err != nil {
		return nil, errno.OperationErr.WithMessage("更新失败")
	}
	current, err := db_user.QueryUserById(s.ctx, req.ID)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToVo(current), nil
}

// QueryUser 分页查询用户
// params:
//   - req: 查询用户请求体
//     required: currentPage, pageSize
//     optional: id, userAccount, userProfile, userRole
//
// returns:
//   - total: total number of matched users
//   - userVos: 用户脱敏信息列表
//   - error: nil on success, non-nil on failure
func (s *UserService) QueryUser(req *user.QueryUserReq) (int64, []*base.UserVo, error) {
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
		UserRole:    req.GetUserRole(),
	}
	total, users, err := db_user.QueryUser(s.ctx, search, currentPage, pageSize)
	if err != nil {
		return 0, nil, errno.NotFoundErr
	}
	return total, ObjsToVos(users), nil
}

// GetUserById 根据id 获取用户
// params:
//   - req: 获取用户请求体
//     required: id
//
// returns:
//   - user: 未脱敏用户信息
//   - error: nil on success, non-nil on failure
func (s *UserService) GetUserById(req *user.GetUserReq) (*base.User, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	current, err := db_user.QueryUserById(s.ctx, req.ID)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToObj(current), nil
}
