package question

import (
	"github.com/jinzhu/gorm"
	"tpcs/internal/pojo/model"
)

// AllQuestionDifficulties 获取题目难度列表
func (d *Dao) AllQuestionDifficulties(db *gorm.DB) ([]model.QuestionDifficulty, error) {
	var questionDifficultyList []model.QuestionDifficulty
	if err := db.Table("question_difficulty_info").
		Order("ID").
		Find(&questionDifficultyList).
		Error; err != nil {
		return nil, err
	}
	return questionDifficultyList, nil
}
