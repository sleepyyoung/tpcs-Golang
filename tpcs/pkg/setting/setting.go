package setting

import (
	"github.com/spf13/viper"
	"log"
)

type Setting struct {
	vp *viper.Viper
}

const (
	ENV_DEV  = "config-dev"
	ENV_TEST = "config-test"
	ENV_PROD = "config-prod"
)

func NewSetting(env string) (*Setting, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("NewSetting error! : %v", r)
		}
	}()

	vp := viper.New()
	vp.SetConfigName(env)
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	vp.ReadInConfig()

	return &Setting{vp}, nil
}
