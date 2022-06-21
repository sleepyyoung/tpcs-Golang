package question

import (
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
)

// GetQuestionTypeByName 通过名称获取题型
func (svc *Service) GetQuestionTypeByName(name string) (*model.QuestionType, error) {
	return svc.dao.GetQuestionTypeByName(name, global.DBEngine)
}

// AllQuestionTypes 获取题型列表（无分页）
func (svc *Service) AllQuestionTypes() ([]model.QuestionType, error) {
	return svc.dao.AllQuestionTypes(global.DBEngine)
}

// QuestionTypeList 题型列表（有分页）
func (svc *Service) QuestionTypeList(param *service.ListRequest) ([]model.QuestionType, int, error) {
	return svc.dao.QuestionTypeList(global.DBEngine, param.Page, param.Limit)
}

// GetQuestionTypeById 通过题型id获取题型
func (svc *Service) GetQuestionTypeById(id int) (*model.QuestionType, error) {
	return svc.dao.GetQuestionTypeById(global.DBEngine, id)
}

// GetQuestionDifficultyById 通过难度id获取题型
func (svc *Service) GetQuestionDifficultyById(id int) (*model.QuestionDifficulty, error) {
	return svc.dao.GetQuestionDifficultyById(global.DBEngine, id)
}

// AddQuestionType 添加题型
func (svc *Service) AddQuestionType(questionTypeName string) pojo.Result {
	err := svc.dao.AddQuestionType(questionTypeName)
	if err != nil {
		global.Logger.Errorf("添加题型失败！原因：%v", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}
