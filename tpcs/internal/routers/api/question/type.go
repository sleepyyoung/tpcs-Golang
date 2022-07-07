package question

import (
	"github.com/gin-gonic/gin"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	questionService "tpcs/internal/service/question"
	"tpcs/pkg/app"
	"tpcs/pkg/logger"
)

type Type struct{}

func NewQuestionType() Type {
	return Type{}
}

// List 题型列表
func (qt Type) List(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	var listRequest service.ListRequest
	err := c.ShouldBindQuery(&listRequest)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionTypeList []model.QuestionType
	questionTypeList, count, err = svc.QuestionTypeList(&listRequest)
	if err != nil {
		logger.Errorf("svc.QuestionTypeList err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var questionTypeMapList []map[string]interface{}
	for _, questionType := range questionTypeList {
		questionTypeMapList = append(questionTypeMapList, map[string]interface{}{
			"id":   questionType.Id,
			"name": questionType.Name,
		})
	}

	response.ToResponse(pojo.Result{
		Code:  0,
		Msg:   "",
		Count: count,
		Data:  questionTypeMapList,
	})
}

// Create 添加题型
func (qt Type) Create(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)
	response.ToResponse(svc.AddQuestionType(c.Param("name")))
}
