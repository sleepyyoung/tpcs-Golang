package user

import (
	"github.com/jinzhu/gorm"
	"tpcs/global"
	"tpcs/internal/pojo/model"
)

// GetUserByUserId 通过用户id查询用户
func (d *Dao) GetUserByUserId(id int) (*model.User, error) {
	var db = global.DBEngine
	var user model.User
	if err := db.Raw("SELECT * FROM USER_INFO WHERE ID = ?", id).Scan(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// GetUserByUsernameAndPassword 通过用户名和密码查询用户
func (d *Dao) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	var db = global.DBEngine
	var user model.User
	if err := db.Raw("SELECT * FROM USER_INFO WHERE USERNAME = ? AND PASSWORD = ?", username, password).Scan(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// GetUserByUsername 通过用户名查询用户
func (d *Dao) GetUserByUsername(username string) (*model.User, error) {
	var db = global.DBEngine
	var user model.User
	if err := db.Raw("SELECT * FROM USER_INFO WHERE USERNAME = ?", username).Scan(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// GetUserByEmail 通过邮箱查询用户
func (d *Dao) GetUserByEmail(email string) (*model.User, error) {
	var db = global.DBEngine
	var user model.User
	if err := db.Raw("SELECT * FROM USER_INFO WHERE EMAIL = ?", email).Scan(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// GetUserById 通过id查询用户
func (d *Dao) GetUserById(id int) (*model.User, error) {
	var db = global.DBEngine
	var user model.User
	if err := db.Raw("SELECT * FROM USER_INFO WHERE ID = ?", id).Scan(&user).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// DeleteUserById 通过id删除用户
func (d *Dao) DeleteUserById(id int) (bool, error) {
	var db = global.DBEngine
	if err := db.Exec("DELETE FROM USER_INFO WHERE ID = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}

//// IsExistsUser 判断用户是否存在
//func (d *Dao) IsExistsUser(username, password string) (exists bool, err error) {
//	user := model.User{Username: username, Password: password}
//	return user.IsExists(d.engine)
//}

// IsAdminByUsernameAndPassword 根据用户名和密码判断用户是否为管理员
func (d *Dao) IsAdminByUsernameAndPassword(username, password string) (bool, error) {
	user, err := d.GetUserByUsernameAndPassword(username, password)
	if err != nil {
		return false, err
	}
	return *user.IsAdministrator, nil
}

// GetAdminList 获取管理员列表
func (d *Dao) GetAdminList() ([]model.User, error) {
	var db = global.DBEngine
	var adminList []model.User
	if err := db.Table("user_info").
		Where("IS_ADMINISTRATOR = ?", 1).
		Order("ID").
		Find(&adminList).
		Error; err != nil {
		return nil, err
	}
	return adminList, nil
}

// IsAdminUserByUserId 根据用户id判断用户是否为管理员
func (d *Dao) IsAdminUserByUserId(db *gorm.DB, id int) (bool, error) {
	var count int
	if err := db.Table("user_info").
		Select("IS_ADMINISTRATOR").
		Where("ID = ?", id).
		Where("IS_ADMINISTRATOR = ?", true).
		Count(&count).Error; err != nil {
		return false, nil
	}
	return count == 1, nil
}

//// ModifyUserPassword 修改用户密码
//func (d *Dao) ModifyUserPassword(user *model.User, newPassword string) (exists bool, err error) {
//	return user.ModifyPassword(d.engine, newPassword)
//}

// ModifyUser 修改用户
func (d *Dao) ModifyUser(user *model.User) (exists bool, err error) {
	return user.Update(d.engine)
}

// CreateUser 创建用户
func (d *Dao) CreateUser(user *model.User) error {
	return user.Create(d.engine)
}

// UserCount 获取用户数量
func (d *Dao) UserCount(db *gorm.DB, isAdmin *bool) (int, error) {
	var count int
	db = db.Table("user_info")
	if isAdmin != nil {
		db = db.Where("IS_ADMINISTRATOR = ?", *isAdmin)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// NotAdminUserList 获取非admin用户列表
func (d *Dao) NotAdminUserList(db *gorm.DB, pageNum, pageSize int) ([]model.User, int, error) {
	var isAdmin = new(bool)
	*isAdmin = false

	count, err := d.UserCount(db, isAdmin)
	if err != nil {
		return nil, 0, err
	}

	var teacherList []model.User
	if err := db.Table("user_info").
		Where("IS_ADMINISTRATOR = ?", *isAdmin).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&teacherList).
		Error; err != nil {
		return nil, 0, err
	}

	return teacherList, count, nil
}
