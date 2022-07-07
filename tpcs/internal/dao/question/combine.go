package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/logger"
)

// QueryIdListByTypeIdAndDifficultyIdAndScore 组卷用，取id
// 按老师要求对分值不做硬性限制
func (d *Dao) QueryIdListByTypeIdAndDifficultyIdAndScore(courseId, typeId, difficultyId int, score float64) ([]int, error) {
	db := global.DBEngine
	var indexQuestionList []model.IndexQuestion
	if err := db.Raw("select ID "+
		"        from question_info "+
		"        where TYPE_ID = ? "+
		"          and DIFFICULTY_ID = ? "+
		//"          and SCORE = ? "+
		"          and COURSE_ID = ? ",
		typeId,
		difficultyId,
		//score,
		courseId,
	).Scan(&indexQuestionList).Error; err != nil {
		return nil, err
	}
	var list []int
	for _, question := range indexQuestionList {
		list = append(list, *question.Id)
	}
	return list, nil
}

// CombinePlanCount 获取组卷方案数量
func (d *Dao) CombinePlanCount() (int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_combine_plan_info").
		Count(&count).
		Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CombinePlanListWithUserId 组卷方案列表
func (d *Dao) CombinePlanListWithUserId(userId, pageNum, pageSize int) ([]model.CombinePlan, int, error) {
	db := global.DBEngine
	count, err := d.CombinePlanCount()
	if err != nil {
		return nil, 0, err
	}

	var combinePlanList []model.CombinePlan
	if err := db.
		Table("question_combine_plan_info").
		Preload("User").
		Preload("Course").
		Select("question_combine_plan_info.ID,"+
			"USER_ID,"+
			"COURSE_ID,"+
			"question_combine_plan_info.NAME,"+
			"PAPER_TITLE,"+
			"PLAN,"+
			"SCORE,"+
			"question_combine_plan_info.NOTE").
		Joins("left join user_info "+
			"on user_info.ID = question_combine_plan_info.USER_ID").
		Joins("left join course_info "+
			"on course_info.ID = question_combine_plan_info.COURSE_ID").
		Offset((pageNum-1)*pageSize).
		Where("USER_ID = ?", userId).
		Limit(pageSize).
		Find(&combinePlanList).
		Error; err != nil {
		return nil, 0, err
	}

	return combinePlanList, count, nil
}

// CombinePlanList 组卷方案列表
func (d *Dao) CombinePlanList(pageNum, pageSize int) ([]model.CombinePlan, int, error) {
	db := global.DBEngine
	count, err := d.CombinePlanCount()
	if err != nil {
		return nil, 0, err
	}

	var combinePlanList []model.CombinePlan
	if err := db.
		Table("question_combine_plan_info").
		Preload("User").
		Preload("Course").
		Select("question_combine_plan_info.ID," +
			"USER_ID," +
			"COURSE_ID," +
			"question_combine_plan_info.NAME," +
			"PAPER_TITLE," +
			"PLAN," +
			"SCORE," +
			"question_combine_plan_info.NOTE").
		Joins("left join user_info " +
			"on user_info.ID = question_combine_plan_info.USER_ID").
		Joins("left join course_info " +
			"on course_info.ID = question_combine_plan_info.COURSE_ID").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&combinePlanList).
		Error; err != nil {
		return nil, 0, err
	}

	return combinePlanList, count, nil
}

// GetCombinePlanByPlanName 通过方案名称获取组卷方案
func (d *Dao) GetCombinePlanByPlanName(planName string) (*model.CombinePlan, error) {
	db := global.DBEngine
	var plan model.CombinePlan
	var count = 0
	db = db.
		Table("question_combine_plan_info").
		Preload("User").
		Preload("Course").
		Select("question_combine_plan_info.ID,"+
			"USER_ID,"+
			"COURSE_ID,"+
			"question_combine_plan_info.NAME,"+
			"PAPER_TITLE,"+
			"PLAN,"+
			"SCORE,"+
			"question_combine_plan_info.NOTE").
		Joins("left join user_info "+
			"on user_info.ID = question_combine_plan_info.USER_ID").
		Joins("left join course_info "+
			"on course_info.ID = question_combine_plan_info.COURSE_ID").
		Where("question_combine_plan_info.NAME = ?", planName)

	db.Count(&count)
	if count > 0 {
		if err := db.Find(&plan).
			Error; err != nil {
			return nil, err
		}
		return &plan, nil
	} else {
		return nil, nil
	}
}

// AddCombinePlan 添加组卷方案
func (d *Dao) AddCombinePlan(plan model.CombinePlan4Add) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_combine_plan_info").Create(&plan).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// DeleteCombinePlan 删除组卷方案
func (d *Dao) DeleteCombinePlan(id int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_combine_plan_info").
		Where("ID = ?", id).
		Delete(&model.CombinePlan{Id: &id}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// EditCombinePlan 修改组卷方案
func (d *Dao) EditCombinePlan(plan model.CombinePlan4Edit) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Model(&plan).
		Where("ID = ?", *plan.Id).
		Updates(plan).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// BatchDeleteCombinePlan 批量删除组卷方案
func (d *Dao) BatchDeleteCombinePlan(ids []int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_combine_plan_info").
		Where("ID in (?) ", ids).
		Delete(&model.CombinePlan{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetCombinePlanById 通过id获取组卷方案
func (d *Dao) GetCombinePlanById(id int) (*model.CombinePlan, error) {
	db := global.DBEngine
	var plan model.CombinePlan
	if err := db.
		Table("question_combine_plan_info").
		Preload("User").
		Preload("Course").
		Select("question_combine_plan_info.ID,"+
			"USER_ID,"+
			"COURSE_ID,"+
			"question_combine_plan_info.NAME,"+
			"PAPER_TITLE,"+
			"PLAN,"+
			"SCORE,"+
			"question_combine_plan_info.NOTE").
		Joins("left join user_info "+
			"on user_info.ID = question_combine_plan_info.USER_ID").
		Joins("left join course_info "+
			"on course_info.ID = question_combine_plan_info.COURSE_ID").
		Where("question_combine_plan_info.ID = ?", id).
		Find(&plan).
		Error; err != nil {
		return nil, err
	}
	return &plan, nil
}
