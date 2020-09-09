package config

import (
	"github.com/spf13/viper"
	"reflect"
)

type I interface {
	SetDefault(key string, value string)             //设置默认值，最低优先级
	Set(key string, value interface{})               //设置配置 原有将覆盖
	BindEnv(key string, envName string)              //绑定环境变量
	GetBool(key string) bool                         // 获取 bool
	GetFloat64(key string) float64                   // 获取 Float64
	GetInt(key string) int                           // 获取Int
	GetIntSlice(key string) []int                    // 获取 []int
	GetString(key string) string                     // 获取 string
	GetStringMap(key string) map[string]interface{}  // 获取 map[string]interface{}
	GetStringMapString(key string) map[string]string // 获取 map[string]string
	GetStringSlice(key string) []string              // 获取 []string
	IsSet(key string) bool                           // 判断配置是否存在
	AllSettings() map[string]interface{}             //获取全部配置文件
	//获取配置信息 可以通过Struct和字符串获取  备注：具体的Struct 写在configStruct.go
	Get(rawVal interface{}, key ...string) (err error)
}
type Config struct {
	I
}
type rxViper struct{}

func New() *Config {
	return &Config{
		I: rxViper{},
	}
}

// 设置默认值，最低优先级
func (c rxViper) SetDefault(key string, value string) {
	viper.SetDefault(key, value)
}

// 设置配置 将覆盖
func (c rxViper) Set(key string, value interface{}) {
	viper.Set(key, true)
}

// 绑定环境变量
func (c rxViper) BindEnv(key string, envName string) {
	viper.BindEnv(key, envName)
}

//获取配置信息
//可以通过Struct和字符串获取
//备注：具体的Struct 写在rxViperStruct.go
func (c rxViper) Get(rawVal interface{}, key ...string) (err error) {
	if rawVal == nil {
		return
	}
	if vType := reflect.TypeOf(rawVal); vType.Name() == "string" {
		viper.GetInt("winter.Age")
		return
	}
	// 读取键名 示例：mysql
	if key != nil {
		viper.UnmarshalKey(key[0], &rawVal)
	}
	if key == nil {
		viper.Unmarshal(&rawVal)
	}
	return
}

// 获取bool
func (c rxViper) GetBool(key string) bool {
	return viper.GetBool(key)
}

// 获取Float64
func (c rxViper) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

// 获取Int
func (c rxViper) GetInt(key string) int {
	return viper.GetInt(key)
}

// 获取 []int
func (c rxViper) GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

// 获取 string
func (c rxViper) GetString(key string) string {
	return viper.GetString(key)
}

// 获取 map[string]interface{}
func (c rxViper) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// 获取 map[string]string
func (c rxViper) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// 获取 []string
func (c rxViper) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// 判断配置是否存在
func (c rxViper) IsSet(key string) bool {
	return viper.IsSet(key)
}

//获取全部配置文件
func (c rxViper) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}

// 配置文件初始化读取
func Init(file string) {
	openFile(file)
}

// 打开配置文件
func openFile(configFile string) (err error) {
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err = viper.ReadInConfig() // Find and read the config file
	return
}
