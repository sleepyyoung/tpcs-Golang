package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
)

// AllQuestionDifficulties 获取题目难度列表（无分页）
func (svc *Service) AllQuestionDifficulties() ([]model.QuestionDifficulty, error) {
	return svc.dao.AllQuestionDifficulties(global.DBEngine)
}
