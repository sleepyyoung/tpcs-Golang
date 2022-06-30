package question

// QuestionScoreList 获取题目分值列表
func (svc *Service) QuestionScoreList(isRemoved bool) ([]float64, error) {
	return svc.dao.QuestionScoreList(isRemoved)
}
