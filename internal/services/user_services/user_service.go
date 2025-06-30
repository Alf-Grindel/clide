package user_services

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"time"
)

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx}
}

// ObjToVo - 转换为脱敏对象
func ObjToVo(oldUser *db_user.User) *base.UserVo {
	if oldUser == nil {
		return nil
	}
	return &base.UserVo{
		ID:          oldUser.Id,
		UserAccount: oldUser.UserAccount,
		UserAvatar:  oldUser.UserAvatar,
		UserProfile: oldUser.UserProfile,
		EditTime:    oldUser.EditTime.Format(time.DateTime),
		CreateTime:  oldUser.CreateTime.Format(time.DateTime),
	}
}

// ObjsToVos - 转换为脱敏列表
func ObjsToVos(oldUsers []*db_user.User) []*base.UserVo {
	if oldUsers == nil {
		return nil
	}
	var res []*base.UserVo
	for _, oldUser := range oldUsers {
		res = append(res, ObjToVo(oldUser))
	}
	return res
}

// ObjToObj - 转化未为脱敏对象
func ObjToObj(oldUser *db_user.User) *base.User {
	if oldUser == nil {
		return nil
	}
	return &base.User{
		ID:          oldUser.Id,
		UserAccount: oldUser.UserAccount,
		UserAvatar:  oldUser.UserAvatar,
		UserProfile: oldUser.UserProfile,
		UserRole:    oldUser.UserRole,
		EditTime:    oldUser.EditTime.Format(time.DateTime),
		CreateTime:  oldUser.CreateTime.Format(time.DateTime),
		UpdateTime:  oldUser.UpdateTime.Format(time.DateTime),
		IsDelete:    constants.IsDeleteMap[oldUser.IsDelete],
	}
}
