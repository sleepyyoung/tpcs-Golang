package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
)

// AllQuestionDifficulties 获取题目难度列表
func (d *Dao) AllQuestionDifficulties() ([]model.QuestionDifficulty, error) {
	db := global.DBEngine
	var questionDifficultyList []model.QuestionDifficulty
	if err := db.Table("question_difficulty_info").
		Order("ID").
		Find(&questionDifficultyList).
		Error; err != nil {
		return nil, err
	}
	return questionDifficultyList, nil
}
