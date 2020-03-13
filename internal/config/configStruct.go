package config

import "rxt/internal/upload/qiniu"

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
		Password: "",
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
		Issuer: "http://student-2c-api.517rxt.test/v1/auth/login",
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
		Host:        "127.0.0.1",
		Port:        3306,
		User:        "root",
		Password:    "123456",
		Database:    "test",
		Parameters:  "false",
		MaxIdleConn: 40,
		MaxOpenConn: 100,
	}
}


// GetQiniu 获取七牛云配置
func GetQiniu() qiniu.Config {
	return qiniu.Config{
		DomainDefault: "api.qiniu.com",
		DomainHTTPS:   "api-statics.qiniu.com",
		AccessKey:     "X8rnYilA5bdpOeJl1w1vaLx9J7AmfaXGYBaMSKqB",
		SecretKey:     "Dd8KGXY_hcVhq_6x8PFfPbJOB-ugyyn1paiTEYNh",
		Bucket:        "api-com",
		Access:        "public",
		CachePath:     "qiniu/cache/",
	}
}
