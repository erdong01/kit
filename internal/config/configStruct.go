package config

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
		Addr:     "101.132.145.239:6379",
		Password: "m7ictwC7b0xnc3H8",
		DB:       1,
	}
}

func GetJwtCnf() *JwtCnf {
	return &JwtCnf{
		Secret: "Y1HH999gmSmGvVddzmhGhThiOEBQBPbj",
		Issuer: "http://student-2c-api.517rxt.test/v1/auth/login",
		Expire: 10,
	}
}

// jwt 配置
type JwtCnf struct {
	Secret string
	Issuer string
	Expire int64
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
