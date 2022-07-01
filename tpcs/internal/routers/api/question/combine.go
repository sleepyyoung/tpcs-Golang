package question

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	courseService "tpcs/internal/service/course"
	questionService "tpcs/internal/service/question"
	"tpcs/pkg/app"
	"tpcs/pkg/util/combine"
)

type Combine struct{}

func NewCombine() Combine {
	return Combine{}
}

// Init 组卷初始化用的数据
func (cb Combine) Init(c *gin.Context) {
	svc := questionService.New(c)
	response := app.NewResponse(c)

	questionTypeList, err := svc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionTypes err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	questionDifficultyList, err := svc.AllQuestionDifficulties()
	if err != nil {
		global.Logger.Errorf("svc.AllQuestionDifficulties err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, questionType := range questionTypeList {
		for _, questionDifficulty := range questionDifficultyList {
			list = append(list, map[string]interface{}{
				"type":         *questionType.Name,
				"typeId":       *questionType.Id,
				"difficulty":   *questionDifficulty.Name,
				"difficultyId": *questionDifficulty.Id,
				"num":          "",
				"score":        "",
				"dScore":       "",
				"tScore":       "",
				"LAY_CHECKED":  true,
			})
		}
	}

	response.ToResponse(pojo.Result{
		Code: 0,
		Msg:  "",
		Data: list,
	})
}

// ExistsInfoByCourse 题库现存题型及数量
func (cb Combine) ExistsInfoByCourse(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	courseId, err := strconv.Atoi(c.Param("courseId"))
	if err != nil {
		global.Logger.Errorf("strconv.Atoi err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	if courseId == 0 {
		response.ToResponse(pojo.Result{
			Code: 1,
			Msg:  "请选择课程！",
			Data: nil,
		})
		return
	}

	existsQuestionInfoList, err := svc.ExistsQuestionInfoList(courseId)
	if err != nil {
		global.Logger.Errorf("svc.ExistsQuestionInfoList err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, questionInfo := range existsQuestionInfoList {
		list = append(list, map[string]interface{}{
			"type":       *questionInfo.QuestionTypeName,
			"difficulty": *questionInfo.QuestionDifficultyName,
			"score":      *questionInfo.Score,
			"num":        *questionInfo.Num,
		})
	}

	if len(list) == 0 {
		response.ToResponse(pojo.Result{
			Code: 1,
			Msg:  "题库中没有该课程相关题目，请重新选择课程或去题库中添加题目！",
			Data: nil,
		})
		return
	}

	response.ToResponse(pojo.Result{
		Code: 0,
		Msg:  "",
		Data: list,
	})
	return
}

// Update 更新组卷页面的表格数据
func (cb Combine) Update(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		global.Logger.Errorf("ioutil.ReadAll err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	tableDataList := make([]map[string]interface{}, 0)
	err = json.Unmarshal(body, &tableDataList)
	if err != nil {
		global.Logger.Errorf("json.Unmarshal err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	list1 := make([]map[string]interface{}, 0)
	tScoreMap := make(map[int]float64)
	for _, m := range tableDataList {
		numS := strings.Trim(m["num"].(string), " ")
		m["LAY_CHECKED"] = true
		typeId := int(m["typeId"].(float64))
		var num int
		if numS == "" {
			num = 0
			//m["score"]=""
		} else {
			num, _ = strconv.Atoi(numS)
		}
		scoreS := strings.Trim(m["score"].(string), " ")
		var score float64
		if scoreS == "" {
			score = 0
			//m["num"]=""
		} else {
			score, _ = strconv.ParseFloat(scoreS, 64)
		}
		m["dScore"] = float64(num) * score
		dScore := m["dScore"].(float64)
		if numS != "" && scoreS != "" {
			tScoreDefault := tScoreMap[typeId]
			tScoreMap[typeId] = dScore + tScoreDefault
		}
		list1 = append(list1, m)
	}

	list2 := make([]map[string]interface{}, 0)
	for _, m := range list1 {
		var dScoreS, tScoreS string
		dScore := m["dScore"].(float64)
		if dScore == 0.0 {
			dScoreS = ""
		} else {
			dScoreS = strconv.FormatFloat(dScore, 'f', -1, 64)
		}
		tScore := tScoreMap[int(m["typeId"].(float64))]
		if tScore == 0.0 {
			tScoreS = ""
		} else {
			tScoreS = strconv.FormatFloat(tScore, 'f', -1, 64)
		}
		questionType, err := svc.GetQuestionTypeById(int(m["typeId"].(float64)))
		if err != nil {
			global.Logger.Errorf("svc.GetQuestionTypeById err: %v", err)
			return
		}
		m = map[string]interface{}{
			"LAY_CHECKED":  true,
			"type":         *questionType.Name,
			"typeId":       m["typeId"],
			"score":        m["score"],
			"scoreId":      m["scoreId"],
			"difficulty":   m["difficulty"],
			"difficultyId": m["difficultyId"],
			"num":          m["num"],
			"dScore":       dScoreS,
			"tScore":       tScoreS,
		}
		list2 = append(list2, m)
	}

	sum := 0.0
	for _, v := range tScoreMap {
		sum += v
	}
	response.ToResponse(map[string]interface{}{
		"table": list2,
		"sum":   sum,
	})
}

func doCombine(c *gin.Context) (bool, string, *pojo.CombineResult) {
	var request questionService.AddCombinePlanRequest
	err := c.ShouldBindWith(&request, binding.FormPost)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		return false, pojo.ResultMsg_FormParseErr, nil
	}

	questionSvc := questionService.New(c.Request.Context())
	courseSvc := courseService.New(c.Request.Context())
	combine.Svc = questionSvc

	tableDataList := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(*request.Plan), &tableDataList)
	if err != nil {
		global.Logger.Errorf("json.Unmarshal err: %v", err)
		return false, pojo.ResultMsg_TryAgainLater, nil
	}

	questionTypeList, err := questionSvc.AllQuestionTypes()
	if err != nil {
		global.Logger.Errorf("questionSvc.AllQuestionTypes err: %v", err)
		return false, pojo.ResultMsg_TryAgainLater, nil
	}

	combineQuestionMap := map[int]*combine.Question{}
	for _, questionType := range questionTypeList {
		combineQuestionMap[*questionType.Id] = nil
	}

	combine.QuestionMap = combineQuestionMap
	for _, m := range tableDataList {
		if "" != fmt.Sprintf("%v", m["score"]) && "" != fmt.Sprintf("%v", m["num"]) {
			typeId, err := strconv.Atoi(fmt.Sprintf("%v", m["typeId"]))
			if err != nil {
				global.Logger.Errorf("strconv.Atoi err: %v", err)
				return false, pojo.ResultMsg_FormParseErr, nil
			}
			difficultyId, err := strconv.Atoi(fmt.Sprintf("%v", m["difficultyId"]))
			if err != nil {
				global.Logger.Errorf("strconv.Atoi err: %v", err)
				return false, pojo.ResultMsg_FormParseErr, nil
			}
			score, err := strconv.ParseFloat(fmt.Sprintf("%v", m["score"]), 64)
			if err != nil {
				global.Logger.Errorf("strconv.ParseFloat err: %v", err)
				return false, pojo.ResultMsg_FormParseErr, nil
			}
			tScore, err := strconv.ParseFloat(fmt.Sprintf("%v", m["tScore"]), 64)
			if err != nil {
				global.Logger.Errorf("strconv.ParseFloat err: %v", err)
				return false, pojo.ResultMsg_FormParseErr, nil
			}
			num, err := strconv.Atoi(fmt.Sprintf("%v", m["num"]))
			if err != nil {
				global.Logger.Errorf("strconv.Atoi err: %v", err)
				return false, pojo.ResultMsg_FormParseErr, nil
			}
			ids, err := questionSvc.QueryIdListByTypeIdAndDifficultyIdAndScore(*request.CourseId, typeId, difficultyId, score)
			if err != nil {
				global.Logger.Errorf("questionSvc.QueryIdListByTypeIdAndDifficultyIdAndScore err: %v", err)
				return false, pojo.ResultMsg_TryAgainLater, nil
			}

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })

			if len(ids) < num {
				course, err := courseSvc.GetCourseById(*request.CourseId)
				if err != nil {
					global.Logger.Errorf("questionSvc.GetCourseById err: %v", err)
					return false, pojo.ResultMsg_TryAgainLater, nil
				}
				questionType, err := questionSvc.GetQuestionTypeById(typeId)
				if err != nil {
					global.Logger.Errorf("questionSvc.GetQuestionTypeById err: %v", err)
					return false, pojo.ResultMsg_TryAgainLater, nil
				}
				questionDifficulty, err := questionSvc.GetQuestionDifficultyById(difficultyId)
				if err != nil {
					global.Logger.Errorf("questionSvc.GetQuestionDifficultyById err: %v", err)
					return false, pojo.ResultMsg_TryAgainLater, nil
				}
				return false, "课程为[<span style=\"color: red;\">" + *course.Name +
					"</span>]、题型为[<span style=\"color: red;\">" + *questionType.Name +
					"</span>]、难度为[<span style=\"color: red;\">" + *questionDifficulty.Name +
					//"</span>]、分值为[<span style=\"color: red;\">" + strconv.FormatFloat(score, 'f', -1, 64) + // 按老师要求对分值不做硬性限制
					"</span>]的题目不足[<span style=\"color: red;\">" + strconv.Itoa(num) + "</span>]道", nil
			}

			subList := ids[:num]
			if _, ok := combineQuestionMap[typeId]; ok {
				combineQuestion := combineQuestionMap[typeId]
				if combineQuestion != nil {
					questionIdList := combineQuestion.QuestionIdList
					if questionIdList == nil {
						questionIdList = subList
					} else {
						questionIdList = append(questionIdList, subList...)
					}
					combineQuestionMap[typeId] = &combine.Question{
						QuestionIdList: questionIdList,
						Score:          score,
						TScore:         tScore,
					}
				} else {
					combineQuestionMap[typeId] = &combine.Question{
						QuestionIdList: subList,
						Score:          score,
						TScore:         tScore,
					}
				}
			}
		}
	}

	var cr pojo.CombineResult
	cr = combine.QuestionCombine(*request.PaperTitle)
	return true, "", &cr
}

// Combine 组卷
func (cb Combine) Combine(c *gin.Context) {
	var request questionService.AddCombinePlanRequest
	err := c.ShouldBindWith(&request, binding.FormPost)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		return
	}

	response := app.NewResponse(c)
	b, s, cr := doCombine(c)
	if b {
		response.ToCombineResultResponse(cr)
	} else {
		response.ToFailCombineResultResponse(s)
	}
}

// CombinePlanList 组卷方案列表
func (cb Combine) CombinePlanList(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)

	var request service.ListRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindQuery err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var count int
	var combinePlanList []model.CombinePlan
	combinePlanList, count, err = svc.CombinePlanList(&request)
	if err != nil {
		global.Logger.Errorf("svc.QuestionList err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var combinePlanMapList []map[string]interface{}
	for _, plan := range combinePlanList {
		combinePlanMapList = append(combinePlanMapList, map[string]interface{}{
			"id":           *plan.Id,
			"name":         *plan.PlanName,
			"user":         *plan.User.Username,
			"course":       *plan.Course.Name,
			"title":        *plan.PaperTitle,
			"plan":         *plan.Plan,
			"score":        *plan.Score,
			"note":         *plan.Note,
			"LAY_DISABLED": !(user.Id != nil && *user.Id == *plan.UserId),
		})
	}

	response.ToResponse(pojo.Result{
		Code:  0,
		Msg:   "",
		Count: count,
		Data:  combinePlanMapList,
	})

}

// AddCombinePlan 添加组卷方案
func (cb Combine) AddCombinePlan(c *gin.Context) {
	questionSvc := questionService.New(c)
	response := app.NewResponse(c)

	var request questionService.AddCombinePlanRequest
	err := c.ShouldBindWith(&request, binding.FormPost)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	b, s, _ := doCombine(c)
	if b {
		plan, err := questionSvc.GetCombinePlanByPlanName(strings.Trim(*request.PlanName, " "))
		if err != nil {
			global.Logger.Errorf("questionSvc.GetCombinePlanByPlanName err: %v\n", err)
			response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
			return
		}
		if plan != nil {
			response.ToFailResultResponse(pojo.ResultMsg_QuestionCombinePlanNameExisted)
			return
		}

		var result pojo.Result
		result = questionSvc.AddCombinePlan(request)
		response.ToResultResponse(&result)
		return
	} else {
		response.ToFailResultResponse(s)
		return
	}
}

// EditCombinePlan 修改组卷方案
func (cb Combine) EditCombinePlan(c *gin.Context) {
	questionSvc := questionService.New(c)
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		global.Logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var request questionService.AddCombinePlanRequest
	err = c.ShouldBindWith(&request, binding.FormPost)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	b, s, _ := doCombine(c)
	if b {
		plan, err := questionSvc.GetCombinePlanByPlanName(strings.Trim(*request.PlanName, " "))
		if err != nil {
			global.Logger.Errorf("questionSvc.GetCombinePlanByPlanName err: %v\n", err)
			response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
			return
		}
		if plan != nil && *plan.Id != id {
			response.ToFailResultResponse(pojo.ResultMsg_QuestionCombinePlanNameExisted)
			return
		}

		var result pojo.Result
		result = questionSvc.EditCombinePlan(id, request)
		response.ToResultResponse(&result)
		return
	} else {
		response.ToFailResultResponse(s)
		return
	}
}

// DeleteCombinePlan 删除组卷方案
func (cb Combine) DeleteCombinePlan(c *gin.Context) {
	questionSvc := questionService.New(c)
	response := app.NewResponse(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		global.Logger.Errorf("strconv.Atoi(c.Param(\"id\")) err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	var result pojo.Result
	result = questionSvc.DeleteCombinePlan(id)
	response.ToResultResponse(&result)
	return
}

// BatchDeleteCombinePlan 批量删除组卷方案
func (cb Combine) BatchDeleteCombinePlan(c *gin.Context) {
	svc := questionService.New(c.Request.Context())
	response := app.NewResponse(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		global.Logger.Errorf("ioutil.ReadAll err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	ids := make([]int, 15)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		global.Logger.Errorf("json.Unmarshal err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	}

	var result pojo.Result
	result = svc.BatchDeleteCombinePlan(ids)
	response.ToResultResponse(&result)
}
