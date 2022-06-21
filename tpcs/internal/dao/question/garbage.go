package question

import (
	"github.com/jinzhu/gorm"
	"tpcs/internal/pojo/model"
)

// RecoverQuestion 从回收站中恢复题目
func (d *Dao) RecoverQuestion(db *gorm.DB, id int) error {
	if err := db.Table("question_info").
		Where("ID = ?", id).
		Update("REMOVED", 0).
		Error; err != nil {
		return err
	}
	return nil
}

// BatchRecoverQuestion 从回收站中批量恢复题目
func (d *Dao) BatchRecoverQuestion(db *gorm.DB, ids []int) error {
	if err := db.Table("question_info").
		Where("ID in (?)", ids).
		Update("REMOVED", 0).
		Error; err != nil {
		return err
	}
	return nil
}

// DeleteQuestion 从回收站中彻底删除题目
func (d *Dao) DeleteQuestion(db *gorm.DB, id int) error {
	if err := db.Table("question_info").
		Where("ID = ?", id).
		Delete(&model.Question{Id: &id}).
		Error; err != nil {
		return err
	}
	return nil
}

// BatchDeleteQuestion 从回收站中批量彻底删除题目
func (d *Dao) BatchDeleteQuestion(db *gorm.DB, ids []int) error {
	if err := db.Table("question_info").
		Where("ID in (?)", ids).
		Delete(&model.Question{}).
		Error; err != nil {
		return err
	}
	return nil
}
