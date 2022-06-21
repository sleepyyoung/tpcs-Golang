package question

import (
	"tpcs/global"
)

// QuestionScoreList 获取题目分值列表
func (svc *Service) QuestionScoreList(isRemoved bool) ([]float64, error) {
	return svc.dao.QuestionScoreList(isRemoved, global.DBEngine)
}
