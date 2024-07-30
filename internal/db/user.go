package db

import (
	"golang.org/x/crypto/bcrypt"

	"QLToolsV2/config"
	"QLToolsV2/internal/model"
)

type User struct {
	model.User
}

// GetFirstUser 获取第一个用户
// func GetFirstUser() (User, error) {
// 	var m User
// 	if err := config.GinDB.Model(&m).First(&m).Error; err != nil {
// 		return m, err
// 	}
// 	return m, nil
// }

// GetUserByUserID 用户ID 获取数据
// func GetUserByUserID(userId string) (User, error) {
// 	var m User
// 	if err := config.GinDB.Model(&m).Where("user_id = ?", userId).
// 		First(&m).Error; err != nil {
// 		return m, err
// 	}
// 	return m, nil
// }

// GetUserByUsername 用户名 获取数据
func GetUserByUsername(username string) (User, error) {
	var m User
	if err := config.GinDB.Model(&m).Where("username = ?", username).
		First(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

// Create 创建数据
func (m *User) Create() error {
	if err := config.GinDB.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

// BcryptHash 使用 bcrypt 对密码进行加密
func (m *User) BcryptHash(password string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	m.PassWord = string(bytes)
}

// BcryptCheck 对比入参密码和数据库的哈希值
func (m *User) BcryptCheck(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.PassWord), []byte(password))
	return err == nil
}

// Update 修改数据
func (m *User) Update(data map[string]any) error {
	if err := config.GinDB.Model(&m).Where("user_id = ?", m.UserID).
		Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
