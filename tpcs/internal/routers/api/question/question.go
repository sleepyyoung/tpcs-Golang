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

//// ********** Public **********//

// Query 结合所有条件查询
func Query(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var request *questionService.QueryQuestionRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindJSON err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestion(request, isRemoved)
	if err != nil {
		logger.Errorf("svc.QueryQuestion err: %v", err)
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

// PreciseQueryByScore 仅凭分值精确查询
func PreciseQueryByScore(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	score, err := strconv.ParseFloat(c.Param("score"), 64)
	if err != nil {
		logger.Errorf("strconv.ParseFloat err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var request service.ListRequest
	err = c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.PreciseQueryQuestionByScore(&request, score, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// IntervalQueryByScore 仅凭分值区间查询
func IntervalQueryByScore(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	min, err := strconv.ParseFloat(c.Param("min"), 64)
	if err != nil {
		logger.Errorf("strconv.ParseFloat err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	max, err := strconv.ParseFloat(c.Param("max"), 64)
	if err != nil {
		logger.Errorf("strconv.ParseFloat err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var request service.ListRequest
	err = c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.IntervalQueryQuestionByScore(&request, min, max, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// QueryByType 通过题型查询
func QueryByType(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	typeId, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		logger.Errorf("strconv.Atoi err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var request service.ListRequest
	err = c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestionByType(&request, typeId, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// QueryByDifficulty 通过难度查询
func QueryByDifficulty(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	difficultyId, err := strconv.Atoi(c.Param("difficulty"))
	if err != nil {
		logger.Errorf("strconv.Atoi err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var request service.ListRequest
	err = c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestionByDifficulty(&request, difficultyId, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// QueryByCourse 通过所属课程查询
func QueryByCourse(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	courseId, err := strconv.Atoi(c.Param("course"))
	if err != nil {
		logger.Errorf("strconv.Atoi err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var request service.ListRequest
	err = c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestionByCourse(&request, courseId, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// QueryByQuestionContent 通过题目内容查询
func QueryByQuestionContent(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	questionContent := c.Param("content")

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestionByQuestionContent(&request, questionContent, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// QueryByAnswerContent 通过答案内容查询
func QueryByAnswerContent(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	answerContent := c.Param("content")

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QueryQuestionByAnswerContent(&request, answerContent, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// GetQuestionByUserId 通过所属用户获取题目
func GetQuestionByUserId(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.GetQuestionByUserId(&request, *user.Id, isRemoved)
	if err != nil {
		logger.Errorf("svc.PreciseQueryQuestionByScore err: %v", err)
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

// GetQuestionDistribution 获取题目分布
func GetQuestionDistribution(c *gin.Context, isRemoved bool) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	courseId, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"courseId\"))[id = %v] err: %v", courseId, err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	allQuestionTypes, err := svc.AllQuestionTypes()
	if err != nil {
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}
	y := make([]string, 0, len(allQuestionTypes))
	for _, questionType := range allQuestionTypes {
		y = append(y, *questionType.Name)
	}

	x, err := svc.QuestionScoreList(isRemoved)
	if err != nil {
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	dExample := make([][]interface{}, 0, 0)
	for iy := 0; iy < len(y); iy++ {
		for ix := 0; ix < len(x); ix++ {
			questionTypeByName, err := svc.GetQuestionTypeByName(y[iy])
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			c, err := svc.GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(
				courseId,
				*questionTypeByName.Id,
				x[ix],
				1,
				isRemoved,
			)
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			if c != 0 {
				list := make([]interface{}, 0, 3)
				list = append(list, x[ix], *questionTypeByName.Name, c)
				dExample = append(dExample, list)
			}
		}
	}

	dMiddle := make([][]interface{}, 0, 0)
	for iy := 0; iy < len(y); iy++ {
		for ix := 0; ix < len(x); ix++ {
			questionTypeByName, err := svc.GetQuestionTypeByName(y[iy])
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			c, err := svc.GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(
				courseId,
				*questionTypeByName.Id,
				x[ix],
				2,
				isRemoved,
			)
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			if c != 0 {
				list := make([]interface{}, 0, 3)
				list = append(list, x[ix], *questionTypeByName.Name, c)
				dMiddle = append(dMiddle, list)
			}
		}
	}

	dHard := make([][]interface{}, 0, 0)
	for iy := 0; iy < len(y); iy++ {
		for ix := 0; ix < len(x); ix++ {
			questionTypeByName, err := svc.GetQuestionTypeByName(y[iy])
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			c, err := svc.GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(
				courseId,
				*questionTypeByName.Id,
				x[ix],
				3,
				isRemoved,
			)
			if err != nil {
				response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
				return
			}
			if c != 0 {
				list := make([]interface{}, 0, 3)
				list = append(list, x[ix], *questionTypeByName.Name, c)
				dHard = append(dHard, list)
			}
		}
	}

	result := make(map[string]interface{})
	result["y"] = y
	result["x"] = x
	result["de"] = dExample
	result["dm"] = dMiddle
	result["dh"] = dHard

	marshalResult, err := json.Marshal(result)
	if err != nil {
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}
	response.ToResponse(pojo.Result{
		Code: 0,
		Msg:  "",
		Data: string(marshalResult),
	})
}

type Question struct{}

func NewQuestion() Question {
	return Question{}
}

// List 题目列表
func (q Question) List(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var questionList []model.Question
	questionList, count, err = svc.QuestionList(&request, false)
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

// Add 添加题目
func (q Question) Add(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	var request *questionService.AddAndModifyQuestionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindJSON err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	if request.QuestionMd == nil || *request.QuestionMd == "" ||
		request.QuestionHtml == nil || *request.QuestionHtml == "" ||
		request.QuestionTxt == nil || *request.QuestionTxt == "" {
		response.ToFailResultResponse(pojo.ResultMsg_QuestionContentNotNone)
		return
	}
	if request.Score == nil || *request.Score <= 0 {
		response.ToFailResultResponse(pojo.ResultMsg_QuestionScoreMoreThan0)
		return
	}
	if request.QuestionTypeId == nil || *request.QuestionTypeId == 0 {
		response.ToFailResultResponse(pojo.ResultMsg_QuestionTypeNotNone)
		return
	}
	if request.QuestionDifficultyId == nil || *request.QuestionDifficultyId == 0 {
		response.ToFailResultResponse(pojo.ResultMsg_QuestionDifficultyNotNone)
		return
	}
	if request.CourseId == nil || *request.CourseId == 0 {
		response.ToFailResultResponse(pojo.ResultMsg_QuestionCourseNotNone)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var result pojo.Result
	result = svc.AddQuestion(user.Id, request)
	response.ToResultResponse(&result)
}

// Edit 编辑题目
func (q Question) Edit(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\"))[id = %v] err: %v", id, err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	questionById, err := svc.GetQuestionById(id, false)
	if err != nil {
		logger.Errorf("svc.GetQuestionById(%v, false) err: %v", id, err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}
	if *user.Id != *questionById.UserId {
		response.ToFailResultResponse(pojo.ResultMsg_InsufficientPermissions)
		return
	}

	var request questionService.AddAndModifyQuestionRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindJSON err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = svc.ModifyQuestion(id, request)
	response.ToResultResponse(&result)
}

// Remove 移除题目
func (q Question) Remove(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var result pojo.Result
	result = svc.RemoveQuestion(id)
	response.ToResultResponse(&result)
}

// BatchRemove 批量移除题目
func (q Question) BatchRemove(c *gin.Context) {
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
	result = svc.BatchRemoveQuestion(ids)
	response.ToResultResponse(&result)
}

// Query 结合所有条件查询
func (q Question) Query(c *gin.Context) {
	Query(c, false)
}

// OnlyMe 只看我的
func (q Question) OnlyMe(c *gin.Context) {
	GetQuestionByUserId(c, false)
}

// PreciseQueryByScore 仅凭分值精确查询
func (q Question) PreciseQueryByScore(c *gin.Context) {
	PreciseQueryByScore(c, false)
}

// IntervalQueryByScore 仅凭分值区间查询
func (q Question) IntervalQueryByScore(c *gin.Context) {
	IntervalQueryByScore(c, false)
}

// QueryByType 通过题型查询
func (q Question) QueryByType(c *gin.Context) {
	QueryByType(c, false)
}

// QueryByDifficulty 通过难度查询
func (q Question) QueryByDifficulty(c *gin.Context) {
	QueryByDifficulty(c, false)
}

// QueryByCourse 通过所属课程查询
func (q Question) QueryByCourse(c *gin.Context) {
	QueryByCourse(c, false)
}

// QueryByQuestionContent 通过题目内容查询
func (q Question) QueryByQuestionContent(c *gin.Context) {
	QueryByQuestionContent(c, false)
}

// QueryByAnswerContent 通过答案内容查询
func (q Question) QueryByAnswerContent(c *gin.Context) {
	QueryByAnswerContent(c, false)
}

// GetQuestionDistribution 获取题目分布
func (q Question) GetQuestionDistribution(c *gin.Context) {
	GetQuestionDistribution(c, false)
}

// GetQuestionGarbageDistribution 获取回收站题目分布
func (q Question) GetQuestionGarbageDistribution(c *gin.Context) {
	GetQuestionDistribution(c, true)
}
