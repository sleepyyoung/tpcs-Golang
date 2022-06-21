package question

import (
	"github.com/jinzhu/gorm"
	"tpcs/internal/pojo/model"
)

// QueryIdListByTypeIdAndDifficultyIdAndScore 组卷用，取id
func (d *Dao) QueryIdListByTypeIdAndDifficultyIdAndScore(db *gorm.DB, courseId, typeId, difficultyId int, score float64) ([]int, error) {
	var indexQuestionList []model.IndexQuestion
	if err := db.Raw("select ID "+
		"        from question_info "+
		"        where TYPE_ID = ? "+
		"          and DIFFICULTY_ID = ? "+
		"          and SCORE = ? "+
		"          and COURSE_ID = ? ",
		typeId, difficultyId, score, courseId).Scan(&indexQuestionList).Error; err != nil {
		return nil, err
	}
	var list []int
	for _, question := range indexQuestionList {
		list = append(list, *question.Id)
	}
	return list, nil
}

// CombinePlanCount 获取组卷方案数量
func (d *Dao) CombinePlanCount(db *gorm.DB) (int, error) {
	var count int
	if err := db.Table("question_combine_plan_info").
		Count(&count).
		Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CombinePlanListWithUserId 组卷方案列表
func (d *Dao) CombinePlanListWithUserId(db *gorm.DB, userId, pageNum, pageSize int) ([]model.CombinePlan, int, error) {
	count, err := d.CombinePlanCount(db)
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
func (d *Dao) CombinePlanList(db *gorm.DB, pageNum, pageSize int) ([]model.CombinePlan, int, error) {
	count, err := d.CombinePlanCount(db)
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

// AddCombinePlan 添加组卷方案
func (d *Dao) AddCombinePlan(db *gorm.DB, plan model.CombinePlan4Add) error {
	if err := db.Table("question_combine_plan_info").Create(&plan).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCombinePlan 删除组卷方案
func (d *Dao) DeleteCombinePlan(db *gorm.DB, id int) error {
	if err := db.Table("question_combine_plan_info").
		Where("ID = ?", id).
		Delete(&model.CombinePlan{Id: &id}).
		Error; err != nil {
		return err
	}
	return nil
}

// EditCombinePlan 修改组卷方案
func (d *Dao) EditCombinePlan(db *gorm.DB, plan model.CombinePlan4Edit) error {
	if err := db.Model(&plan).
		Where("ID = ?", *plan.Id).
		Updates(plan).
		Error; err != nil {
		return err
	}
	return nil
}

// BatchDeleteCombinePlan 批量删除组卷方案
func (d *Dao) BatchDeleteCombinePlan(db *gorm.DB, ids []int) error {
	if err := db.Table("question_combine_plan_info").
		Where("ID in (?) ", ids).
		Delete(&model.CombinePlan{}).
		Error; err != nil {
		return err
	}
	return nil
}

// GetCombinePlanById 通过id获取组卷方案
func (d *Dao) GetCombinePlanById(db *gorm.DB, id int) (*model.CombinePlan, error) {
	var plan model.CombinePlan
	if err := db.
		Table("question_combine_plan_info").
		Preload("User").
		Preload("Course").
		Select("question_combine_plan_info.ID,"+
			"USER_ID,"+
			"COURSE_ID,"+
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
