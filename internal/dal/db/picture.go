package db

import (
	"context"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"time"
)

type Picture struct {
	Id           int64     `json:"id"`
	Url          string    `json:"url"`
	PicName      string    `json:"pic_name"`
	Introduction *string   `json:"introduction"`
	Category     *string   `json:"category"`
	Tags         *string   `json:"tags"`
	PicSize      int64     `json:"pic_size"`
	PicWidth     int32     `json:"pic_width"`
	PicHeight    int32     `json:"pic_height"`
	PicScale     float64   `json:"pic_scale"`
	PicFormat    string    `json:"pic_format"`
	UserId       int64     `json:"user_id"`
	EditTime     time.Time `json:"edit_time"`
	CreateTime   time.Time `json:"create_time" gorm:"<-:false"`
	UpdateTime   time.Time `json:"update_time" gorm:"<-:false"`
	IsDelete     int       `json:"is_delete"`
}

func (p Picture) TableName() string {
	return constants.PictureTableName
}

func CreatePicture(ctx context.Context, picture *Picture) (*Picture, error) {
	if picture.Id == 0 {
		id, err := utils.GenerateId()
		if err != nil {
			return nil, err
		}
		picture.Id = id
	}
	res := DB.WithContext(ctx).Omit("edit_time", "is_delete").Create(&picture)
	if err := res.Error; err != nil {
		return nil, err
	}
	return QueryPictureById(ctx, picture.Id)
}

func UpdatePicture(ctx context.Context, picture *Picture) (*Picture, error) {
	res := DB.WithContext(ctx).Model(&Picture{}).Where("id = ? and is_delete = 0", picture.Id).Updates(&picture)
	if err := res.Error; err != nil {
		return nil, err
	}
	return QueryPictureById(ctx, picture.Id)
}

func QueryPictureById(ctx context.Context, id int64) (*Picture, error) {
	picture := &Picture{}
	res := DB.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&picture)
	if err := res.Error; err != nil {
		return nil, err
	}
	return picture, nil
}
