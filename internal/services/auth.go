package services

import (
	"github.com/Alf-Grindel/clide/internal/model"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

// GetLoginUserIdRole - 从请求上下文中获取login userid userRole
// params:
//   - c: 请求上下文
//
// returns:
//   - loginUser: model.LoginUser
//   - error: nil on success, non-nil on failure
func GetLoginUserIdRole(c *app.RequestContext) (*model.LoginUser, error) {
	userIdVal, exits := c.Get("user_id")
	userId, ok := userIdVal.(int64)
	if !exits || !ok {
		return nil, errno.NotLoginErr
	}
	userRoleVal, exits := c.Get("user_role")
	userRole, ok := userRoleVal.(string)
	if !exits || !ok {
		return nil, errno.NotLoginErr
	}
	return &model.LoginUser{
		Id:   userId,
		Role: userRole,
	}, nil
}
