package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"tpcs/global"
	"tpcs/internal/middleware"
	auth_ "tpcs/internal/routers/api/auth"
	course_ "tpcs/internal/routers/api/course"
	file_ "tpcs/internal/routers/api/file"
	question_ "tpcs/internal/routers/api/question"
	user_ "tpcs/internal/routers/api/user"
	_4xxPageRouter "tpcs/internal/routers/page/4xx"
	coursePageRouter "tpcs/internal/routers/page/course"
	filePageRouter "tpcs/internal/routers/page/file"
	lrmfPageRouter "tpcs/internal/routers/page/lrmf"
	mpagePageRouter "tpcs/internal/routers/page/mpage"
	questionPageRouter "tpcs/internal/routers/page/question"
	userPageRouter "tpcs/internal/routers/page/user"
)

func NewRouter() *gin.Engine {
	store := cookie.NewStore([]byte("secret"))
	r := gin.New()
	r.Use(sessions.SessionsMany(global.AppSetting.SessionNames, store))
	r.SetFuncMap(template.FuncMap{
		"equal": func(a, b int) bool {
			return strconv.Itoa(a) == strconv.Itoa(b)
		},
	})

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	}

	r.LoadHTMLGlob("resources/templates/*/*")

	r.StaticFS("/static", http.Dir("resources/static"))
	r.StaticFile("/favicon.ico", "resources/static/favicon.ico")

	auth := auth_.NewAuth()
	user := user_.NewUser()
	file := file_.NewFile()
	course := course_.NewCourse()
	question := question_.NewQuestion()
	questionType := question_.NewQuestionType()
	garbage := question_.NewGarbage()
	combine := question_.NewCombine()
	teacher := user_.NewTeacher()

	r.GET("/404", _4xxPageRouter.NotFoundHandler)
	r.GET("/403", _4xxPageRouter.ForbiddenHandler)

	// -----------------------------------
	r.POST("/login", auth.Login)
	r.POST("/logout", auth.Logout)
	r.POST("/register", auth.Register)
	r.POST("/forgot-password", auth.ForgotPassword)
	r.GET("/login", lrmfPageRouter.LoginHandler)
	r.GET("/register", lrmfPageRouter.RegisterHandler)
	r.GET("/forgot-password/:username", lrmfPageRouter.ForgotPasswordHandler)

	registerSendVcodeApi := r.Group("/send-vcode")
	registerSendVcodeApi.POST("/register", auth.SendVcode4Register)
	registerSendVcodeApi.POST("/forgot", auth.SendVcode4Forgot)

	modifyPasswordApi := r.Group("/modify-password")
	modifyPasswordApi.Use(middleware.UserApiInterceptor())
	modifyPasswordApi.POST("/", auth.ModifyPassword)

	withAuthPage := r.Group("/")
	withAuthPage.Use(middleware.UserPageInterceptor())
	withAuthPage.GET("/modify-password/:username", lrmfPageRouter.ModifyPasswordHandler)
	withAuthPage.GET("/", mpagePageRouter.MainHandler)
	withAuthPage.GET("/index", mpagePageRouter.IndexHandler)
	withAuthPage.GET("/user/:username", userPageRouter.UserHandler)
	withAuthPage.GET("/question-list", questionPageRouter.QuestionListHandler)
	withAuthPage.GET("/question-add", questionPageRouter.AddQuestionHandler)
	withAuthPage.GET("/question-edit/:id", questionPageRouter.EditQuestionHandler)
	withAuthPage.GET("/question-detail/:id", questionPageRouter.QuestionDetailHandler)
	withAuthPage.GET("/question-distribution", questionPageRouter.QuestionDistributionHandler)
	withAuthPage.GET("/file-list", filePageRouter.FileListHandler)
	withAuthPage.GET("/file-upload", filePageRouter.FileUploadHandler)
	withAuthPage.GET("/question-garbage", questionPageRouter.QuestionGarbageHandler)
	withAuthPage.GET("/question-garbage-detail/:id", questionPageRouter.QuestionGarbageDetailHandler)
	withAuthPage.GET("/question-garbage-distribution", questionPageRouter.QuestionGarbageDistributionHandler)

	withAdminAuthPage := r.Group("/")
	withAdminAuthPage.GET("/teacher-audit", userPageRouter.TeacherAuditByJWTHandler)
	withAdminAuthPage.Use(middleware.AdminPageInterceptor())
	withAdminAuthPage.GET("/course-list", coursePageRouter.CourseListHandler)
	withAdminAuthPage.GET("/teacher-list", userPageRouter.TeacherListHandler)
	withAdminAuthPage.GET("/teacher-audit/:id", userPageRouter.TeacherAuditBySessionHandler)
	withAdminAuthPage.GET("/question-type-list", questionPageRouter.QuestionTypeListHandler)

	withTeacherAuthPage := r.Group("/")
	withTeacherAuthPage.Use(middleware.TeacherPageInterceptor())
	withTeacherAuthPage.GET("/question-combine", questionPageRouter.CombineHandler)
	withTeacherAuthPage.GET("/question-combine-auto/:id", questionPageRouter.AutoCombineHandler)
	withTeacherAuthPage.GET("/question-combine-plan-list", questionPageRouter.CombinePlanHandler)
	withTeacherAuthPage.GET("/question-combine-plan-add", questionPageRouter.AddCombinePlanHandler)
	withTeacherAuthPage.GET("/question-combine-plan-edit/:id", questionPageRouter.EditCombinePlanHandler)
	withTeacherAuthPage.GET("/question-combine-plan-detail/:id", questionPageRouter.CombinePlanDetailHandler)

	userApi := r.Group("/api/user")
	userApi.Use(middleware.UserApiInterceptor())
	userApi.PUT("/:userId/username", user.UpdateUsername)
	userApi.PUT("/:userId/email", user.UpdateEmail)
	userApi.POST("/:userId/send-vcode", user.SendVcode)

	fileApi := r.Group("/api/files")
	fileApi.Use(middleware.UserApiInterceptor())
	fileApi.GET("/", file.List)
	fileApi.POST("/", file.Upload)
	fileApi.POST("/md-img", file.Upload4MdImg)
	fileApi.GET("/status", file.UploadStatus)
	fileApi.GET("/download/:fileId", file.Download)
	fileApi.DELETE("/:id", file.Delete)
	fileApi.DELETE("/batch", file.BatchDelete)

	courseApi := r.Group("/api/courses")
	courseApi.Use(middleware.AdminApiInterceptor())
	courseApi.GET("/", course.List)
	courseApi.POST("/:name", course.Create)

	teacherApi := r.Group("/api/teachers")
	teacherApi.POST("/audit2", teacher.Audit2)
	teacherApi.Use(middleware.AdminApiInterceptor())
	teacherApi.GET("/", teacher.List)
	teacherApi.POST("/freeze", teacher.Freeze)
	teacherApi.POST("/thaw", teacher.Thaw)
	teacherApi.POST("/audit", teacher.Audit)

	questionTypeApi := r.Group("/api/question-types")
	questionTypeApi.Use(middleware.AdminApiInterceptor())
	questionTypeApi.GET("/", questionType.List)
	questionTypeApi.POST("/:name", questionType.Create)

	combineApi := r.Group("/api/combine")
	combineApi.Use(middleware.TeacherApiInterceptor())
	combineApi.GET("/init", combine.Init)
	combineApi.GET("/exists/:courseId", combine.ExistsInfoByCourse)
	combineApi.PUT("/", combine.Update)
	combineApi.POST("/", combine.Combine)
	combinePlanApi := combineApi.Group("/plans")
	combinePlanApi.GET("/", combine.CombinePlanList)
	combinePlanApi.POST("/", combine.AddCombinePlan)
	combinePlanApi.PUT("/:id", combine.EditCombinePlan)
	combinePlanApi.DELETE("/:id", combine.DeleteCombinePlan)
	combinePlanApi.DELETE("/batch", combine.BatchDeleteCombinePlan)

	questionApi := r.Group("/api/questions")
	questionApi.Use(middleware.UserApiInterceptor())
	questionApi.GET("/", question.List)
	questionApi.POST("/", question.Add)
	questionApi.PUT("/:id", question.Edit)
	questionApi.DELETE("/:id", question.Remove)
	questionApi.DELETE("/batch", question.BatchRemove)
	questionApi.GET("/distribution/:courseId", question.GetQuestionDistribution)
	queryQuestionApi := questionApi.Group("/query")
	queryQuestionApi.GET("/", question.Query)
	queryQuestionApi.GET("/only-me", question.OnlyMe)
	queryQuestionByScoreApi := queryQuestionApi.Group("/score")
	queryQuestionByScoreApi.GET("/interval/:min/:max", question.IntervalQueryByScore)
	queryQuestionByScoreApi.GET("/precise/:score", question.PreciseQueryByScore)
	queryQuestionByTypeApi := queryQuestionApi.Group("/type")
	queryQuestionByTypeApi.GET("/:type", question.QueryByType)
	queryQuestionByDifficultyApi := queryQuestionApi.Group("/difficulty")
	queryQuestionByDifficultyApi.GET("/:difficulty", question.QueryByDifficulty)
	queryQuestionByCourseApi := queryQuestionApi.Group("/course")
	queryQuestionByCourseApi.GET("/:course", question.QueryByCourse)
	queryQuestionByQuestionContentApi := queryQuestionApi.Group("/question-content")
	queryQuestionByQuestionContentApi.GET("/:content", question.QueryByQuestionContent)
	queryQuestionByAnswerContentApi := queryQuestionApi.Group("/answer-content")
	queryQuestionByAnswerContentApi.GET("/:content", question.QueryByAnswerContent)

	questionGarbageApi := questionApi.Group("/garbages")
	questionGarbageApi.Use(middleware.UserApiInterceptor())
	questionGarbageApi.GET("/", garbage.List)
	questionGarbageApi.POST("/:id", garbage.Recover)
	questionGarbageApi.DELETE("/:id", garbage.Delete)
	questionGarbageApi.PUT("/batch", garbage.BatchRecover)
	questionGarbageApi.DELETE("/batch", garbage.BatchDelete)
	questionGarbageApi.GET("/distribution/:courseId", question.GetQuestionGarbageDistribution)
	queryQuestionGarbageApi := questionGarbageApi.Group("/query")
	queryQuestionGarbageApi.GET("/", garbage.Query)
	queryQuestionGarbageByScoreApi := queryQuestionGarbageApi.Group("/score")
	queryQuestionGarbageByScoreApi.GET("/interval/:min/:max", garbage.IntervalQueryByScore)
	queryQuestionGarbageByScoreApi.GET("/precise/:score", garbage.PreciseQueryByScore)
	queryQuestionGarbageByTypeApi := queryQuestionGarbageApi.Group("/type")
	queryQuestionGarbageByTypeApi.GET("/:type", garbage.QueryByType)
	queryQuestionGarbageByDifficultyApi := queryQuestionGarbageApi.Group("/difficulty")
	queryQuestionGarbageByDifficultyApi.GET("/:difficulty", garbage.QueryByDifficulty)
	queryQuestionGarbageByCourseApi := queryQuestionGarbageApi.Group("/course")
	queryQuestionGarbageByCourseApi.GET("/:course", garbage.QueryByCourse)
	queryQuestionGarbageByQuestionContentApi := queryQuestionGarbageApi.Group("/question-content")
	queryQuestionGarbageByQuestionContentApi.GET("/:content", garbage.QueryByQuestionContent)
	queryQuestionGarbageByAnswerContentApi := queryQuestionGarbageApi.Group("/answer-content")
	queryQuestionGarbageByAnswerContentApi.GET("/:content", garbage.QueryByAnswerContent)

	return r
}
