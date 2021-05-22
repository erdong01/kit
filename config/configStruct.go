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
		Secret: "",
		Issuer: "",
		Expire: 10,
		Ttl:    36000,
	}
}

// MySQL mysql数据库配置
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
		Password:    "",
		Database:    "",
		Parameters:  "false",
		MaxIdleConn: 40,
		MaxOpenConn: 100,
	}
}

type Mongo struct {
	Host []string
	Uri  string
}

// GetMongo
func GetMongo() Mongo {
	return Mongo{
		Uri: "mongodb://111.229.20.134:27017",
	}
}

// GetQiniu 获取七牛云配置
func GetQiniu() qiniu.Config {
	return qiniu.Config{
		DomainDefault: "",
		DomainHTTPS:   "",
		AccessKey:     "",
		SecretKey:     "",
		Bucket:        "",
		Access:        "",
		CachePath:     "",
	}
}

// mysql数据库配置
type Etcd struct {
	Addr []string
}

func GetEtcd() *Etcd {
	return &Etcd{
		Addr: []string{"127.0.0.1:2379"},
	}
}
