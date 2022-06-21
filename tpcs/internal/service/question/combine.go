package question

import (
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
)

// QueryIdListByTypeIdAndDifficultyIdAndScore 组卷用，取id
func (svc *Service) QueryIdListByTypeIdAndDifficultyIdAndScore(courseId, typeId, difficultyId int, score float64) ([]int, error) {
	return svc.dao.QueryIdListByTypeIdAndDifficultyIdAndScore(global.DBEngine, courseId, typeId, difficultyId, score)
}

// CombinePlanListWithUserId 组卷方案列表
func (svc *Service) CombinePlanListWithUserId(userId int, param *service.ListRequest) ([]model.CombinePlan, int, error) {
	return svc.dao.CombinePlanListWithUserId(global.DBEngine, userId, param.Page, param.Limit)
}

// CombinePlanList 组卷方案列表
func (svc *Service) CombinePlanList(param *service.ListRequest) ([]model.CombinePlan, int, error) {
	return svc.dao.CombinePlanList(global.DBEngine, param.Page, param.Limit)
}

// AddCombinePlan 添加组卷方案
func (svc *Service) AddCombinePlan(request AddCombinePlanRequest) pojo.Result {
	err := svc.dao.AddCombinePlan(global.DBEngine, model.CombinePlan4Add{
		UserId:     request.UserId,
		CourseId:   request.CourseId,
		PaperTitle: request.PaperTitle,
		Plan:       request.Plan,
		Score:      request.Score,
		Note:       request.Note,
	})
	if err != nil {
		global.Logger.Errorf("questionSvc.AddCombinePlan err: %v", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// EditCombinePlan 修改组卷方案
func (svc *Service) EditCombinePlan(id int, request AddCombinePlanRequest) pojo.Result {
	err := svc.dao.EditCombinePlan(global.DBEngine, model.CombinePlan4Edit{
		Id:         &id,
		UserId:     request.UserId,
		CourseId:   request.CourseId,
		PaperTitle: request.PaperTitle,
		Plan:       request.Plan,
		Score:      request.Score,
		Note:       request.Note,
	})
	if err != nil {
		global.Logger.Errorf("questionSvc.EditCombinePlan err: %v", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// DeleteCombinePlan 删除组卷方案
func (svc *Service) DeleteCombinePlan(id int) pojo.Result {
	err := svc.dao.DeleteCombinePlan(global.DBEngine, id)
	if err != nil {
		global.Logger.Errorf("questionSvc.DeleteCombinePlan err: %v", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// BatchDeleteCombinePlan 批量删除组卷方案
func (svc *Service) BatchDeleteCombinePlan(ids []int) pojo.Result {
	err := svc.dao.BatchDeleteCombinePlan(global.DBEngine, ids)
	if err != nil {
		global.Logger.Errorf("批量删除组卷方案失败！原因：%v\n", err)
		return pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_TryAgainLater}
	}
	return pojo.Result{Success: pojo.ResultSuccess_True}
}

// GetCombinePlanById 通过id获取组卷方案
func (svc *Service) GetCombinePlanById(id int) (*model.CombinePlan, error) {
	return svc.dao.GetCombinePlanById(global.DBEngine, id)
}
