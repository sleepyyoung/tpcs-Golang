package model

type Course struct {
	Id   *int    `gorm:"column:ID;primaryKey" json:"id"`
	Name *string `gorm:"column:NAME" json:"name"`
}

func (c *Course) TableName() string {
	return "course_info"
}
