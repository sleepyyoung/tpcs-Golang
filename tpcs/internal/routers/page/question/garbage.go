package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	courseService "tpcs/internal/service/course"
	questionService "tpcs/internal/service/question"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/logger"
)

// QuestionGarbageHandler 题目回收站
func QuestionGarbageHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	userSvc := userService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	isAdminUser, err := userSvc.IsAdminByUsernameAndPassword(*user.Username, *user.Password)

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		logger.Errorf("svc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	c.HTML(http.StatusOK, "question-garbage.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
		"isAdminUser":            isAdminUser,
	})
}

// QuestionGarbageDetailHandler 回收站题目详情
func QuestionGarbageDetailHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	question, err := questionSvc.GetQuestionById(id, true)
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, "question-detail.tmpl", gin.H{"question": question})
}

// QuestionGarbageDistributionHandler 回收站题目分布
func QuestionGarbageDistributionHandler(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{})
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, "question-garbage-distribution.tmpl", gin.H{"courseList": courseList})
}
