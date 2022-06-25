package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/internal/routers"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/email"
	"tpcs/pkg/file"
	"tpcs/pkg/logger"
	"tpcs/pkg/setting"
)

func init() {
	// 存session用
	gob.Register(model.User{})
	gob.Register(file.UploadStatus{})
	setupSetting()
	setupLogger()
	setupEmail()
	setupDBEngine()
	setupRedisClient()
	go registerRequestListener()
}

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr: ":" +
			global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

func setupSetting() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupSetting() error! : %v", r)
		}
	}()

	sets, _ := setting.NewSetting(setting.ENV_DEV)
	_ = sets.ReadSection("Server", &global.ServerSetting)
	_ = sets.ReadSection("App", &global.AppSetting)
	_ = sets.ReadSection("Database", &global.DatabaseSetting)
	_ = sets.ReadSection("Email", &global.EmailSetting)
	_ = sets.ReadSection("Redis", &global.RedisSetting)
	_ = sets.ReadSection("JWT", &global.JWTSetting)

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
}

func setupLogger() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupLogger() error! : %v", r)
		}
	}()

	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath +
			"/" +
			global.AppSetting.LogFileName +
			global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	global.Logger = global.Logger.WithCallersFrames()
}

func setupEmail() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupEmail() error! : %v", r)
		}
	}()

	global.Email = email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
}

func setupDBEngine() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupDBEngine() error! : %v", r)
		}
	}()

	global.DBEngine, _ = model.NewDBEngine(global.DatabaseSetting)
}

func setupRedisClient() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupRedisClient() error! : %v", r)
		}
	}()

	global.RedisClient, _ = model.NewRedisClient(global.RedisSetting)
}

// Redis订阅Topic:register（注册用）
func registerRequestListener() {
	userSvc := userService.New(context.Background())

	pubsub := global.RedisClient.Subscribe(global.RedisSetting.Topic)
	_, err := pubsub.Receive()
	if err != nil {
		global.Logger.Errorf("global.RedisClient.Subscribe error: %v\n", err)
	}

	var adminList []model.User
	adminList, err = userSvc.GetAdminList()
	if err != nil {
		return
	}

	var newUser0 *model.User
	for msg := range pubsub.Channel() {
		err := json.Unmarshal([]byte(msg.Payload), &newUser0)
		if err != nil {
			global.Logger.Errorf("json.Unmarshal error: %v"+
				"", err)
		}

		newUser, err := userSvc.GetUserByUsername(*newUser0.Username)
		if err != nil {
			global.Logger.Errorf("userSvc.GetUserByUsername error: %v", err)
			return
		}

		for _, admin := range adminList {
			token, err := app.GenerateToken(
				global.JWTSetting.Key,
				global.JWTSetting.Secret,
				*admin.Username,
				*newUser0.Username,
			)
			if err != nil {
				global.Logger.Errorf("[admin: %v]app.GenerateToken error: %v\n", *admin.Username, err)
				return
			}

			err = global.Email.SendMail(
				[]string{*admin.Email},
				fmt.Sprintf("TPCS：您有一个注册用户待审核"),
				fmt.Sprintf("您收到了新的注册请求：<br/>用户名：%v<br/>邮箱：%v<br/>备注：<br/><textarea disabled>%v</textarea><br/>请（在校园网环境下）点击链接确认是否审核通过：<a href=\"%vteacher-audit?token=%v\">%vteacher-audit?token=%v</a>，<br/>或请到 TPCS -> 教师管理 页面进行审核。",
					*newUser.Username,
					*newUser.Email,
					*newUser.Note,
					global.AppSetting.URL,
					token,
					global.AppSetting.URL,
					token,
				),
			)
			if err != nil {
				global.Logger.Errorf("向管理员[%v]发送有用户注册请求的邮件失败！error: %v\n", *admin.Email, err)
			}
		}

	}
}
