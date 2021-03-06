package course

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/logger"
)

// CourseCount 获取课程数量
func (d *Dao) CourseCount() (int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("course_info").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CourseList 获取课程列表
func (d *Dao) CourseList(pageNum, pageSize int) ([]model.Course, int, error) {
	db := global.DBEngine
	count, err := d.CourseCount()
	if err != nil {
		return nil, 0, err
	}

	var courseList []model.Course
	db = db.Table("course_info")
	if pageNum > 0 && pageSize > 0 {
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
	}
	db = db.Find(&courseList)

	if err := db.Error; err != nil {
		return nil, 0, err
	}

	return courseList, count, nil
}

// CourseIsExistsByCourseName 通过课程名判断课程是否存在
func (d *Dao) CourseIsExistsByCourseName(name string) (bool, error) {
	db := global.DBEngine
	var count int
	if err := db.
		Table("course_info").
		Where("NAME = ?", name).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count == 1, nil
}

// CreateCourse 添加课程
func (d *Dao) CreateCourse(name string) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	err := tx.Create(&model.Course{Name: &name}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetCourseById 通过id获取课程
func (d *Dao) GetCourseById(id int) (*model.Course, error) {
	db := global.DBEngine
	var course model.Course
	if err := db.Table("course_info").Where("ID = ?", id).Find(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
