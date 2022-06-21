package model

import "github.com/jinzhu/gorm"

type Course struct {
	Id   *int    `gorm:"column:ID;primaryKey" json:"id"`
	Name *string `gorm:"column:NAME" json:"name"`
}

func (c *Course) TableName() string {
	return "course_info"
}

func (c *Course) Create(db *gorm.DB) error {
	return db.Create(&c).Error
}
