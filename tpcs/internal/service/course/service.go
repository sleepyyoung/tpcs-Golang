package course

import (
	"context"
	"tpcs/global"
	"tpcs/internal/dao/course"
)

type Service struct {
	ctx context.Context
	dao *course.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = course.New(global.DBEngine)
	return svc
}
