package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	courseService "tpcs/internal/service/course"
	questionService "tpcs/internal/service/question"
	userService "tpcs/internal/service/user"
)

// CombineHandler 手动组卷页面
func CombineHandler(c *gin.Context) {
	courseSvc := courseService.New(c.Request.Context())

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		global.Logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	c.HTML(http.StatusOK, "question-combine.tmpl", gin.H{
		"courseList": courseList,
	})
}

// AutoCombineHandler 自动组卷页面
func AutoCombineHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		global.Logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	c.HTML(http.StatusOK, "question-combine-auto1.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
	})
}

// CombinePlanHandler 组卷方案页面
func CombinePlanHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		global.Logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	c.HTML(http.StatusOK, "question-combine-plan-list.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
	})
}

// AddCombinePlanHandler 新增组卷方案页面
func AddCombinePlanHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())
	userSvc := userService.New(c.Request.Context())

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		global.Logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user, _ := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
	userId := *user.Id

	c.HTML(http.StatusOK, "question-combine-plan-add.tmpl", gin.H{
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
		"userId":                 userId,
	})
}

// EditCombinePlanHandler 编辑组卷方案页面
func EditCombinePlanHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())
	userSvc := userService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		global.Logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		return
	}

	plan, err := questionSvc.GetCombinePlanById(id)
	if err != nil {
		global.Logger.Errorf("questionSvc.GetCombinePlanById err: %v", err)
		return
	}

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionTypes err: %v", err)
		return
	}

	questionDifficultyList, err := questionSvc.AllQuestionDifficulties()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		return
	}

	courseList, _, err := courseSvc.CourseList(&service.ListRequest{Page: 0, Limit: 0})
	if err != nil {
		global.Logger.Errorf("svc.CourseList err: %v", err)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user, _ := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
	userId := *user.Id

	c.HTML(http.StatusOK, "question-combine-plan-edit.tmpl", gin.H{
		"plan":                   plan,
		"questionTypeList":       questionTypeList,
		"questionDifficultyList": questionDifficultyList,
		"courseList":             courseList,
		"userId":                 userId,
	})
}

// CombinePlanDetailHandler 组卷方案详情页面
func CombinePlanDetailHandler(c *gin.Context) {
	questionSvc := questionService.New(c.Request.Context())
	userSvc := userService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		global.Logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		return
	}

	plan, err := questionSvc.GetCombinePlanById(id)
	if err != nil {
		global.Logger.Errorf("questionSvc.GetCombinePlanById err: %v", err)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user, _ := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
	userId := *user.Id

	c.HTML(http.StatusOK, "question-combine-plan-detail.tmpl", gin.H{
		"plan":   plan,
		"userId": userId,
	})
}
