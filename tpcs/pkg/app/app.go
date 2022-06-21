package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tpcs/internal/pojo"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResultResponse(result *pojo.Result) {
	r.Ctx.JSON(http.StatusOK, result)
}

func (r *Response) ToSuccessResultResponse() {
	r.Ctx.JSON(http.StatusOK, &pojo.Result{Success: pojo.ResultSuccess_True})
}

func (r *Response) ToFailResultResponse(msg string) {
	r.Ctx.JSON(http.StatusOK, &pojo.Result{Success: pojo.ResultSuccess_False, Msg: msg})
}

func (r *Response) ToCombineResultResponse(result *pojo.CombineResult) {
	r.Ctx.JSON(http.StatusOK, result)
}

func (r *Response) ToFailCombineResultResponse(msg string) {
	r.Ctx.JSON(http.StatusOK, &pojo.CombineResult{Success: pojo.ResultSuccess_False, PaperHtml: msg, AnswerHtml: msg})
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}
