package db_picture

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

type Picture struct {
	Id            int64     `json:"id"`
	Url           string    `json:"url"`
	PicName       string    `json:"pic_name"`
	Introduction  string    `json:"introduction"`
	Category      string    `json:"category"`
	Tags          string    `json:"tags"`
	PicSize       int64     `json:"pic_size"`
	PicWidth      int32     `json:"pic_width"`
	PicHeight     int32     `json:"pic_height"`
	PicScale      float64   `json:"pic_scale"`
	PicFormat     string    `json:"pic_format"`
	UserId        int64     `json:"user_id"`
	EditTime      time.Time `json:"edit_time"`
	CreateTime    time.Time `json:"create_time" gorm:"<-:false"`
	UpdateTime    time.Time `json:"update_time" gorm:"<-:false"`
	IsDelete      int       `json:"is_delete"`
	ReviewStatus  int       `json:"review_status"`
	ReviewMessage string    `json:"review_message"`
	ReviewId      int64     `json:"review_id"`
	ReviewTime    time.Time `json:"review_time"`
}

func (p Picture) TableName() string {
	return constants.PictureTableName
}

// CreatePicture - create picture
// params:
//   - picture:
//     required: url, picName, picSize, picWidth, picHeight, picScale, picFormat, userId
//     optional: introduction, category, tags
//
// returns:
//   - pictureId
//   - error: nil on success, non-nil on failure
func CreatePicture(ctx context.Context, picture *Picture) (int64, error) {
	id, err := utils.GenerateId()
	if err != nil {
		hlog.Errorf("dal - CreatePicture: generate picture id failed, %s\n", err)
		return 0, err
	}
	picture.Id = id
	omitFields := []string{"edit_time", "is_delete"}
	if picture.Introduction == "" {
		omitFields = append(omitFields, "introduction")
	}
	if picture.Category == "" {
		omitFields = append(omitFields, "category")
	}
	if picture.Tags == "" {
		omitFields = append(omitFields, "tags")
	}
	if picture.ReviewMessage == "" {
		omitFields = append(omitFields, "review_message")
	}
	if picture.ReviewId == 0 {
		omitFields = append(omitFields, "review_id")
	}
	if picture.ReviewTime.IsZero() {
		omitFields = append(omitFields, "review_time")
	}
	res := db.DB.WithContext(ctx).Omit(omitFields...).Create(&picture)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - CreatePicture: create picture into db failed, %s\n", err)
		return 0, err
	}
	return id, nil
}

// DeletePicture -  delete picture based on given id
// params:
//   - pictureId
//
// returns:
//   - error: nil on success, non-nil on failure
func DeletePicture(ctx context.Context, id int64) error {
	res := db.DB.Model(&Picture{}).WithContext(ctx).Where("id = ? and is_delete = 0", id).Update("is_delete", 1)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - DeletePicture: delete picture failed, %s\n", err)
		return err
	}
	return nil
}

// UpdatePicture - update picture
// params:
//   - picture
//     required: pictureId
//     optional: url, picName, picSize, picWidth, picHeight, picScale, picFormat, userId, introduction, category, tags
//
// returns:
//   - error: nil on success, non-nil on failure
func UpdatePicture(ctx context.Context, picture *Picture) error {
	res := db.DB.WithContext(ctx).Model(&Picture{}).Where("id = ? and is_delete = 0", picture.Id).Updates(&picture)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - UpdatePicture: update picture failed, %s\n", err)
		return err
	}
	return nil
}

// QueryPictureById - query picture based on given id
// params:
//   - pictureId
//
// returns:
//   - picture
//   - error: nil on success, non-nil on failure
func QueryPictureById(ctx context.Context, id int64) (*Picture, error) {
	picture := &Picture{}
	res := db.DB.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&picture)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - QueryPictureById: query picture by id failed, %s\n", err)
		return nil, err
	}
	return picture, nil
}

// QueryPicture - query picture based on given filter
// params:
//   - picture
//     required: reviewStatus
//     optional: id, picName, introduction, category, picSize, picWidth, PicHeight, picScale, picFormat, userId,
//     optional: reviewMessage, reviewId
//   - searchText: match picName or introduction (optional)
//   - tags: tags list (must all match) optional
//   - currentPage (required)
//   - pageSize (required)
//
// returns:
//   - total: total number of matched picture
//   - pictures: list of picture matching the criteria
//   - error: nil on success, non-nil on failure
func QueryPicture(ctx context.Context, picture *Picture, searchText string, tags []string, currentPage, pageSize int64) (int64, []*Picture, error) {
	var pictures []*Picture
	res := db.DB.WithContext(ctx).Model(&Picture{}).Where("is_delete = 0 ")
	if picture.Id != 0 {
		res = res.Where("id = ?", picture.Id)
	}
	if picture.PicSize != 0 {
		res = res.Where("pic_size = ?", picture.PicSize)
	}
	if picture.PicWidth != 0 {
		res = res.Where("pic_width = ?", picture.PicWidth)
	}
	if picture.PicHeight != 0 {
		res = res.Where("pic_height = ?", picture.PicHeight)
	}
	if picture.PicScale != 0.0 {
		res = res.Where("pic_scale = ?", picture.PicScale)
	}
	if picture.UserId != 0 {
		res = res.Where("user_id = ?", picture.UserId)
	}
	if picture.ReviewId != 0 {
		res = res.Where("review_id = ?", picture.ReviewId)
	}
	if picture.ReviewStatus != -1 {
		res = res.Where("review_status = ?", picture.ReviewStatus)
	}

	if picture.PicName != "" {
		res = res.Where("pic_name like ?", "%"+picture.PicName+"%")
	}
	if picture.Introduction != "" {
		res = res.Where("introduction like ? ", "%"+picture.Introduction+"%")
	}
	if picture.Category != "" {
		res = res.Where("category like ?", "%"+picture.Category+"%")
	}
	if picture.PicFormat != "" {
		res = res.Where("pic_format like ?", "%"+picture.PicFormat+"%")
	}
	if picture.ReviewMessage != "" {
		res = res.Where("review_message like ?", "%"+picture.ReviewMessage+"%")
	}
	if searchText != "" {
		res = res.Where("pic_name like ? or introduction like ?", "%"+searchText+"%", "%"+searchText+"%")
	}

	if tags != nil {
		for _, tag := range tags {
			res = res.Where(" tags like ? ", "%\""+tag+"\"%")
		}
	}

	var total int64
	if err := res.Count(&total).Error; err != nil {
		hlog.Errorf("dal - QueryPicture: count match picture failed, %s\n", err)
		return 0, nil, err
	}

	offset := (currentPage - 1) * pageSize
	if err := res.Offset(int(offset)).Limit(int(pageSize)).Find(&pictures).Error; err != nil {
		hlog.Errorf("dal - QueryPicture: query picture failed, %s\n", err)
		return 0, nil, err
	}
	return total, pictures, nil
}
