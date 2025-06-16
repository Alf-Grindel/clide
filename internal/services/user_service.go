package services

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
	"time"
)

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx}
}

// UserRegister 注册
// Param: userRegisterReq
// return: user_id
func (s *UserService) UserRegister(req *user.UserRegisterReq) (int64, error) {
	if req == nil {
		return -1, errno.ParamErr
	}
	account := req.UserAccount
	_, err := db.QueryUserByAccount(s.ctx, account)
	if err == nil {
		return -1, errno.ParamErr.WithMessage("账号已存在")
	}
	password, err := utils.GeneratePassword(req.UserPassword)
	if err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr.WithMessage("注册失败")
	}
	id, err := db.CreateUser(s.ctx, &db.User{
		UserAccount:  account,
		UserPassword: password,
	})
	if err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr.WithMessage("注册失败")
	}
	return id, nil
}

// UserLogin 登录
// Param: userLoginReq
// return: user_vo 脱敏
func (s *UserService) UserLogin(req *user.UserLoginReq, c *app.RequestContext) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	account := req.UserAccount
	current, err := db.QueryUserByAccount(s.ctx, account)
	if err != nil {
		return nil, errno.ParamErr.WithMessage("用户不存在或密码错误")
	}
	password := req.UserPassword
	if !utils.ComparePassword(current.UserPassword, password) {
		return nil, errno.ParamErr.WithMessage("用户不存在或密码错误")
	}
	currentByte, err := sonic.Marshal(&current)
	if err != nil {
		hlog.Error("marshal user failed,", err)
		return nil, errno.SystemErr.WithMessage("登录失败")
	}
	session := sessions.Default(c)
	session.Set(constants.UserLoginState, currentByte)
	err = session.Save()
	if err != nil {
		hlog.Error("save session failed,", err)
		return nil, errno.SystemErr.WithMessage("登录失败")
	}
	return userObjToVo(current), nil
}

// GetLoginUser 获取登录用户
// Param: requestContext
// return: user_vo 脱敏
func (s *UserService) GetLoginUser(c *app.RequestContext) (*base.UserVo, error) {
	session := sessions.Default(c)
	currentByte, ok := session.Get(constants.UserLoginState).([]byte)
	if !ok {
		hlog.Error("can not get session value")
		return nil, errno.NotLoginErr
	}
	var current *db.User
	if err := sonic.Unmarshal(currentByte, &current); err != nil {
		hlog.Error("unmarshal user failed,", err)
		return nil, errno.NotLoginErr
	}
	return userObjToVo(current), nil
}

// UserLogout 登出
// Param: requestContext
// return:
func (s *UserService) UserLogout(c *app.RequestContext) error {
	_, err := s.GetLoginUser(c)
	if err != nil {
		return errno.NotLoginErr
	}
	session := sessions.Default(c)
	session.Delete(constants.UserLoginState)
	err = session.Save()
	if err != nil {
		hlog.Error("save session failed,", err)
		return errno.NotLoginErr
	}
	return nil
}

// EditUser 修改用户信息 [用户] - 只允许修改自己的信息
// Param: userEditReq
// return: user_vo 脱敏
func (s *UserService) EditUser(req *user.UserEditReq, c *app.RequestContext) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	login, err := s.GetLoginUser(c)
	if err != nil {
		return nil, errno.NotLoginErr
	}
	current, err := db.QueryUserById(s.ctx, login.ID)
	if err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	if req.UserPassword != nil {
		password, err := utils.GeneratePassword(req.GetUserPassword())
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr
		}
		current.UserPassword = password
	}
	current.UserAvatar = req.UserAvatar
	current.UserProfile = req.UserProfile
	current.EditTime = time.Now()
	current, err = db.UpdateUser(s.ctx, current)
	if err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr.WithMessage("更新失败")
	}
	return userObjToVo(current), nil
}

// SearchUsers 搜素用户 分页
// Param: userSearchReq
// return: list<user_vo> 脱敏
func (s *UserService) SearchUsers(req *user.UserSearchReq) (int64, []*base.UserVo, error) {
	if req == nil {
		return -1, nil, errno.ParamErr
	}
	search := &db.User{}
	if req.ID != nil {
		search.Id = req.GetID()
	}
	if req.UserAccount != nil {
		search.UserAccount = req.GetUserAccount()
	}
	search.UserProfile = req.UserProfile
	var page int64
	if req.CurrentPage != nil {
		page = req.GetCurrentPage()
		if page < 1 {
			page = constants.CurrentPage
		}
	}
	total, users, err := db.QueryUser(s.ctx, search, page)
	if err != nil {
		hlog.Error(err)
		return -1, nil, errno.SystemErr.WithMessage("查询失败")
	}
	return total, userObjsToVos(users), nil
}

// AddUser 添加用户 [管理员]
// Param: addUserReq
// return: user_id
func (s *UserService) AddUser(req *user.AddUserReq) (int64, error) {
	if req == nil {
		return -1, errno.ParamErr
	}
	password, err := utils.GeneratePassword(constants.DefaultPassword)
	if err != nil {
		hlog.Error(err)
		return -1, errno.SystemErr
	}
	role := "user"
	if req.UserRole != nil {
		role = req.GetUserRole()
	}
	current := &db.User{
		UserAccount:  req.UserAccount,
		UserPassword: password,
		UserAvatar:   req.UserAvatar,
		UserProfile:  req.UserProfile,
		UserRole:     role,
	}

	id, err := db.AddUser(s.ctx, current)
	if err != nil {
		hlog.Error(err)
		return -1, errno.OperationErr
	}
	return id, nil
}

// DeleteUser 删除用户 [管理员]
// Param: deleteUserReq
// return
func (s *UserService) DeleteUser(req *user.DeleteUserReq) error {
	if req == nil {
		return errno.ParamErr
	}
	if err := db.DeleteUser(s.ctx, req.ID); err != nil {
		hlog.Error(err)
		return errno.OperationErr
	}
	return nil
}

// UpdateUser 更新用户 [管理员]
// Param: updateUserReq
// return: user_vo 脱敏
func (s *UserService) UpdateUser(req *user.UpdateUserReq) (*base.UserVo, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	current, err := db.QueryUserById(s.ctx, req.ID)
	if err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr
	}
	if req.UserPassword != nil {
		password, err := utils.GeneratePassword(req.GetUserPassword())
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr
		}
		current.UserPassword = password
	}
	current.UserAvatar = req.UserAvatar
	current.UserProfile = req.UserProfile
	if req.UserRole != nil {
		current.UserRole = req.GetUserRole()
	}
	current.EditTime = time.Now()
	current, err = db.UpdateUser(s.ctx, current)
	if err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr.WithMessage("更新失败")
	}
	return userObjToVo(current), nil
}

// QueryUser 查询用户 [管理员]
// Param: queryUserReq
// return: list<user_vo> 脱敏
func (s *UserService) QueryUser(req *user.QueryUserReq) (int64, []*base.UserVo, error) {
	if req == nil {
		return -1, nil, errno.ParamErr
	}
	search := &db.User{}
	if req.ID != nil {
		search.Id = req.GetID()
	}
	if req.UserAccount != nil {
		search.UserAccount = req.GetUserAccount()
	}
	if req.UserRole != nil {
		search.UserRole = req.GetUserRole()
	}
	search.UserProfile = req.UserProfile
	var page int64
	if req.CurrentPage != nil {
		page = req.GetCurrentPage()
		if page < 1 {
			page = constants.CurrentPage
		}
	}
	total, users, err := db.QueryUser(s.ctx, search, page)
	if err != nil {
		hlog.Error(err)
		return -1, nil, errno.SystemErr.WithMessage("查询失败")
	}
	return total, userObjsToVos(users), nil
}

// GetUser 根据id 获取用户 [管理员]
// Param: getUserReq
// return: user 未脱敏
func (s *UserService) GetUser(req *user.GetUserReq) (*base.User, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	current, err := db.QueryUserById(s.ctx, req.ID)
	if err != nil {
		hlog.Error(err)
		return nil, errno.SystemErr.WithMessage("查询失败")
	}
	return userObjToObj(current), nil
}

func userObjToVo(current *db.User) *base.UserVo {
	if current == nil {
		return nil
	}
	avatar := ""
	if current.UserAvatar != nil {
		avatar = *current.UserAvatar
	}
	profile := ""
	if current.UserProfile != nil {
		profile = *current.UserProfile
	}
	return &base.UserVo{
		ID:          current.Id,
		UserAccount: current.UserAccount,
		UserAvatar:  avatar,
		UserProfile: profile,
		EditTime:    current.EditTime.Format(time.DateTime),
		CreateTime:  current.CreateTime.Format(time.DateTime),
	}
}

func userObjsToVos(users []*db.User) []*base.UserVo {
	if users == nil {
		return nil
	}
	var res []*base.UserVo
	for _, current := range users {
		res = append(res, userObjToVo(current))
	}
	return res
}

func userObjToObj(current *db.User) *base.User {
	if current == nil {
		return nil
	}
	isDeleteMap := map[int]string{
		0: "未删除",
		1: "已删除",
	}
	avatar := ""
	if current.UserAvatar != nil {
		avatar = *current.UserAvatar
	}
	profile := ""
	if current.UserProfile != nil {
		profile = *current.UserProfile
	}
	return &base.User{
		ID:          current.Id,
		UserAccount: current.UserAccount,
		UserAvatar:  avatar,
		UserProfile: profile,
		UserRole:    current.UserRole,
		EditTime:    current.EditTime.Format(time.DateTime),
		CreateTime:  current.CreateTime.Format(time.DateTime),
		UpdateTime:  current.UpdateTime.Format(time.DateTime),
		IsDelete:    isDeleteMap[current.IsDelete],
	}
}
