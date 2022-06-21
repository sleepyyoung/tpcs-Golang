package model

import "fmt"

// CombinePlan 组卷方案
type CombinePlan struct {
	Id         *int    `gorm:"column:ID;primaryKey" json:"id"`
	User       *User   `gorm:"column:USER_ID;foreignKey:UserId;references:USER_ID" json:"user"`
	UserId     *int    `gorm:"column:USER_ID"`
	Course     *Course `gorm:"column:COURSE_ID;foreignKey:CourseId;references:COURSE_ID" json:"course"`
	CourseId   *int    `gorm:"column:COURSE_ID"`
	PaperTitle *string `gorm:"column:PAPER_TITLE" json:"paperTitle"`
	Plan       *string `gorm:"column:PLAN" json:"plan"`
	Score      *int    `gorm:"column:SCORE" json:"score"`
	Note       *string `gorm:"column:NOTE" json:"note"`
}

// CombinePlan4Add 添加组卷方案的结构体
type CombinePlan4Add struct {
	UserId     *int    `gorm:"column:USER_ID"`
	CourseId   *int    `gorm:"column:COURSE_ID"`
	PaperTitle *string `gorm:"column:PAPER_TITLE" json:"paperTitle"`
	Plan       *string `gorm:"column:PLAN" json:"plan"`
	Score      *int    `gorm:"column:SCORE" json:"score"`
	Note       *string `gorm:"column:NOTE" json:"note"`
}

// CombinePlan4Edit 编辑组卷方案的结构体
type CombinePlan4Edit struct {
	Id         *int    `gorm:"column:ID;primaryKey" json:"id"`
	UserId     *int    `gorm:"column:USER_ID"`
	CourseId   *int    `gorm:"column:COURSE_ID"`
	PaperTitle *string `gorm:"column:PAPER_TITLE" json:"paperTitle"`
	Plan       *string `gorm:"column:PLAN" json:"plan"`
	Score      *int    `gorm:"column:SCORE" json:"score"`
	Note       *string `gorm:"column:NOTE" json:"note"`
}

func (cp *CombinePlan) TableName() string {
	return "question_combine_plan_info"
}

func (cp4a *CombinePlan4Add) TableName() string {
	return "question_combine_plan_info"
}

func (cp4e *CombinePlan4Edit) TableName() string {
	return "question_combine_plan_info"
}

func (cp *CombinePlan) String() string {
	result := "{"
	if cp.Id != nil {
		result += fmt.Sprintf("Id: %v,", *cp.Id)
	} else {
		result += fmt.Sprintf("Id: %v,", "<nil>")
	}
	if cp.UserId != nil {
		result += fmt.Sprintf("UserId: %v,", *cp.UserId)
	} else {
		result += fmt.Sprintf("UserId: %v,", "<nil>")
	}
	if cp.User != nil {
		result += fmt.Sprintf("User: %v,", *cp.User)
	} else {
		result += fmt.Sprintf("User: %v,", "<nil>")
	}
	if cp.CourseId != nil {
		result += fmt.Sprintf("CourseId: %v,", *cp.CourseId)
	} else {
		result += fmt.Sprintf("CourseId: %v,", "<nil>")
	}
	if cp.Course != nil {
		result += fmt.Sprintf("Course: %v,", *cp.Course)
	} else {
		result += fmt.Sprintf("Course: %v,", "<nil>")
	}
	if cp.PaperTitle != nil {
		result += fmt.Sprintf("PaperTitle: %v,", *cp.PaperTitle)
	} else {
		result += fmt.Sprintf("PaperTitle: %v,", "<nil>")
	}
	if cp.Plan != nil {
		result += fmt.Sprintf("Plan: %v,", *cp.Plan)
	} else {
		result += fmt.Sprintf("Plan: %v,", "<nil>")
	}
	if cp.Note != nil {
		result += fmt.Sprintf("Note: %v,", *cp.Note)
	} else {
		result += fmt.Sprintf("Note: %v,", "<nil>")
	}
	result += "}"

	return result
}
