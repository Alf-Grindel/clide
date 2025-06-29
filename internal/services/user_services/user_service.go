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
func ObjToVo(current *db_user.User) *base.UserVo {
	if current == nil {
		return nil
	}
	return &base.UserVo{
		ID:          current.Id,
		UserAccount: current.UserAccount,
		UserAvatar:  current.UserAvatar,
		UserProfile: current.UserProfile,
		EditTime:    current.EditTime.Format(time.DateTime),
		CreateTime:  current.CreateTime.Format(time.DateTime),
	}
}

// ObjsToVos - 转换为脱敏列表
func ObjsToVos(currents []*db_user.User) []*base.UserVo {
	if currents == nil {
		return nil
	}
	var res []*base.UserVo
	for _, current := range currents {
		res = append(res, ObjToVo(current))
	}
	return res
}

// ObjToObj - 转化未为脱敏对象
func ObjToObj(current *db_user.User) *base.User {
	if current == nil {
		return nil
	}
	return &base.User{
		ID:          current.Id,
		UserAccount: current.UserAccount,
		UserAvatar:  current.UserAvatar,
		UserProfile: current.UserProfile,
		UserRole:    current.UserRole,
		EditTime:    current.EditTime.Format(time.DateTime),
		CreateTime:  current.CreateTime.Format(time.DateTime),
		UpdateTime:  current.UpdateTime.Format(time.DateTime),
		IsDelete:    constants.IsDeleteMap[current.IsDelete],
	}
}
