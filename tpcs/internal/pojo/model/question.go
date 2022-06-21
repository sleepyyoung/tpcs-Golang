package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type QuestionType struct {
	Id   *int    `gorm:"column:ID;primaryKey" json:"id"`
	Name *string `gorm:"column:NAME" json:"name"`
}

func (q *QuestionType) TableName() string {
	return "question_type_info"
}

func (q *QuestionType) Create(db *gorm.DB) error {
	return db.Create(&q).Error
}

type QuestionDifficulty struct {
	Id   *int    `gorm:"column:ID;primaryKey" json:"id"`
	Name *string `gorm:"column:NAME" json:"name"`
}

func (q *QuestionDifficulty) TableName() string {
	return "question_difficulty_info"
}

type Question struct {
	Id                   *int                `gorm:"column:ID;primaryKey" json:"id"`
	Score                *float64            `gorm:"column:SCORE" json:"score"`
	UserId               *int                `gorm:"column:USER_ID"`
	User                 *User               `gorm:"column:USER_ID;foreignKey:UserId;references:USER_ID" json:"user"`
	QuestionTypeId       *int                `gorm:"column:question_type_id"`
	QuestionType         *QuestionType       `gorm:"column:TYPE_ID;foreignKey:QuestionTypeId;references:TYPE_ID" json:"type"`
	QuestionDifficultyId *int                `gorm:"column:question_difficulty_id"`
	QuestionDifficulty   *QuestionDifficulty `gorm:"column:DIFFICULTY_ID;foreignKey:QuestionDifficultyId;references:DIFFICULTY_ID" json:"difficulty"`
	CourseId             *int                `gorm:"column:course_id"`
	Course               *Course             `gorm:"column:COURSE_ID;foreignKey:CourseId;references:COURSE_ID" json:"course"`
	QuestionMd           *string             `gorm:"column:QUESTION_MD" json:"questionMd"`
	QuestionTxt          *string             `gorm:"column:QUESTION_TXT" json:"questionTxt"`
	QuestionHtml         *string             `gorm:"column:QUESTION_HTML" json:"questionHtml"`
	AnswerMd             *string             `gorm:"column:ANSWER_MD" json:"answerMd"`
	AnswerTxt            *string             `gorm:"column:ANSWER_TXT" json:"answerTxt"`
	AnswerHtml           *string             `gorm:"column:ANSWER_HTML" json:"answerHtml"`
}

func (cp *Question) TableName() string {
	return "question_info"
}

func (cp *Question) String() string {
	result := "{"
	if cp.Id != nil {
		result += fmt.Sprintf("Id: %v,", *cp.Id)
	} else {
		result += fmt.Sprintf("Id: %v,", "<nil>")
	}
	if cp.Score != nil {
		result += fmt.Sprintf("Score: %v,", *cp.Score)
	} else {
		result += fmt.Sprintf("Score: %v,", "<nil>")
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
	if cp.QuestionTypeId != nil {
		result += fmt.Sprintf("QuestionTypeId: %v,", *cp.QuestionTypeId)
	} else {
		result += fmt.Sprintf("QuestionTypeId: %v,", "<nil>")
	}
	if cp.QuestionType != nil {
		result += fmt.Sprintf("QuestionType: %v,", *cp.QuestionType)
	} else {
		result += fmt.Sprintf("QuestionType: %v,", "<nil>")
	}
	if cp.QuestionDifficultyId != nil {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", *cp.QuestionDifficultyId)
	} else {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", "<nil>")
	}
	if cp.QuestionDifficulty != nil {
		result += fmt.Sprintf("QuestionDifficulty: %v,", *cp.QuestionDifficulty)
	} else {
		result += fmt.Sprintf("QuestionDifficulty: %v,", "<nil>")
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
	if cp.QuestionMd != nil {
		result += fmt.Sprintf("QuestionMd: %v,", *cp.QuestionMd)
	} else {
		result += fmt.Sprintf("QuestionMd: %v,", "<nil>")
	}
	if cp.QuestionTxt != nil {
		result += fmt.Sprintf("QuestionTxt: %v,", *cp.QuestionTxt)
	} else {
		result += fmt.Sprintf("QuestionTxt: %v,", "<nil>")
	}
	if cp.QuestionHtml != nil {
		result += fmt.Sprintf("QuestionHtml: %v,", *cp.QuestionHtml)
	} else {
		result += fmt.Sprintf("QuestionHtml: %v,", "<nil>")
	}
	if cp.AnswerMd != nil {
		result += fmt.Sprintf("AnswerMd: %v,", *cp.AnswerMd)
	} else {
		result += fmt.Sprintf("AnswerMd: %v,", "<nil>")
	}
	if cp.AnswerTxt != nil {
		result += fmt.Sprintf("AnswerTxt: %v,", *cp.AnswerTxt)
	} else {
		result += fmt.Sprintf("AnswerTxt: %v,", "<nil>")
	}
	if cp.AnswerHtml != nil {
		result += fmt.Sprintf("AnswerHtml: %v,", *cp.AnswerHtml)
	} else {
		result += fmt.Sprintf("AnswerHtml: %v,", "<nil>")
	}
	result += "}"

	return result
}

// IndexQuestion “题目列表”页面用
type IndexQuestion struct {
	Id                   *int     `gorm:"column:ID;primaryKey" json:"id"`
	Score                *float64 `gorm:"column:SCORE" json:"score"`
	QuestionTypeId       *int     `gorm:"column:question_type_id"`
	QuestionDifficultyId *int     `gorm:"column:question_difficulty_id"`
	CourseId             *int     `gorm:"column:course_id"`
	QuestionMd           *string  `gorm:"column:QUESTION_MD" json:"questionMd"`
	QuestionTxt          *string  `gorm:"column:QUESTION_TXT" json:"questionTxt"`
	QuestionHtml         *string  `gorm:"column:QUESTION_HTML" json:"questionHtml"`
	AnswerMd             *string  `gorm:"column:ANSWER_MD" json:"answerMd"`
	AnswerTxt            *string  `gorm:"column:ANSWER_TXT" json:"answerTxt"`
	AnswerHtml           *string  `gorm:"column:ANSWER_HTML" json:"answerHtml"`
}

func (iq IndexQuestion) TableName() string {
	return "question_info"
}

// AddQuestion “添加”页面用
type AddQuestion struct {
	Score                *float64 `json:"score" gorm:"column:SCORE"`
	UserId               *int     `gorm:"column:USER_ID" json:"user"`
	QuestionTypeId       *int     `json:"type" gorm:"column:TYPE_ID"`
	QuestionDifficultyId *int     `json:"difficulty" gorm:"column:DIFFICULTY_ID"`
	CourseId             *int     `json:"course" gorm:"column:COURSE_ID"`
	QuestionMd           *string  `json:"questionMd" gorm:"column:QUESTION_MD"`
	QuestionTxt          *string  `json:"questionTxt" gorm:"column:QUESTION_TXT"`
	QuestionHtml         *string  `json:"questionHtml" gorm:"column:QUESTION_HTML"`
	AnswerMd             *string  `json:"answerMd" gorm:"column:ANSWER_MD"`
	AnswerTxt            *string  `json:"answerTxt" gorm:"column:ANSWER_TXT"`
	AnswerHtml           *string  `json:"answerHtml" gorm:"column:ANSWER_HTML"`
}

func (aq *AddQuestion) TableName() string {
	return "question_info"
}

func (aq *AddQuestion) String() string {
	result := "{"
	if aq.Score != nil {
		result += fmt.Sprintf("Score: %v,", *aq.Score)
	} else {
		result += fmt.Sprintf("Score: %v,", "<nil>")
	}
	if aq.UserId != nil {
		result += fmt.Sprintf("UserId: %v,", *aq.UserId)
	} else {
		result += fmt.Sprintf("UserId: %v,", "<nil>")
	}
	if aq.QuestionTypeId != nil {
		result += fmt.Sprintf("QuestionTypeId: %v,", *aq.QuestionTypeId)
	} else {
		result += fmt.Sprintf("QuestionTypeId: %v,", "<nil>")
	}
	if aq.QuestionDifficultyId != nil {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", *aq.QuestionDifficultyId)
	} else {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", "<nil>")
	}
	if aq.CourseId != nil {
		result += fmt.Sprintf("CourseId: %v,", *aq.CourseId)
	} else {
		result += fmt.Sprintf("CourseId: %v,", "<nil>")
	}
	if aq.QuestionMd != nil {
		result += fmt.Sprintf("QuestionMd: %v,", *aq.QuestionMd)
	} else {
		result += fmt.Sprintf("QuestionMd: %v,", "<nil>")
	}
	if aq.QuestionTxt != nil {
		result += fmt.Sprintf("QuestionTxt: %v,", *aq.QuestionTxt)
	} else {
		result += fmt.Sprintf("QuestionTxt: %v,", "<nil>")
	}
	if aq.QuestionHtml != nil {
		result += fmt.Sprintf("QuestionHtml: %v,", *aq.QuestionHtml)
	} else {
		result += fmt.Sprintf("QuestionHtml: %v,", "<nil>")
	}
	if aq.AnswerMd != nil {
		result += fmt.Sprintf("AnswerMd: %v,", *aq.AnswerMd)
	} else {
		result += fmt.Sprintf("AnswerMd: %v,", "<nil>")
	}
	if aq.AnswerTxt != nil {
		result += fmt.Sprintf("AnswerTxt: %v,", *aq.AnswerTxt)
	} else {
		result += fmt.Sprintf("AnswerTxt: %v,", "<nil>")
	}
	if aq.AnswerHtml != nil {
		result += fmt.Sprintf("AnswerHtml: %v", *aq.AnswerHtml)
	} else {
		result += fmt.Sprintf("AnswerHtml: %v", "<nil>")
	}
	result += "}"

	return result
}

// ModifyQuestion “修改”页面用
type ModifyQuestion struct {
	Score                *float64 `json:"score" gorm:"column:SCORE"`
	QuestionTypeId       *int     `json:"type" gorm:"column:TYPE_ID"`
	QuestionDifficultyId *int     `json:"difficulty" gorm:"column:DIFFICULTY_ID"`
	CourseId             *int     `json:"course" gorm:"column:COURSE_ID"`
	QuestionMd           *string  `json:"questionMd" gorm:"column:QUESTION_MD"`
	QuestionTxt          *string  `json:"questionTxt" gorm:"column:QUESTION_TXT"`
	QuestionHtml         *string  `json:"questionHtml" gorm:"column:QUESTION_HTML"`
	AnswerMd             *string  `json:"answerMd" gorm:"column:ANSWER_MD"`
	AnswerTxt            *string  `json:"answerTxt" gorm:"column:ANSWER_TXT"`
	AnswerHtml           *string  `json:"answerHtml" gorm:"column:ANSWER_HTML"`
}

func (mq *ModifyQuestion) TableName() string {
	return "question_info"
}

func (mq *ModifyQuestion) String() string {
	result := "{"
	if mq.Score != nil {
		result += fmt.Sprintf("Score: %v,", *mq.Score)
	} else {
		result += fmt.Sprintf("Score: %v,", "<nil>")
	}
	if mq.QuestionTypeId != nil {
		result += fmt.Sprintf("QuestionTypeId: %v,", *mq.QuestionTypeId)
	} else {
		result += fmt.Sprintf("QuestionTypeId: %v,", "<nil>")
	}
	if mq.QuestionDifficultyId != nil {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", *mq.QuestionDifficultyId)
	} else {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", "<nil>")
	}
	if mq.CourseId != nil {
		result += fmt.Sprintf("CourseId: %v,", *mq.CourseId)
	} else {
		result += fmt.Sprintf("CourseId: %v,", "<nil>")
	}
	if mq.QuestionMd != nil {
		result += fmt.Sprintf("QuestionMd: %v,", *mq.QuestionMd)
	} else {
		result += fmt.Sprintf("QuestionMd: %v,", "<nil>")
	}
	if mq.QuestionTxt != nil {
		result += fmt.Sprintf("QuestionTxt: %v,", *mq.QuestionTxt)
	} else {
		result += fmt.Sprintf("QuestionTxt: %v,", "<nil>")
	}
	if mq.QuestionHtml != nil {
		result += fmt.Sprintf("QuestionHtml: %v,", *mq.QuestionHtml)
	} else {
		result += fmt.Sprintf("QuestionHtml: %v,", "<nil>")
	}
	if mq.AnswerMd != nil {
		result += fmt.Sprintf("AnswerMd: %v,", *mq.AnswerMd)
	} else {
		result += fmt.Sprintf("AnswerMd: %v,", "<nil>")
	}
	if mq.AnswerTxt != nil {
		result += fmt.Sprintf("AnswerTxt: %v,", *mq.AnswerTxt)
	} else {
		result += fmt.Sprintf("AnswerTxt: %v,", "<nil>")
	}
	if mq.AnswerHtml != nil {
		result += fmt.Sprintf("AnswerHtml: %v", *mq.AnswerHtml)
	} else {
		result += fmt.Sprintf("AnswerHtml: %v", "<nil>")
	}
	result += "}"

	return result
}

// ExistsQuestionInfo “题库现存”api用
type ExistsQuestionInfo struct {
	Score                  *float64 `gorm:"column:SCORE"`
	QuestionTypeName       *string  `gorm:"column:type"`
	QuestionDifficultyName *string  `gorm:"column:difficulty"`
	Num                    *int     `gorm:"column:num"`
}

func (eqi *ExistsQuestionInfo) TableName() string {
	return "question_info"
}
