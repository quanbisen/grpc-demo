package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"user/internal/service"
)

type User struct {
	UserId   uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Nickname string
	Password string
}

const (
	PasswordCost = 12 // 密码加密难度
)

func (user *User) CheckUserExist(req *service.UserRequest) bool {
	if err := DB.Where("username=?", req.Username).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

// ShowUserInfo 获取用户信息
func (user *User) ShowUserInfo(req *service.UserRequest) error {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}
	return errors.New("username not exist")
}

// UserCreate 创建用户
func (user *User) UserCreate(req *service.UserRequest) error {
	var count int64
	DB.Where("username=?", req.Username).Count(&count)
	if count != 0 {
		return errors.New("username exist")
	}
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.SetPassword(req.Password)
	err := DB.Create(user).Error
	return err
}

func (user *User) SetPassword(str string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), PasswordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(str))
	return err == nil
}

func BuildUser(item User) *service.UserModel {
	return &service.UserModel{
		UserId:   uint64(item.UserId),
		Username: item.Username,
		Nickname: item.Nickname,
	}
}
