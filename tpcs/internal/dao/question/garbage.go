package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/logger"
)

// RecoverQuestion 从回收站中恢复题目
func (d *Dao) RecoverQuestion(id int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID = ?", id).
		Update("REMOVED", 0).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// BatchRecoverQuestion 从回收站中批量恢复题目
func (d *Dao) BatchRecoverQuestion(ids []int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID in (?)", ids).
		Update("REMOVED", 0).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// DeleteQuestion 从回收站中彻底删除题目
func (d *Dao) DeleteQuestion(id int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID = ?", id).
		Delete(&model.Question{Id: &id}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// BatchDeleteQuestion 从回收站中批量彻底删除题目
func (d *Dao) BatchDeleteQuestion(ids []int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID in (?)", ids).
		Delete(&model.Question{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
