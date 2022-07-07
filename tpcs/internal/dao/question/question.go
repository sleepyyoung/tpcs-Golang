package question

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/pkg/logger"
)

// GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty 通过课程id、题型、分值、难度获取题目数量
func (d *Dao) GetQuestionCountByCourseIdAndQuestionTypeAndScoreAndDifficulty(courseId int, typeId int, score float64, difficultyId int, isRemoved bool) (int, error) {
	db := global.DBEngine
	var count int
	db = db.Table("question_info").
		Where("TYPE_ID = ?", typeId).
		Where("SCORE = ?", score).
		Where("DIFFICULTY_ID = ?", difficultyId).
		Where("REMOVED = ?", isRemoved)
	if courseId > 0 {
		db = db.Where("COURSE_ID = ?", courseId)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsQuestionInfoList 题库现存题目信息
func (d *Dao) ExistsQuestionInfoList(courseId int) ([]model.ExistsQuestionInfo, error) {
	db := global.DBEngine
	var existsQuestionInfoList []model.ExistsQuestionInfo
	if err := db.Raw("select * "+
		"from ( "+
		"         select type.ID         type_id, "+
		"                type.NAME       type, "+
		"                difficulty.ID   difficulty_id, "+
		"                difficulty.NAME difficulty, "+
		"                SCORE, "+
		"                count(*)        num "+
		"         from question_info "+
		"                  left join question_type_info type on type.ID = question_info.TYPE_ID "+
		"                  left join question_difficulty_info difficulty on difficulty.ID = question_info.DIFFICULTY_ID "+
		"         where COURSE_ID = ? "+
		"           and REMOVED = 0 "+
		"         group by type.ID, difficulty.ID, SCORE "+
		"     ) temp "+
		"order by temp.type_id, temp.difficulty_id ", courseId).
		Scan(&existsQuestionInfoList).
		Error; err != nil {
		return nil, err
	}
	return existsQuestionInfoList, nil
}

// QuestionCount 获取题目数量
func (d *Dao) QuestionCount(isRemoved bool) (int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_info").
		Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return 0, err
	}
	return count, nil
}

// QuestionList 获取题目列表
func (d *Dao) QuestionList(isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	count, err := d.QuestionCount(isRemoved)
	if err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	if err := db.Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID,"+
			"SCORE,"+
			"USER_ID,"+
			"question_type.ID         question_type_id,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"course.ID                course_id,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// AddQuestion 添加题目
func (d *Dao) AddQuestion(question model.AddQuestion) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").Create(&question).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetQuestionById 通过id获取题目
func (d *Dao) GetQuestionById(isRemoved bool, id int) (*model.Question, error) {
	db := global.DBEngine
	var question model.Question
	if err := db.Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Preload("User").
		Select("question_info.ID,"+
			"SCORE,"+
			"USER_ID,"+
			"user_info.USERNAME,"+
			"question_type.ID         question_type_id,"+
			"question_type.NAME,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"question_difficulty.NAME,"+
			"course.ID                course_id,"+
			"course.NAME,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Joins("left join user_info "+
			"on user_info.ID = question_info.USER_ID").
		Where("question_info.ID = ?", id).
		Where("REMOVED = ?", isRemoved).
		Find(&question).
		Error; err != nil {
		return nil, err
	}
	return &question, nil
}

// GetQuestionByUserId 通过所属用户获取题目
func (d *Dao) GetQuestionByUserId(isRemoved bool, id int, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_info").
		Where("USER_ID = ?", id).
		Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var question []model.Question
	if err := db.Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID,"+
			"SCORE,"+
			"USER_ID,"+
			"question_type.ID         question_type_id,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"course.ID                course_id,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("USER_ID = ?", id).
		Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&question).
		Error; err != nil {
		return nil, 0, err
	}
	return question, count, nil
}

// ModifyQuestion 修改题目
func (d *Dao) ModifyQuestion(id int, question model.ModifyQuestion) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Model(&question).
		Where("ID = ?", id).
		Updates(question).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// RemoveQuestion 通过id移除题目
func (d *Dao) RemoveQuestion(id int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID = ?", id).
		Update("REMOVED", 1).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// BatchRemoveQuestion 通过id批量移除题目
func (d *Dao) BatchRemoveQuestion(ids []int) error {
	db := global.DBEngine
	tx := db.Begin()
	if err := tx.Error; err != nil {
		logger.Errorf("事务开启异常: %v\n", err)
		return err
	}

	if err := tx.Table("question_info").
		Where("ID in (?)", ids).
		Update("REMOVED", 1).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// QueryQuestion 综合所有条件查询
func (d *Dao) QueryQuestion(score, min, max *float64,
	typeId, difficultyId, courseId *int,
	questionContent, answerContent *string,
	pageNum, pageSize int, isRemoved bool) ([]model.Question, int, error) {
	db := global.DBEngine
	var questionList []model.Question
	db = db.Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID," +
			"SCORE," +
			"USER_ID," +
			"question_type.ID         question_type_id," +
			"question_difficulty.ID   question_difficulty_id," +
			"course.ID                course_id," +
			"QUESTION_MD," +
			"QUESTION_TXT," +
			"QUESTION_HTML," +
			"ANSWER_MD," +
			"ANSWER_TXT," +
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty " +
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type " +
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course " +
			"on course.ID = question_info.COURSE_ID")
	if score != nil {
		db = db.Where("SCORE = ?", *score)
	} else {
		if min == nil && max != nil {
			db = db.Where("SCORE <= ?", *max)
		} else if min != nil && max == nil {
			db = db.Where("SCORE >= ?", *min)
		} else if min != nil && max != nil {
			db = db.Where("SCORE >= ? && SCORE <= ?", *min, *max)
		}
	}
	if typeId != nil {
		db = db.Where("question_type.ID = ?", typeId)
	}
	if difficultyId != nil {
		db = db.Where("question_difficulty.ID = ?", difficultyId)
	}
	if courseId != nil {
		db = db.Where("course.ID = ?", courseId)
	}
	if questionContent != nil {
		db = db.Where("MATCH(QUESTION_TXT) AGAINST(? IN BOOLEAN MODE)", *questionContent)
	}
	if answerContent != nil {
		db = db.Where("MATCH(ANSWER_TXT) AGAINST(? IN BOOLEAN MODE)", *answerContent)
	}
	db = db.Where("REMOVED = ?", isRemoved)

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&questionList)
	if err := db.Error; err != nil {
		return nil, count, err
	}

	return questionList, count, nil
}

// PreciseQueryQuestionByScore 仅凭分值精确查询
func (d *Dao) PreciseQueryQuestionByScore(score float64, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_info").
		Where("SCORE = ?", score).
		Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	if err := db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID,"+
			"SCORE, "+
			"USER_ID, "+
			"question_type.ID         question_type_id,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"course.ID                course_id,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("SCORE = ?", score).
		Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// IntervalQueryQuestionByScore 仅凭分值区间查询
func (d *Dao) IntervalQueryQuestionByScore(min, max float64, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	db = db.Table("question_info").Where("REMOVED = ?", isRemoved)

	if min < 0 && max >= 0 {
		db = db.Where("SCORE <= ?", max)
	} else if min >= 0 && max < 0 {
		db = db.Where("SCORE >= ?", min)
	} else if min >= 0 && max >= 0 {
		db = db.Where("SCORE >= ? and SCORE <= ?", min, max)
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db = db.Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID," +
			"SCORE, " +
			"USER_ID, " +
			"question_type.ID         question_type_id," +
			"question_difficulty.ID   question_difficulty_id," +
			"course.ID                course_id," +
			"QUESTION_MD," +
			"QUESTION_TXT," +
			"QUESTION_HTML," +
			"ANSWER_MD," +
			"ANSWER_TXT," +
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty " +
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type " +
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course " +
			"on course.ID = question_info.COURSE_ID")

	if min < 0 && max >= 0 {
		db = db.Where("SCORE <= ?", max)
	} else if min >= 0 && max < 0 {
		db = db.Where("SCORE >= ?", min)
	} else if min >= 0 && max >= 0 {
		db = db.Where("SCORE >= ? and SCORE <= ?", min, max)
	}

	var questionList []model.Question
	db = db.Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList)

	if err := db.Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// QueryQuestionByType 通过题型查询
func (d *Dao) QueryQuestionByType(typeId int, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	db = db.Table("question_info")
	if typeId > 0 {
		db = db.Where("TYPE_ID = ?", typeId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	db = db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID," +
			"SCORE, " +
			"USER_ID, " +
			"question_type.ID         question_type_id, " +
			"question_difficulty.ID   question_difficulty_id, " +
			"course.ID                course_id, " +
			"QUESTION_MD," +
			"QUESTION_TXT," +
			"QUESTION_HTML," +
			"ANSWER_MD," +
			"ANSWER_TXT," +
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty " +
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type " +
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course " +
			"on course.ID = question_info.COURSE_ID")
	if typeId > 0 {
		db = db.Where("TYPE_ID = ?", typeId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// QueryQuestionByDifficulty 通过难度查询
func (d *Dao) QueryQuestionByDifficulty(difficultyId int, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	db = db.Table("question_info")
	if difficultyId > 0 {
		db = db.Where("DIFFICULTY_ID = ?", difficultyId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	db = db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID," +
			"SCORE, " +
			"USER_ID, " +
			"question_type.ID         question_type_id, " +
			"question_difficulty.ID   question_difficulty_id, " +
			"course.ID                course_id, " +
			"QUESTION_MD," +
			"QUESTION_TXT," +
			"QUESTION_HTML," +
			"ANSWER_MD," +
			"ANSWER_TXT," +
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty " +
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type " +
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course " +
			"on course.ID = question_info.COURSE_ID")
	if difficultyId > 0 {
		db = db.Where("DIFFICULTY_ID = ?", difficultyId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// QueryQuestionByCourse 通过所属课程查询
func (d *Dao) QueryQuestionByCourse(courseId int, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	db = db.Table("question_info")
	if courseId > 0 {
		db = db.Where("COURSE_ID = ?", courseId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	db = db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID," +
			"SCORE," +
			"USER_ID," +
			"question_type.ID         question_type_id," +
			"question_difficulty.ID   question_difficulty_id," +
			"course.ID                course_id," +
			"QUESTION_MD," +
			"QUESTION_TXT," +
			"QUESTION_HTML," +
			"ANSWER_MD," +
			"ANSWER_TXT," +
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty " +
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type " +
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course " +
			"on course.ID = question_info.COURSE_ID")
	if courseId > 0 {
		db = db.Where("COURSE_ID = ?", courseId)
	}
	if err := db.Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// QueryQuestionByQuestionContent 通过题目内容查询
func (d *Dao) QueryQuestionByQuestionContent(questionContent string, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_info").
		Where("MATCH(QUESTION_TXT) AGAINST(? IN BOOLEAN MODE)", questionContent).
		Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	if err := db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID,"+
			"SCORE,"+
			"USER_ID,"+
			"question_type.ID         question_type_id,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"course.ID                course_id,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("MATCH(QUESTION_TXT) AGAINST(? IN BOOLEAN MODE)", questionContent).
		Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}

// QueryQuestionByAnswerContent 通过答案内容查询
func (d *Dao) QueryQuestionByAnswerContent(answerContent string, isRemoved bool, pageNum, pageSize int) ([]model.Question, int, error) {
	db := global.DBEngine
	var count int
	if err := db.Table("question_info").
		Where("MATCH(ANSWER_TXT) AGAINST(? IN BOOLEAN MODE)", answerContent).
		Where("REMOVED = ?", isRemoved).
		Count(&count).
		Error; err != nil {
		return nil, 0, err
	}

	var questionList []model.Question
	if err := db.
		Table("question_info").
		Preload("QuestionType").
		Preload("QuestionDifficulty").
		Preload("Course").
		Select("question_info.ID,"+
			"SCORE,"+
			"USER_ID,"+
			"question_type.ID         question_type_id,"+
			"question_difficulty.ID   question_difficulty_id,"+
			"course.ID                course_id,"+
			"QUESTION_MD,"+
			"QUESTION_TXT,"+
			"QUESTION_HTML,"+
			"ANSWER_MD,"+
			"ANSWER_TXT,"+
			"ANSWER_HTML").
		Joins("left join question_difficulty_info question_difficulty "+
			"on question_info.DIFFICULTY_ID = question_difficulty.ID").
		Joins("left join question_type_info question_type "+
			"on question_info.TYPE_ID = question_type.ID").
		Joins("left join course_info course "+
			"on course.ID = question_info.COURSE_ID").
		Where("MATCH(ANSWER_TXT) AGAINST(? IN BOOLEAN MODE)", answerContent).
		Where("REMOVED = ?", isRemoved).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&questionList).
		Error; err != nil {
		return nil, 0, err
	}

	return questionList, count, nil
}
