package file

import (
	"context"
	"tpcs/global"
	"tpcs/internal/dao/file"
)

type Service struct {
	ctx context.Context
	dao *file.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = file.New(global.DBEngine)
	return svc
}
