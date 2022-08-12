package config

type Core struct {
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
		Uri: "mongodb://127.0.0.1:27017",
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
