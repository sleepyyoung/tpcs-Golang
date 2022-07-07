package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/logger"
)

// GetQuestionTypeByName 通过名称获取题型
func (d *Dao) GetQuestionTypeByName(name string) (*model.QuestionType, error) {
	db := global.DBEngine
	var questionType model.QuestionType
	if err := db.Table("question_type_info").
		Where("NAME = ?", name).
		Find(&questionType).
		Error; err != nil {
		return nil, err
	}
	return &questionType, nil
}

// AllQuestionTypes 获取题型列表（无分页）
func (d *Dao) AllQuestionTypes() ([]model.QuestionType, error) {
	db := global.DBEngine
	var questionTypeList []model.QuestionType
	if err := db.Table("question_type_info").
		Order("ID").
		Find(&questionTypeList).
		Error; err != nil {
		return nil, err
	}
	return questionTypeList, nil
}

// QuestionTypeList 题型列表（有分页）
func (d *Dao) QuestionTypeList(pageNum, pageSize int) ([]model.QuestionType, int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_type_info").Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionTypeList []model.QuestionType
	db = db.Table("question_type_info")
	if pageNum > 0 && pageSize > 0 {
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
	}
	db = db.Find(&questionTypeList)

	if err := db.Error; err != nil {
		return nil, 0, err
	}

	return questionTypeList, count, nil
}

// GetQuestionTypeById 通过题型id获取题型
func (d *Dao) GetQuestionTypeById(id int) (*model.QuestionType, error) {
	db := global.DBEngine
	var questionType model.QuestionType
	if err := db.Table("question_type_info").
		Where("ID = ?", id).
		Find(&questionType).
		Error; err != nil {
		return nil, err
	}

	return &questionType, nil
}

// GetQuestionDifficultyById 通过难度id获取题型
func (d *Dao) GetQuestionDifficultyById(id int) (*model.QuestionDifficulty, error) {
	db := global.DBEngine
	var questionDifficulty model.QuestionDifficulty
	if err := db.Table("question_difficulty_info").
		Where("ID = ?", id).
		Find(&questionDifficulty).
		Error; err != nil {
		return nil, err
	}

	return &questionDifficulty, nil
}

// AddQuestionType 添加题型
func (d *Dao) AddQuestionType(questionTypeName string) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	err := tx.Create(&model.QuestionType{Name: &questionTypeName}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
