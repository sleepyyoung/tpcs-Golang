package file

import (
	"github.com/jinzhu/gorm"
	"tpcs/global"
	"tpcs/internal/dao/user"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/file"
)

// GetFileById 通过id获取文件
func (d *Dao) GetFileById(db *gorm.DB, id int) (*model.File, error) {
	var f model.File
	if err := db.Table("file_info").
		Select("ID, USER_ID, TRUTH_NAME, FILE_NAME, GMT_CREATE").
		Where("ID = ?", id).
		Find(&f).
		Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// FileCount 获取文件数量
func (d *Dao) FileCount(db *gorm.DB) (int, error) {
	var count int
	if err := db.Table("file_info").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FileList 获取文件列表
func (d *Dao) FileList(db *gorm.DB, userId, pageNum, pageSize int) ([]model.File, int, error) {
	isAdminUser, err := user.New(global.DBEngine).IsAdminUserByUserId(db, userId)
	if err != nil {
		global.Logger.Errorf("d.IsAdminUserByUserId err: %v", err)
		return nil, 0, err
	}

	count, err := d.FileCount(db)
	if err != nil {
		return nil, 0, err
	}

	var paperList []model.File
	db = db.Select("ID, USER_ID, TRUTH_NAME, FILE_NAME, GMT_CREATE")
	if !isAdminUser {
		db = db.Where("USER_ID = ?", userId)
	}
	db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&paperList)

	if err := db.Error; err != nil {
		return nil, 0, err
	}
	return paperList, count, nil
}

// AddFile 添加文件
func (d *Dao) AddFile(db *gorm.DB, file model.File4Add) error {
	return file.Create(db)
}

// DeleteFile 删除文件
func (d *Dao) DeleteFile(db *gorm.DB, id int) error {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		global.Logger.Errorf("事务开启异常: %v\n", err)
	}

	if err := tx.Table("file_info").
		Where("ID = ?", id).
		Delete(&model.File{Id: &id}).
		Error; err != nil {
		global.Logger.Errorf("删除文件异常，事务回滚。异常原因: %v\n", err)
		tx.Rollback()
		return err
	}

	fileById, err := d.GetFileById(global.DBEngine, id)
	if err != nil {
		global.Logger.Errorf("通过Id获取文件异常，事务回滚。异常原因: %v\n", err)
		tx.Rollback()
		return err
	}
	err = file.DeleteFile(global.AppSetting.UploadDir + "/" + *fileById.FileName)
	if err != nil {
		global.Logger.Errorf("文件物理删除失败，事务回滚。失败原因: %v\n", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// BatchDeleteFile 批量删除文件
func (d *Dao) BatchDeleteFile(db *gorm.DB, ids []int) error {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		global.Logger.Errorf("事务开启异常: %v\n", err)
	}

	for _, id := range ids {
		fileById, err := d.GetFileById(global.DBEngine, id)
		if err != nil {
			global.Logger.Errorf("通过Id获取文件异常，事务回滚。异常原因: %v\n", err)
			tx.Rollback()
			return err
		}

		err = file.DeleteFile(global.AppSetting.UploadDir + "/" + *fileById.FileName)
		if err != nil {
			global.Logger.Errorf("文件物理删除失败，事务回滚。失败原因: %v\n", err)
			tx.Rollback()
			return err
		}
	}

	if err := db.Table("file_info").
		Where("ID in (?)", ids).
		Delete(&model.File{}).
		Error; err != nil {
		global.Logger.Errorf("删除文件异常，事务回滚。异常原因: %v\n", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}