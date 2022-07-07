package question

import (
	"tpcs/internal/pojo"
	"tpcs/pkg/logger"
)

// RecoverQuestion 从回收站中恢复题目
func (svc *Service) RecoverQuestion(id int) pojo.Result {
	err := svc.dao.RecoverQuestion(id)
	if err != nil {
		logger.Errorf("恢复题目失败！原因：%v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// BatchRecoverQuestion 从回收站中批量恢复题目
func (svc *Service) BatchRecoverQuestion(ids []int) pojo.Result {
	err := svc.dao.BatchRecoverQuestion(ids)
	if err != nil {
		logger.Errorf("批量恢复题目失败！原因：%v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// DeleteQuestion 从回收站中彻底删除题目
func (svc *Service) DeleteQuestion(id int) pojo.Result {
	err := svc.dao.DeleteQuestion(id)
	if err != nil {
		logger.Errorf("彻底删除题目失败！原因：%v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// BatchDeleteQuestion 从回收站中批量彻底删除题目
func (svc *Service) BatchDeleteQuestion(ids []int) pojo.Result {
	err := svc.dao.BatchDeleteQuestion(ids)
	if err != nil {
		logger.Errorf("批量彻底删除题目失败！原因：%v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}
