package question

import (
	"strconv"
	"strings"
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
)

// GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty 通过课程id、题型、分值、难度获取题目数量
func (svc *Service) GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(courseId int, typeId int, score float64, difficultyId int, isRemoved bool) (int, error) {
	return svc.dao.GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(courseId, typeId, score, difficultyId, isRemoved)
}

// ExistsQuestionInfoList 题库现存题目信息
func (svc *Service) ExistsQuestionInfoList(courseId int) ([]model.ExistsQuestionInfo, error) {
	return svc.dao.ExistsQuestionInfoList(courseId)
}

// QuestionList 获取题目列表
func (svc *Service) QuestionList(param *service.ListRequest, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QuestionList(isRemoved, param.Page, param.Limit)
}

// AddQuestion 添加题目
func (svc *Service) AddQuestion(userid *int, param *AddAndModifyQuestionRequest) pojo.Result {
	err := svc.dao.AddQuestion(model.AddQuestion{
		Score:                param.Score,
		UserId:               userid,
		QuestionTypeId:       param.QuestionTypeId,
		QuestionDifficultyId: param.QuestionDifficultyId,
		CourseId:             param.CourseId,
		QuestionMd:           param.QuestionMd,
		QuestionTxt:          param.QuestionTxt,
		QuestionHtml:         param.QuestionHtml,
		AnswerMd:             param.AnswerMd,
		AnswerTxt:            param.AnswerTxt,
		AnswerHtml:           param.AnswerHtml,
	})
	if err != nil {
		global.Logger.Errorf("添加题目失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// GetQuestionById 通过id获取题目
func (svc *Service) GetQuestionById(id int, isRemoved bool) (*model.Question, error) {
	return svc.dao.GetQuestionById(isRemoved, id)
}

// GetQuestionByUserId 通过所属用户获取题目
func (svc *Service) GetQuestionByUserId(param *service.ListRequest, id int, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.GetQuestionByUserId(isRemoved, id, param.Page, param.Limit)
}

// ModifyQuestion 修改题目
func (svc *Service) ModifyQuestion(id int, param AddAndModifyQuestionRequest) pojo.Result {
	err := svc.dao.ModifyQuestion(id, model.ModifyQuestion{
		Score:                param.Score,
		QuestionTypeId:       param.QuestionTypeId,
		QuestionDifficultyId: param.QuestionDifficultyId,
		CourseId:             param.CourseId,
		QuestionMd:           param.QuestionMd,
		QuestionTxt:          param.QuestionTxt,
		QuestionHtml:         param.QuestionHtml,
		AnswerMd:             param.AnswerMd,
		AnswerTxt:            param.AnswerTxt,
		AnswerHtml:           param.AnswerHtml,
	})
	if err != nil {
		global.Logger.Errorf("修改题目失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// RemoveQuestion 通过id移除题目
func (svc *Service) RemoveQuestion(id int) pojo.Result {
	err := svc.dao.RemoveQuestion(id)
	if err != nil {
		global.Logger.Errorf("移除题目失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// BatchRemoveQuestion 通过id批量移除题目
func (svc *Service) BatchRemoveQuestion(ids []int) pojo.Result {
	err := svc.dao.BatchRemoveQuestion(ids)
	if err != nil {
		global.Logger.Errorf("批量移除题目失败！原因: %v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// QueryQuestion 综合所有条件查询
func (svc *Service) QueryQuestion(param *QueryQuestionRequest, isRemoved bool) ([]model.Question, int, error) {
	var score, min, max *float64
	var typeId, difficultyId, courseId *int
	var questionContent, answerContent *string
	score, min, max = new(float64), new(float64), new(float64)
	typeId, difficultyId, courseId = new(int), new(int), new(int)
	questionContent, answerContent = new(string), new(string)

	if param.Score != nil {
		scoreS := strings.Trim(*param.Score, " ")
		if scoreS != "" {
			s, err := strconv.ParseFloat(scoreS, 64)
			if err != nil {
				score = nil

				ss := strings.Split(scoreS, ",")
				minS, maxS := ss[0], ss[1]
				i, err := strconv.ParseFloat(minS, 64)
				if err != nil {
					return nil, 0, err
				}
				a, err := strconv.ParseFloat(maxS, 64)
				if err != nil {
					return nil, 0, err
				}

				*min, *max = i, a
			} else {
				*score = s
			}
		} else {
			score, min, max = nil, nil, nil
		}
	} else {
		score, min, max = nil, nil, nil
	}
	if param.QuestionTypeId != nil && *param.QuestionTypeId != 0 {
		*typeId = *param.QuestionTypeId
	} else {
		typeId = nil
	}
	if param.QuestionDifficultyId != nil && *param.QuestionDifficultyId != 0 {
		*difficultyId = *param.QuestionDifficultyId
	} else {
		difficultyId = nil
	}
	if param.CourseId != nil && *param.CourseId != 0 {
		*courseId = *param.CourseId
	} else {
		courseId = nil
	}
	if param.QuestionContent != nil && strings.Trim(*param.QuestionContent, " ") != "" {
		*questionContent = strings.Trim(*param.QuestionContent, " ")
	} else {
		questionContent = nil
	}
	if param.AnswerContent != nil && strings.Trim(*param.AnswerContent, " ") != "" {
		*answerContent = strings.Trim(*param.AnswerContent, " ")
	} else {
		answerContent = nil
	}

	return svc.dao.QueryQuestion(score, min, max,
		typeId, difficultyId, courseId,
		questionContent, answerContent, param.ListRequest.Page, param.ListRequest.Limit, isRemoved)
}

// PreciseQueryQuestionByScore 仅凭分值精确查询
func (svc *Service) PreciseQueryQuestionByScore(param *service.ListRequest, score float64, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.PreciseQueryQuestionByScore(score, isRemoved, param.Page, param.Limit)
}

// IntervalQueryQuestionByScore 仅凭分值区间查询
func (svc *Service) IntervalQueryQuestionByScore(param *service.ListRequest, min, max float64, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.IntervalQueryQuestionByScore(min, max, isRemoved, param.Page, param.Limit)
}

// QueryQuestionByType 通过题型查询
func (svc *Service) QueryQuestionByType(param *service.ListRequest, typeId int, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QueryQuestionByType(typeId, isRemoved, param.Page, param.Limit)
}

// QueryQuestionByDifficulty 通过难度查询
func (svc *Service) QueryQuestionByDifficulty(param *service.ListRequest, difficultyId int, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QueryQuestionByDifficulty(difficultyId, isRemoved, param.Page, param.Limit)
}

// QueryQuestionByCourse 通过所属课程查询
func (svc *Service) QueryQuestionByCourse(param *service.ListRequest, courseId int, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QueryQuestionByCourse(courseId, isRemoved, param.Page, param.Limit)
}

// QueryQuestionByQuestionContent 通过题目内容查询
func (svc *Service) QueryQuestionByQuestionContent(param *service.ListRequest, questionContent string, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QueryQuestionByQuestionContent(questionContent, isRemoved, param.Page, param.Limit)
}

// QueryQuestionByAnswerContent 通过答案内容查询
func (svc *Service) QueryQuestionByAnswerContent(param *service.ListRequest, answerContent string, isRemoved bool) ([]model.Question, int, error) {
	return svc.dao.QueryQuestionByAnswerContent(answerContent, isRemoved, param.Page, param.Limit)
}
