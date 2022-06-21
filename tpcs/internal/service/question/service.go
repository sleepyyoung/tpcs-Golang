package question

import (
	"context"
	"fmt"
	"tpcs/global"
	"tpcs/internal/dao/question"
	"tpcs/internal/service"
)

type Service struct {
	ctx context.Context
	dao *question.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = question.New(global.DBEngine)
	return svc
}

// AddAndModifyQuestionRequest ”添加题目“和”编辑题目“请求体
type AddAndModifyQuestionRequest struct {
	Score                *float64 `json:"score" binding:"required"`
	QuestionTypeId       *int     `json:"type" binding:"required"`
	QuestionDifficultyId *int     `json:"difficulty" binding:"required"`
	CourseId             *int     `json:"course" binding:"required"`
	QuestionMd           *string  `json:"questionMd" binding:"required"`
	QuestionTxt          *string  `json:"questionTxt" binding:"required"`
	QuestionHtml         *string  `json:"questionHtml" binding:"required"`
	AnswerMd             *string  `json:"answerMd"`
	AnswerTxt            *string  `json:"answerTxt"`
	AnswerHtml           *string  `json:"answerHtml"`
}

func (amqr AddAndModifyQuestionRequest) String() string {
	result := "{"

	if amqr.Score != nil {
		result += fmt.Sprintf("Score: %v,", *amqr.Score)
	} else {
		result += fmt.Sprintf("Score: %v,", "<nil>")
	}
	if amqr.QuestionTypeId != nil {
		result += fmt.Sprintf("QuestionTypeId: %v,", *amqr.QuestionTypeId)
	} else {
		result += fmt.Sprintf("QuestionTypeId: %v,", "<nil>")
	}
	if amqr.QuestionDifficultyId != nil {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", *amqr.QuestionDifficultyId)
	} else {
		result += fmt.Sprintf("QuestionDifficultyId: %v,", "<nil>")
	}
	if amqr.CourseId != nil {
		result += fmt.Sprintf("CourseId: %v,", *amqr.CourseId)
	} else {
		result += fmt.Sprintf("CourseId: %v,", "<nil>")
	}
	if amqr.QuestionMd != nil {
		result += fmt.Sprintf("QuestionMd: %v,", *amqr.QuestionMd)
	} else {
		result += fmt.Sprintf("QuestionMd: %v,", "<nil>")
	}
	if amqr.QuestionTxt != nil {
		result += fmt.Sprintf("QuestionTxt: %v,", *amqr.QuestionTxt)
	} else {
		result += fmt.Sprintf("QuestionTxt: %v,", "<nil>")
	}
	if amqr.QuestionHtml != nil {
		result += fmt.Sprintf("QuestionHtml: %v,", *amqr.QuestionHtml)
	} else {
		result += fmt.Sprintf("QuestionHtml: %v,", "<nil>")
	}
	if amqr.AnswerMd != nil {
		result += fmt.Sprintf("AnswerMd: %v,", *amqr.AnswerMd)
	} else {
		result += fmt.Sprintf("AnswerMd: %v,", "<nil>")
	}
	if amqr.AnswerTxt != nil {
		result += fmt.Sprintf("AnswerTxt: %v,", *amqr.AnswerTxt)
	} else {
		result += fmt.Sprintf("AnswerTxt: %v,", "<nil>")
	}
	if amqr.AnswerHtml != nil {
		result += fmt.Sprintf("AnswerHtml: %v", *amqr.AnswerHtml)
	} else {
		result += fmt.Sprintf("AnswerHtml: %v", "<nil>")
	}
	result += "}"

	return result
}

// QueryQuestionRequest ”结合所有条件查询“请求体
type QueryQuestionRequest struct {
	ListRequest          service.ListRequest
	Score                *string `json:"score" form:"score"`
	QuestionTypeId       *int    `json:"type" form:"type"`
	QuestionDifficultyId *int    `json:"difficulty" form:"difficulty"`
	CourseId             *int    `json:"course" form:"course"`
	QuestionContent      *string `json:"questionContent" form:"questionContent"`
	AnswerContent        *string `json:"answerContent" form:"answerContent"`
}

// AddCombinePlanRequest 添加组卷方案请求体
type AddCombinePlanRequest struct {
	ListRequest service.ListRequest
	UserId      *int    `json:"user" form:"user"`
	CourseId    *int    `json:"course" form:"course"`
	PaperTitle  *string `json:"paperTitle" form:"title"`
	Plan        *string `json:"plan" form:"plan"`
	Score       *int    `json:"score" form:"score"`
	Note        *string `json:"note" form:"note"`
}

func (acpr AddCombinePlanRequest) String() string {
	result := "{"
	if acpr.UserId != nil {
		result += fmt.Sprintf("UserId: %v,", *acpr.UserId)
	} else {
		result += fmt.Sprintf("UserId: %v,", "<nil>")
	}
	if acpr.CourseId != nil {
		result += fmt.Sprintf("CourseId: %v,", *acpr.CourseId)
	} else {
		result += fmt.Sprintf("CourseId: %v,", "<nil>")
	}
	if acpr.PaperTitle != nil {
		result += fmt.Sprintf("PaperTitle: %v,", *acpr.PaperTitle)
	} else {
		result += fmt.Sprintf("PaperTitle: %v,", "<nil>")
	}
	if acpr.Plan != nil {
		result += fmt.Sprintf("Plan: %v,", *acpr.Plan)
	} else {
		result += fmt.Sprintf("Plan: %v,", "<nil>")
	}
	if acpr.Note != nil {
		result += fmt.Sprintf("Note: %v,", *acpr.Note)
	} else {
		result += fmt.Sprintf("Note: %v,", "<nil>")
	}
	result += "}"

	return result
}
