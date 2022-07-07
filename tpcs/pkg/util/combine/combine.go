package combine

import (
	"sort"
	"strconv"
	"strings"
	"tpcs/internal/pojo"
	questionService "tpcs/internal/service/question"
	"tpcs/pkg/logger"
)

type Question struct {
	// 该题型题目列表
	QuestionIdList []int
	// 单个题目分值
	Score float64
	// 该题型总分值
	TScore float64
}

var (
	Svc         questionService.Service
	QuestionMap map[int]*Question
	bigNum      = []string{"一、", "二、", "三、", "四、", "五、", "六、", "七、", "八、", "九、", "十、",
		"十一、", "十二、", "十三、", "十四、", "十五、", "十六、", "十七、", "十八、", "十九、", "二十、"}
)

type sortable []int

func (p sortable) Len() int {
	return len(p)
}
func (p sortable) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p sortable) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func QuestionCombine(paperTitle string) pojo.CombineResult {
	bigNumFlag := 0
	paperHtml, paperHtmlTitle, answerHtml, answerHtmlTitle := "", "", "", ""
	paperHtmlTitle += "<h1 align=\"center\">" + paperTitle + "</h1><hr/>"
	answerHtmlTitle += "<h1 align=\"center\">" + paperTitle + "【答案】</h1><hr/>"

	list := make([]int, 0, len(QuestionMap))
	for k := range QuestionMap {
		list = append(list, k)
	}

	sort.Sort(sortable(list))
	for _, id := range list {
		combineQuestion := QuestionMap[id]
		if combineQuestion == nil {
			continue
		}
		questionType, err := Svc.GetQuestionTypeById(id)
		if err != nil {
			logger.Errorf("Svc.GetQuestionTypeById err: %v", err)
			return pojo.CombineResult{
				Success:    pojo.ResultSuccess_False,
				PaperHtml:  pojo.ResultMsg_TryAgainLater,
				AnswerHtml: pojo.ResultMsg_TryAgainLater,
			}
		}
		questionTypeName := *questionType.Name
		questionIdList := combineQuestion.QuestionIdList
		score, tScore := combineQuestion.Score, combineQuestion.TScore
		paperHtml += "<h5>" + bigNum[bigNumFlag] + questionTypeName +
			"题（每题" + strconv.FormatFloat(score, 'f', -1, 64) + "分，共" +
			strconv.FormatFloat(tScore, 'f', -1, 64) + "分）</h5>" + questionCombine(questionIdList)
		answerHtml += "<h5>" + bigNum[bigNumFlag] + questionTypeName + "题" + answerCombine(questionIdList)
		bigNumFlag++
	}

	if "" == paperHtml {
		return pojo.CombineResult{
			Success:    pojo.ResultSuccess_True,
			PaperHtml:  "",
			AnswerHtml: "",
		}
	} else {
		paperHtmlTitle += paperHtml
		answerHtmlTitle += answerHtml
	}

	return pojo.CombineResult{
		Success:    pojo.ResultSuccess_True,
		PaperHtml:  paperHtmlTitle,
		AnswerHtml: answerHtmlTitle,
	}
}

func questionCombine(list []int) string {
	html := ""
	for i := 1; i <= len(list); i++ {
		question, err := Svc.GetQuestionById(list[i-1], false)
		if err != nil {
			logger.Errorf("Svc.GetQuestionById err: %v", err)
			return err.Error()
		}
		questionHtml := *question.QuestionHtml
		questionHtml = strings.Replace(questionHtml, "<p>", "<p>"+strconv.Itoa(i)+"、", 1)
		html += questionHtml
	}
	return html
}

func answerCombine(list []int) string {
	html := ""
	for i := 1; i <= len(list); i++ {
		question, err := Svc.GetQuestionById(list[i-1], false)
		if err != nil {
			logger.Errorf("Svc.GetQuestionById err: %v", err)
			return err.Error()
		}
		answerHtml := *question.AnswerHtml
		html += "<h5>" + strconv.Itoa(i) + "、</h5>" + answerHtml
	}
	return html
}
