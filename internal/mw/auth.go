package mw

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		currentByte, ok := session.Get(constants.UserLoginState).([]byte)
		if !ok {
			resp := errno.BuildBaseResp(errno.NotLoginErr)
			c.JSON(200, resp)
			c.Abort()
			return
		}
		var user *db.User
		if err := sonic.Unmarshal(currentByte, &user); err != nil {
			resp := errno.BuildBaseResp(errno.NotLoginErr)
			c.JSON(200, resp)
			c.Abort()
			return
		}
		if user.UserRole != "admin" {
			resp := errno.BuildBaseResp(errno.NoAuthErr)
			c.JSON(200, resp)
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
