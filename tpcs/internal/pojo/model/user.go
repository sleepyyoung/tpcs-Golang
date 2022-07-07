package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id              *int    `gorm:"column:ID;primaryKey" json:"id"`
	Username        *string `gorm:"column:USERNAME;" json:"username"`
	Password        *string `gorm:"column:PASSWORD;" json:"password"`
	Email           *string `gorm:"column:EMAIL;" json:"email"`
	Note            *string `gorm:"column:NOTE;" json:"note"`
	Status          *int    `gorm:"column:STATUS;" json:"status"`
	IsAdministrator *bool   `gorm:"column:IS_ADMINISTRATOR;" json:"isAdministrator"`
	*Model
}

func (u *User) TableName() string {
	return "user_info"
}

func (u *User) String() string {
	result := "{"
	if u.Id != nil {
		result += fmt.Sprintf("Id: %v,", *u.Id)
	} else {
		result += fmt.Sprintf("Id: %v,", "<nil>")
	}
	if u.Username != nil {
		result += fmt.Sprintf("Username: %v,", *u.Username)
	} else {
		result += fmt.Sprintf("Username: %v,", "<nil>")
	}
	if u.Password != nil {
		result += fmt.Sprintf("Password: %v,", *u.Password)
	} else {
		result += fmt.Sprintf("Password: %v,", "<nil>")
	}
	if u.Email != nil {
		result += fmt.Sprintf("Email: %v,", *u.Email)
	} else {
		result += fmt.Sprintf("Email: %v,", "<nil>")
	}
	if u.Note != nil {
		result += fmt.Sprintf("Note: %v,", *u.Note)
	} else {
		result += fmt.Sprintf("Note: %v,", "<nil>")
	}
	if u.Status != nil {
		result += fmt.Sprintf("Status: %v,", *u.Status)
	} else {
		result += fmt.Sprintf("Status: %v,", "<nil>")
	}
	if u.IsAdministrator != nil {
		result += fmt.Sprintf("IsAdministrator: %v", *u.IsAdministrator)
	} else {
		result += fmt.Sprintf("IsAdministrator: %v", "<nil>")
	}
	result += "}"

	return result
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// Update 修改用户
func (u *User) Update(db *gorm.DB) (bool, error) {
	if u.Id == nil && u.Username == nil {
		return false, nil
	}
	db = db.Model(&u)
	if u.Id != nil {
		db = db.Where("ID = ?", *u.Id)
	}
	if u.Username != nil {
		db = db.Update("USERNAME", *u.Username)
		if err := db.Error; err != nil {
			return false, err
		}
	}
	if u.Password != nil {
		db = db.Update("PASSWORD", *u.Password)
		if err := db.Error; err != nil {
			return false, err
		}
	}
	if u.Email != nil {
		db = db.Update("EMAIL", *u.Email)
		if err := db.Error; err != nil {
			return false, err
		}
	}
	if u.Status != nil {
		db = db.Update("STATUS", *u.Status)
		if err := db.Error; err != nil {
			return false, err
		}
	}

	return true, nil
}

// Create 创建用户
func (u *User) Create(db *gorm.DB) error {
	return db.Create(&u).Error
}
