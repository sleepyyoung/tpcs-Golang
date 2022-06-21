package service

import (
	"context"
	"tpcs/global"
	"tpcs/internal/dao/question"
)

type Service struct {
	ctx context.Context
	dao *question.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = question.New(global.DBEngine)
	return svc
}

// ListRequest 分页列表的请求体
type ListRequest struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=15"`
}
