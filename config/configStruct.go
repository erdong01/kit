package config

import "github.com/erDong01/micro-kit/upload/qiniu"

type Core struct {
}

// reids 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

//reids默认配置
func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "test",
		DB:       1,
	}
}

// jwt 配置
type JwtCnf struct {
	Secret string
	Issuer string
	Expire int64
	Ttl    int
}

func GetJwtCnf() *JwtCnf {
	return &JwtCnf{
		Secret: "Y1HH999gmSmGvVddzmhGhThiOEBQBPbj",
		Issuer: "",
		Expire: 10,
		Ttl:    36000,
	}
}

// mysql数据库配置
type MySQL struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	Parameters  string
	MaxIdleConn int
	MaxOpenConn int
}

func GetMySQL() *MySQL {
	return &MySQL{
		Host:        "101.132.145.239",
		Port:        33057,
		User:        "api.rxt",
		Password:    "a18@990Wll0C3",
		Database:    "api.517rxt.com",
		Parameters:  "false",
		MaxIdleConn: 40,
		MaxOpenConn: 100,
	}
}


// GetQiniu 获取七牛云配置
func GetQiniu() qiniu.Config {
	return qiniu.Config{
		DomainDefault: "api-statics.renxuetang.com",
		DomainHTTPS:   "api-statics.renxuetang.com",
		AccessKey:     "X8rnYilA5bdpOeJl1w1vaLx9J7AmfaXGYBaMSKqB",
		SecretKey:     "Dd8KGXY_hcVhq_6x8PFfPbJOB-ugyyn1paiTEYNh",
		Bucket:        "api-renxuetang-com",
		Access:        "public",
		CachePath:     "rxedu_test/cache/",
	}
}
