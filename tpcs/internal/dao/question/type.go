package question

import (
	"github.com/jinzhu/gorm"
	"tpcs/global"
	"tpcs/internal/pojo/model"
)

// GetQuestionTypeByName 通过名称获取题型
func (d *Dao) GetQuestionTypeByName(name string, db *gorm.DB) (*model.QuestionType, error) {
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
func (d *Dao) AllQuestionTypes(db *gorm.DB) ([]model.QuestionType, error) {
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
func (d *Dao) QuestionTypeList(db *gorm.DB, pageNum, pageSize int) ([]model.QuestionType, int, error) {
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
func (d *Dao) GetQuestionTypeById(db *gorm.DB, id int) (*model.QuestionType, error) {
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
func (d *Dao) GetQuestionDifficultyById(db *gorm.DB, id int) (*model.QuestionDifficulty, error) {
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
	var questionType = model.QuestionType{Name: &questionTypeName}
	err := questionType.Create(global.DBEngine)
	if err != nil {
		return err
	}
	return nil
}
