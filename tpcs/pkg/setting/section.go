package setting

import "time"

type AppSettingS struct {
	URL             string
	LogSavePath     string
	LogFileName     string
	LogFileExt      string
	UploadDir       string
	FileUploadPath  string
	MdImgUploadPath string
	SessionNames    []string
}

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseSettingS struct {
	DBType       string
	Username     string
	Password     string
	Host         string
	DBName       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}

type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

type RedisSettingS struct {
	Addr     string
	Password string
	DB       int
	Topic    string
}

type JWTSettingS struct {
	Key    string
	Secret string
	Issuer string
	Expire time.Duration
}
