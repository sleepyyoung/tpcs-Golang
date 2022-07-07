package question

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	questionService "tpcs/internal/service/question"
	"tpcs/pkg/app"
	"tpcs/pkg/logger"
)

//// *********** Garbage *********** //

type Garbage struct{}

func NewGarbage() Garbage {
	return Garbage{}
}

// List 回收站题目列表
func (g Garbage) List(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var listRequest service.ListRequest
	err := c.ShouldBindQuery(&listRequest)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QuestionList(&listRequest, true)
	if err != nil {
		logger.Errorf("svc.QuestionList err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var questionMapList []map[string]interface{}
	for _, question := range questionList {
		questionMapList = append(questionMapList, map[string]interface{}{
			"id":           question.Id,
			"score":        question.Score,
			"type":         question.QuestionType.Name,
			"difficulty":   question.QuestionDifficulty.Name,
			"course":       question.Course.Name,
			"content":      question.QuestionTxt,
			"LAY_DISABLED": !(user.Id != nil && *user.Id == *question.UserId),
		})
	}

	response.ToResponse(pojo.Result{
		Code:  0,
		Msg:   "",
		Count: count,
		Data:  questionMapList,
	})
}

// Recover 恢复题目
func (g Garbage) Recover(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var result pojo.Result
	result = svc.RecoverQuestion(id)
	response.ToResultResponse(&result)
}

// Delete 从回收站中彻底删除题目
func (g Garbage) Delete(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var result pojo.Result
	result = svc.DeleteQuestion(id)
	response.ToResultResponse(&result)
}

// BatchRecover 批量恢复题目
func (g Garbage) BatchRecover(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("ioutil.ReadAll err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	ids := make([]int, 15)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		logger.Errorf("json.Unmarshal err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = svc.BatchRecoverQuestion(ids)
	response.ToResultResponse(&result)
}

// BatchDelete 批量删除题目
func (g Garbage) BatchDelete(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("ioutil.ReadAll err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	ids := make([]int, 15)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		logger.Errorf("json.Unmarshal err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = svc.BatchDeleteQuestion(ids)
	response.ToResultResponse(&result)
}

// Query 结合所有条件查询
func (g Garbage) Query(c *gin.Context) {
	Query(c, true)
}

// PreciseQueryByScore 仅凭分值精确查询
func (g Garbage) PreciseQueryByScore(c *gin.Context) {
	PreciseQueryByScore(c, true)
}

// IntervalQueryByScore 仅凭分值区间查询
func (g Garbage) IntervalQueryByScore(c *gin.Context) {
	IntervalQueryByScore(c, true)
}

// QueryByType 通过题型查询
func (g Garbage) QueryByType(c *gin.Context) {
	QueryByType(c, true)
}

// QueryByDifficulty 通过难度查询
func (g Garbage) QueryByDifficulty(c *gin.Context) {
	QueryByDifficulty(c, true)
}

// QueryByCourse 通过所属课程查询
func (g Garbage) QueryByCourse(c *gin.Context) {
	QueryByCourse(c, true)
}

// QueryByQuestionContent 通过题目内容查询
func (g Garbage) QueryByQuestionContent(c *gin.Context) {
	QueryByQuestionContent(c, true)
}

// QueryByAnswerContent 通过答案内容查询
func (g Garbage) QueryByAnswerContent(c *gin.Context) {
	QueryByAnswerContent(c, true)
}
