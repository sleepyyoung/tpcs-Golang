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

// QuestionListHandler 题目列表
func QuestionListHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	userSvc := userService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	isAdminUser, err := userSvc.IsAdminByUsernameAndPassword(*user.Username, *user.Password)

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		logger.Errorf("questionSvc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		logger.Errorf("questionSvc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		logger.Errorf("questionSvc.CourseList err: %v", err)
		return
	}

	c.HTML(http.StatusOK, "question-list.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
		"isAdminUser":            isAdminUser,
	})
}

// AddQuestionHandler 添加题目
func AddQuestionHandler(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())
	questionSvc := questionService.New(c.Request.Context())

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

	c.HTML(http.StatusOK, "question-add.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
	})
}

// EditQuestionHandler 编辑题目
func EditQuestionHandler(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())
	questionSvc := questionService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	question, err := questionSvc.GetQuestionById(id, false)
	if err != nil {
		return
	}
	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{})
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "question-edit.tmpl", gin.H{
		"question":               &question,
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
	})
}

// QuestionDetailHandler 题目详情
func QuestionDetailHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	question, err := questionSvc.GetQuestionById(id, false)
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, "question-detail.tmpl", gin.H{"question": question})
}

// QuestionDistributionHandler 题目分布
func QuestionDistributionHandler(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{})
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, "question-distribution.tmpl", gin.H{"courseList": courseList})
}
