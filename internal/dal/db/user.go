package db

import (
	"context"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"time"
)

type User struct {
	Id           int64     `json:"id"`
	UserAccount  string    `json:"user_account"`
	UserPassword string    `json:"user_password"`
	UserAvatar   *string   `json:"user_avatar"`
	UserProfile  *string   `json:"user_profile"`
	UserRole     string    `json:"user_role"`
	EditTime     time.Time `json:"edit_time"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	IsDelete     int       `json:"is_delete"`
}

func (u User) TableName() string {
	return constants.UserTableName
}

func CreateUser(ctx context.Context, user *User) (int64, error) {
	id, err := utils.GenerateId()
	if err != nil {
		return -1, err
	}
	user.Id = id
	res := DB.WithContext(ctx).Select("id", "user_account", "user_password").Create(user)
	if err := res.Error; err != nil {
		return -1, err
	}
	return user.Id, nil
}

func QueryUserById(ctx context.Context, id int64) (*User, error) {
	user := &User{}
	res := DB.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&user)
	if err := res.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func QueryUserByAccount(ctx context.Context, account string) (*User, error) {
	user := &User{}
	res := DB.WithContext(ctx).Where("user_account = ? and is_delete = 0", account).First(&user)
	if err := res.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func QueryUser(ctx context.Context, user *User, page int64) (int64, []*User, error) {
	var users []*User
	res := DB.WithContext(ctx).Model(&User{}).Where("is_delete = 0")
	if user.Id > 0 {
		res = res.Where("id = ?", user.Id)
	}
	if user.UserAccount != "" {
		res = res.Where("user_account = ?", user.UserAccount)
	}
	if user.UserProfile != nil {
		res = res.Where("user_profile like ?", "%"+*user.UserProfile+"%")
	}
	if user.UserRole != "" {
		res = res.Where("user_role = ?", user.UserRole)
	}

	var total int64
	if err := res.Count(&total).Error; err != nil {
		return -1, nil, err
	}

	offset := (page - 1) * constants.PageSize
	if err := res.Offset(int(offset)).Limit(constants.PageSize).Find(&users).Error; err != nil {
		return -1, nil, err
	}
	return total, users, nil
}

func DeleteUser(ctx context.Context, id int64) error {
	user := &User{}
	res := DB.Model(&user).WithContext(ctx).Where("id = ? and is_delete = 0", id).Update("is_delete", 1)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context, user *User) (*User, error) {
	res := DB.WithContext(ctx).Updates(user)
	if err := res.Error; err != nil {
		return nil, err
	}
	res = DB.WithContext(ctx).Where("id = ? and is_delete = 0", user.Id)
	current := &User{}
	if err := res.First(&current).Error; err != nil {
		return nil, err
	}
	return current, nil
}

func AddUser(ctx context.Context, user *User) (int64, error) {
	id, err := utils.GenerateId()
	if err != nil {
		return -1, err
	}
	user.Id = id
	res := DB.WithContext(ctx).Select("id", "user_account", "user_password", "user_avatar", "user_profile", "user_role").Create(user)
	if err := res.Error; err != nil {
		return -1, err
	}
	return user.Id, nil
}
