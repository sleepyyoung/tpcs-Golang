package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/internal/routers"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/email"
	"tpcs/pkg/logger"
	"tpcs/pkg/setting"
	"tpcs/pkg/upload"
)

func init() {
	// 存session用
	gob.Register(model.User{})
	gob.Register(upload.UploadStatus{})
	setupSetting()
	setupLogger()
	setupEmail()
	setupDBEngine()
	setupRedisClient()
	go registerRequestListener()
}

// CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tpcs .
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("项目启动失败！: %v\n", r)
		}
	}()

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupSetting() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("项目初始化失败！: %v\n", r)
		}
	}()

	sets, err := setting.NewSetting(setting.ENV_DOCKER)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("App", &global.AppSetting)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("Redis", &global.RedisSetting)
	if err != nil {
		panic(err)
	}
	err = sets.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		panic(err)
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
}

func setupLogger() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("日志组件初始化失败！: %v\n", r)
		}
	}()

	err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
}

func setupEmail() {
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
			log.Fatalf("数据库初始化失败！: %v\n", r)
		}
	}()

	engine, err := model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		panic(err)
	}
	global.DBEngine = engine
}

func setupRedisClient() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Redis初始化失败！: %v\n", r)
		}
	}()

	client := redis.NewClient(
		&redis.Options{
			Addr:     global.RedisSetting.Addr,
			Password: global.RedisSetting.Password,
			DB:       global.RedisSetting.DB,
		},
	)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	global.RedisClient = client
}

// Redis订阅Topic:register（注册用）
func registerRequestListener() {
	userSvc := userService.New(context.Background())

	pubsub := global.RedisClient.Subscribe(global.RedisSetting.Topic)
	_, err := pubsub.Receive()
	if err != nil {
		logger.Errorf("global.RedisClient.Subscribe error: %v\n", err)
	}

	var adminList []model.User
	adminList, err = userSvc.GetAdminList()
	if err != nil {
		logger.Errorf("userSvc.GetAdminList error: %v\n", err)
		return
	}

	var newUser0 *model.User
	for msg := range pubsub.Channel() {
		err := json.Unmarshal([]byte(msg.Payload), &newUser0)
		if err != nil {
			logger.Errorf("json.Unmarshal error: %v\n", err)
		}

		newUser, err := userSvc.GetUserByUsername(*newUser0.Username)
		if err != nil {
			logger.Errorf("userSvc.GetUserByUsername error: %v\n", err)
			return
		}

		for _, admin := range adminList {
			token, err := app.GenerateToken(
				global.JWTSetting.Key,
				global.JWTSetting.Secret,
				*admin.Username,
				*newUser.Username,
			)
			if err != nil {
				logger.Errorf("[admin: %v]app.GenerateToken error: %v\n", *admin.Username, err)
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
				logger.Errorf("向管理员[%v]发送有用户注册请求的邮件失败！error: %v\n", *admin.Email, err)
			}
		}
	}
}
