package course

import (
	"github.com/gin-gonic/gin"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	courseService "tpcs/internal/service/course"
	"tpcs/pkg/app"
	"tpcs/pkg/logger"
)

type Course struct{}

func NewCourse() Course {
	return Course{}
}

func (course Course) List(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())
	response := app.NewResponse(c)

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var courseList []model.Course
	courseList, count, err = courseSvc.CourseList(&request)
	if err != nil {
		logger.Errorf("svc.SelectPaperListByUserId err: %v", err)
		return
	}

	var courseMapList []map[string]interface{}
	for _, course := range courseList {
		courseMapList = append(courseMapList, map[string]interface{}{
			"id":   course.Id,
			"name": course.Name,
		})
	}

	response.ToResponse(pojo.Result{
		Code:  0,
		Msg:   "",
		Count: count,
		Data:  courseMapList,
	})
}

func (course Course) Create(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())
	response := app.NewResponse(c)
	response.ToResponse(courseSvc.CreateCourse(c.Param("name")))
}
