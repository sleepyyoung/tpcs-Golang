package course

import (
	"github.com/jinzhu/gorm"
	"tpcs/internal/pojo/model"
)

// CourseCount 获取课程数量
func (d *Dao) CourseCount(db *gorm.DB) (int, error) {
	var count int
	if err := db.Table("course_info").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CourseList 获取课程列表
func (d *Dao) CourseList(db *gorm.DB, pageNum, pageSize int) ([]model.Course, int, error) {
	count, err := d.CourseCount(db)
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
func (d *Dao) CourseIsExistsByCourseName(db *gorm.DB, name string) (bool, error) {
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
	return (&model.Course{Name: &name}).Create(d.engine)
}

// GetCourseById 通过id获取课程
func (d *Dao) GetCourseById(db *gorm.DB, id int) (*model.Course, error) {
	var course model.Course
	if err := db.Table("course_info").Where("ID = ?", id).Find(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
