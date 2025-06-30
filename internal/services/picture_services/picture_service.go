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
func ObjToVo(oldPicture *db_picture.Picture, user *db_user.User) *base.PictureVo {
	if oldPicture == nil || user == nil {
		return nil
	}

	var tagsList []string
	if oldPicture.Tags != "" {
		if err := sonic.Unmarshal([]byte(oldPicture.Tags), &tagsList); err != nil {
			hlog.Errorf("picture_services - ObjToVo: unmarshal tags failed, %s\n", err)
			return nil
		}
	}

	currentUser := user_services.ObjToVo(user)

	return &base.PictureVo{
		ID:           oldPicture.Id,
		URL:          oldPicture.Url,
		PicName:      oldPicture.PicName,
		Introduction: oldPicture.Introduction,
		Category:     oldPicture.Category,
		Tags:         tagsList,
		PicSize:      oldPicture.PicSize,
		PicWidth:     oldPicture.PicWidth,
		PicHeight:    oldPicture.PicHeight,
		PicScale:     oldPicture.PicScale,
		PicFormat:    oldPicture.PicFormat,
		EditTime:     oldPicture.EditTime.Format(time.DateTime),
		CreateTime:   oldPicture.CreateTime.Format(time.DateTime),
		UserId:       oldPicture.UserId,
		User:         currentUser,
	}
}

// ObjsToVos - 转化为脱敏列表
func ObjsToVos(ctx context.Context, oldPictures []*db_picture.Picture) []*base.PictureVo {
	if oldPictures == nil {
		return nil
	}
	var pictures []*base.PictureVo
	for _, oldPicture := range oldPictures {
		user, err := db_user.QueryUserById(ctx, oldPicture.UserId)
		if err != nil {
			return nil
		}
		pictures = append(pictures, ObjToVo(oldPicture, user))
	}
	return pictures
}

// ObjToObj - 转化为未脱敏对象
func ObjToObj(oldPicture *db_picture.Picture, oldUser *db_user.User) *base.Picture {
	if oldPicture == nil {
		return nil
	}

	var tagsList []string
	if oldPicture.Tags != "" {
		if err := sonic.Unmarshal([]byte(oldPicture.Tags), &tagsList); err != nil {
			hlog.Errorf("picture_services - ObjToVo: unmarshal tags failed, %s\n", err)
			return nil
		}
	}

	user := user_services.ObjToObj(oldUser)

	return &base.Picture{
		ID:            oldPicture.Id,
		URL:           oldPicture.Url,
		PicName:       oldPicture.PicName,
		Introduction:  oldPicture.Introduction,
		Category:      oldPicture.Category,
		Tags:          tagsList,
		PicSize:       oldPicture.PicSize,
		PicWidth:      oldPicture.PicWidth,
		PicHeight:     oldPicture.PicHeight,
		PicScale:      oldPicture.PicScale,
		PicFormat:     oldPicture.PicFormat,
		EditTime:      oldPicture.EditTime.Format(time.DateTime),
		CreateTime:    oldPicture.CreateTime.Format(time.DateTime),
		UpdateTime:    oldPicture.UpdateTime.Format(time.DateTime),
		IsDelete:      constants.IsDeleteMap[oldPicture.IsDelete],
		UserId:        oldPicture.UserId,
		User:          user,
		ReviewStatus:  constants.ReviewStatusMap[oldPicture.ReviewStatus],
		ReviewMessage: oldPicture.ReviewMessage,
		ReviewId:      oldPicture.ReviewId,
		ReviewTime:    oldPicture.ReviewTime.Format(time.DateTime),
	}
}

// ObjsToObjs - 转化为未脱敏列表
func ObjsToObjs(ctx context.Context, oldPictures []*db_picture.Picture) []*base.Picture {
	if oldPictures == nil {
		return nil
	}
	var pictures []*base.Picture
	for _, oldPicture := range oldPictures {
		user, err := db_user.QueryUserById(ctx, oldPicture.UserId)
		if err != nil {
			return nil
		}
		pictures = append(pictures, ObjToObj(oldPicture, user))
	}
	return pictures
}
