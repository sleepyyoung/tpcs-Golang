package course

import (
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
)

// CourseList 获取课程列表
func (svc *Service) CourseList(param *service.ListRequest) ([]model.Course, int, error) {
	return svc.dao.CourseList(global.DBEngine, param.Page, param.Limit)
}

// CreateCourse 添加课程
func (svc *Service) CreateCourse(name string) pojo.Result {
	exists, _ := svc.dao.CourseIsExistsByCourseName(global.DBEngine, name)
	if exists {
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_CourseExisted}
	}
	err := svc.dao.CreateCourse(name)
	if err != nil {
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// GetCourseById 通过id获取课程
func (svc *Service) GetCourseById(id int) (*model.Course, error) {
	return svc.dao.GetCourseById(global.DBEngine, id)
}
