package db_user

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

type User struct {
	Id           int64     `json:"id"`
	UserAccount  string    `json:"user_account"`
	UserPassword string    `json:"user_password"`
	UserAvatar   string    `json:"user_avatar"`
	UserProfile  string    `json:"user_profile"`
	UserRole     string    `json:"user_role"`
	EditTime     time.Time `json:"edit_time"`
	CreateTime   time.Time `json:"create_time" gorm:"<-:false"`
	UpdateTime   time.Time `json:"update_time" gorm:"<-:false"`
	IsDelete     int       `json:"is_delete"`
}

func (u User) TableName() string {
	return constants.UserTableName
}

// CreateUser - create user when user register
// params:
//   - required: userAccount, userPassword
//
// returns:
//   - userID
//   - error: nil on success, non-nil on failure
func CreateUser(ctx context.Context, userAccount, userPassword string) (int64, error) {
	id, err := utils.GenerateId()
	if err != nil {
		hlog.Errorf("dal - CreateUser: generate user id failed, %s\n", err)
		return 0, err
	}
	user := &User{
		Id:           id,
		UserAccount:  userAccount,
		UserPassword: userPassword,
	}
	res := db.DB.WithContext(ctx).Select("id", "user_account", "user_password").Create(user)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - CreateUser: insert user into db failed, %s\n", err)
		return 0, err
	}
	return user.Id, nil
}

// AddUser - create user by admin
// params:
//   - required: userAccount, userPassword
//   - optional: userAvatar, userProfile, userRole
//
// returns:
//   - userID
//   - error: nil on success, non-nil on failure
func AddUser(ctx context.Context, user *User) (int64, error) {
	id, err := utils.GenerateId()
	if err != nil {
		hlog.Errorf("dal - AddUser: generate user id failed, %s\n", err)
		return 0, err
	}
	user.Id = id
	omitFields := []string{"edit_time", "is_delete"}
	if user.UserAvatar == "" {
		omitFields = append(omitFields, "user_avatar")
	}
	if user.UserProfile == "" {
		omitFields = append(omitFields, "user_profile")
	}
	if user.UserRole == "" {
		omitFields = append(omitFields, "user_role")
	}
	res := db.DB.WithContext(ctx).Omit(omitFields...).Create(user)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - AddUser: insert user into db failed, %s\n", err)
		return 0, err
	}
	return user.Id, nil
}

// UpdateUser - update user
// params:
//   - required: id
//   - optional: userPassword, userAvatar, userProfile, userRole, editTime
//
// returns:
//   - error: nil on success, non-nil on failure

func UpdateUser(ctx context.Context, user *User) error {
	res := db.DB.Model(&User{}).WithContext(ctx).Where("id = ? and is_delete = 0", user.Id).Updates(user)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - UpdateUser: updates user failed, %s\n", err)
		return err
	}
	return nil
}

// DeleteUser - delete user
// params:
//   - required: id
//
// returns:
//   - error: nil on success, non-nil on failure
func DeleteUser(ctx context.Context, id int64) error {
	res := db.DB.Model(&User{}).WithContext(ctx).Where("id = ? and is_delete = 0", id).Update("is_delete", 1)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - DeleteUser: delete user failed, %s\n", err)
		return err
	}
	return nil
}

// QueryUserById - query user based on given user id
// params:
//   - required: id
//
// returns:
//   - user
//   - error: nil on success, non-nil on failure
func QueryUserById(ctx context.Context, id int64) (*User, error) {
	user := &User{}
	res := db.DB.WithContext(ctx).Where("id = ? and is_delete = 0", id).First(&user)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - QueryUserById: query user failed, %s\n", err)
		return nil, err
	}
	return user, nil
}

// QueryUserByAccount - query user based on given user account
// params:
//   - required: userAccount
//
// returns:
//   - user
//   - error: nil on success, non-nil on failure
func QueryUserByAccount(ctx context.Context, account string) (*User, error) {
	user := &User{}
	res := db.DB.WithContext(ctx).Where("user_account = ? and is_delete = 0", account).First(&user)
	if err := res.Error; err != nil {
		hlog.Errorf("dal - QueryUserByAccount: query user failed, %s\n", err)
		return nil, err
	}
	return user, nil
}

// QueryUser - query users based on the given filters
// params:
//   - required: currentPage, pageSize
//   - optional: id, userAccount, userProfile, userRole
//
// returns:
//   - total: total number of matched users
//   - users: list of users matching the criteria
//   - error: nil on success, non-nil on failure
func QueryUser(ctx context.Context, user *User, currentPage, pageSize int64) (int64, []*User, error) {
	var users []*User
	res := db.DB.WithContext(ctx).Model(&User{}).Where("is_delete = 0")
	if user.Id != 0 {
		res = res.Where("id = ?", user.Id)
	}
	if user.UserAccount != "" {
		res = res.Where("user_account = ?", user.UserAccount)
	}
	if user.UserProfile != "" {
		res = res.Where("user_profile like ?", "%"+user.UserProfile+"%")
	}
	if user.UserRole != "" {
		res = res.Where("user_role = ?", user.UserRole)
	}

	var total int64
	if err := res.Count(&total).Error; err != nil {
		hlog.Errorf("dal - QueryUser: count match user failed, %s\n", err)
		return -1, nil, err
	}

	offset := (currentPage - 1) * pageSize
	if err := res.Offset(int(offset)).Limit(int(pageSize)).Find(&users).Error; err != nil {
		hlog.Errorf("dal - QueryUser: Query user failed, %s\n", err)
		return -1, nil, err
	}
	return total, users, nil
}
