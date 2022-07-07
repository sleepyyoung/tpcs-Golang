package global

import (
	"tpcs/pkg/email"
	"tpcs/pkg/setting"
)

var (
	AppSetting      *setting.AppSettingS
	ServerSetting   *setting.ServerSettingS
	DatabaseSetting *setting.DatabaseSettingS
	//Logger          *log.Logger
	Email        *email.Email
	EmailSetting *setting.EmailSettingS
	RedisSetting *setting.RedisSettingS
	JWTSetting   *setting.JWTSettingS
)
