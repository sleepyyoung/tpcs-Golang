package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

const (
	ENV_DEV  = "config-dev"
	ENV_PROD = "config-prod"
)

func NewSetting(env string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName(env)
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}
