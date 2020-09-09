package log

// 日志文件配置
type Config struct {
	Writer string
	Format string
	Level  string
	File   *ConfigFile
}

type ConfigFile struct {
	Name             string
	Path             string
	RotationOpen     bool
	RotationUnit     string
	RotationTime     int
	rotationTimeSave int
}

func getConfig() *Config {
	return &Config{
		Writer: "file,stderr",
		Format: "json",
		Level:  "debug",
		File: &ConfigFile{
			Name:             "api",
			Path:             "./log",
			RotationOpen:     false,
			RotationUnit:     "hour",
			RotationTime:     1,
			rotationTimeSave: 6,
		},
	}
}
