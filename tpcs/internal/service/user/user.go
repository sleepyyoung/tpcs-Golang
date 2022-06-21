package user

import (
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
)

// GetUserByUserId 通过用户id查询用户
func (svc *Service) GetUserByUserId(id int) (*model.User, error) {
	return svc.dao.GetUserByUserId(id)
}

// GetUserByUsernameAndPassword 通过用户名和密码查询用户
func (svc *Service) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	return svc.dao.GetUserByUsernameAndPassword(username, password)
}

// GetUserByUsername 通过用户名查询用户
func (svc *Service) GetUserByUsername(username string) (*model.User, error) {
	return svc.dao.GetUserByUsername(username)
}

// GetUserByEmail 通过邮箱查询用户
func (svc *Service) GetUserByEmail(email string) (*model.User, error) {
	return svc.dao.GetUserByEmail(email)
}

// GetUserById 通过id查询用户
func (svc *Service) GetUserById(id int) (*model.User, error) {
	return svc.dao.GetUserById(id)
}

// DeleteUserById 通过id删除用户
func (svc *Service) DeleteUserById(id int) (bool, error) {
	return svc.dao.DeleteUserById(id)
}

// IsAdminByUsernameAndPassword 根据用户名和密码判断用户是否为管理员
func (svc *Service) IsAdminByUsernameAndPassword(username, password string) (bool, error) {
	return svc.dao.IsAdminByUsernameAndPassword(username, password)
}

// GetAdminList 获取管理员列表
func (svc *Service) GetAdminList() ([]model.User, error) {
	return svc.dao.GetAdminList()
}

// ModifyUser 修改用户信息
func (svc *Service) ModifyUser(user *model.User) (bool, error) {
	return svc.dao.ModifyUser(user)
}

// CreateUser 创建用户
func (svc *Service) CreateUser(u *model.User) error {
	return svc.dao.CreateUser(u)
}

// NotAdminUserList 获取非admin用户列表
func (svc *Service) NotAdminUserList(param *service.ListRequest) ([]model.User, int, error) {
	return svc.dao.NotAdminUserList(global.DBEngine, param.Page, param.Limit)
}
