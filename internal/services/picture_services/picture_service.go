package picture_services

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/services/user_services"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

type PictureService struct {
	ctx context.Context
}

func NewPictureService(ctx context.Context) *PictureService {
	return &PictureService{
		ctx: ctx,
	}
}

// ObjToVo - 转化为脱敏对象
func ObjToVo(current *db_picture.Picture, user *db_user.User) *base.PictureVo {
	if current == nil || user == nil {
		return nil
	}

	var tagsList []string
	if current.Tags != "" {
		if err := sonic.Unmarshal([]byte(current.Tags), &tagsList); err != nil {
			hlog.Errorf("picture_services - ObjToVo: unmarshal tags failed, %s\n", err)
			return nil
		}
	}

	currentUser := user_services.ObjToVo(user)

	return &base.PictureVo{
		ID:           current.Id,
		URL:          current.Url,
		PicName:      current.PicName,
		Introduction: current.Introduction,
		Category:     current.Category,
		Tags:         tagsList,
		PicSize:      current.PicSize,
		PicWidth:     current.PicWidth,
		PicHeight:    current.PicHeight,
		PicScale:     current.PicScale,
		PicFormat:    current.PicFormat,
		EditTime:     current.EditTime.Format(time.DateTime),
		CreateTime:   current.CreateTime.Format(time.DateTime),
		UserId:       current.UserId,
		User:         currentUser,
	}
}

// ObjsToVos - 转化为脱敏列表
func ObjsToVos(ctx context.Context, currents []*db_picture.Picture) []*base.PictureVo {
	if currents == nil {
		return nil
	}
	var pictures []*base.PictureVo
	for _, current := range currents {
		user, err := db_user.QueryUserById(ctx, current.UserId)
		if err != nil {
			return nil
		}
		pictures = append(pictures, ObjToVo(current, user))
	}
	return pictures
}

// ObjToObj - 转化为未脱敏列表
func ObjToObj(current *db_picture.Picture, user *db_user.User) *base.Picture {
	if current == nil {
		return nil
	}

	var tagsList []string
	if current.Tags != "" {
		if err := sonic.Unmarshal([]byte(current.Tags), &tagsList); err != nil {
			hlog.Errorf("picture_services - ObjToVo: unmarshal tags failed, %s\n", err)
			return nil
		}
	}

	currentUser := user_services.ObjToObj(user)

	return &base.Picture{
		ID:           current.Id,
		URL:          current.Url,
		PicName:      current.PicName,
		Introduction: current.Introduction,
		Category:     current.Category,
		Tags:         tagsList,
		PicSize:      current.PicSize,
		PicWidth:     current.PicWidth,
		PicHeight:    current.PicHeight,
		PicScale:     current.PicScale,
		PicFormat:    current.PicFormat,
		EditTime:     current.EditTime.Format(time.DateTime),
		CreateTime:   current.CreateTime.Format(time.DateTime),
		UpdateTime:   current.UpdateTime.Format(time.DateTime),
		IsDelete:     constants.IsDeleteMap[current.IsDelete],
		UserId:       current.UserId,
		User:         currentUser,
	}
}
