package main

import (
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
		Addr:           ":" + global.ServerSetting.HttpPort,
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

	sets, _ := setting.NewSetting()
	_ = sets.ReadSection("Server", &global.ServerSetting)
	_ = sets.ReadSection("App", &global.AppSetting)
	_ = sets.ReadSection("Database", &global.DatabaseSetting)
	//_ = sets.ReadSection("JWT", &global.JWTSetting)
	_ = sets.ReadSection("Email", &global.EmailSetting)
	_ = sets.ReadSection("Redis", &global.RedisSetting)

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	//global.JWTSetting.Expire *= time.Second
}

func setupLogger() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("setupLogger() error! : %v", r)
		}
	}()

	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
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
	var db = global.DBEngine

	pubsub := global.RedisClient.Subscribe(global.RedisSetting.Topic)
	_, err := pubsub.Receive()
	if err != nil {
		global.Logger.Errorf("Subscribe error: %v", err)
	}

	var adminList []model.User
	if err := db.Table("user_info").
		Where("IS_ADMINISTRATOR = ?", 1).
		Order("ID").
		Find(&adminList).
		Error; err != nil {
		global.Logger.Errorf("获取管理员列表出错: %v", err)
	}
	var adminEmailList []string
	adminEmailList = make([]string, 0, len(adminList))
	for _, admin := range adminList {
		adminEmailList = append(adminEmailList, *admin.Email)
	}

	var newUser *model.User
	for msg := range pubsub.Channel() {
		err := json.Unmarshal([]byte(msg.Payload), &newUser)
		if err != nil {
			global.Logger.Errorf("json.Unmarshal error: %v", err)
		}

		for _, email_ := range adminEmailList {
			err = global.Email.SendMail(
				[]string{email_},
				fmt.Sprintf("TPCS：您有一个注册用户待审核"),
				fmt.Sprintf("新用户 <strong>%v</strong> 发起了注册请求，请尽快前往 <strong>TPCS & 教师管理</strong> 进行审核", *newUser.Username),
			)
			if err != nil {
				global.Logger.Errorf("向管理员[%v]发送有用户注册请求的邮件失败！error: %v\n", email_, err)
			}
		}
	}
}
