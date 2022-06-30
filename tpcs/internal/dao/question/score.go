package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
)

// QuestionScoreList 获取题目分值列表
func (d *Dao) QuestionScoreList(isRemoved bool) ([]float64, error) {
	db := global.DBEngine
	var questionList []model.Question
	if err := db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("distinct SCORE").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("REMOVED = ?", isRemoved).
		Order("SCORE").
		Find(&questionList).
		Error; err != nil {
		return nil, err
	}

	scoreList := make([]float64, 0, len(questionList))
	for _, question := range questionList {
		scoreList = append(scoreList, *question.Score)
	}

	return scoreList, nil
}
