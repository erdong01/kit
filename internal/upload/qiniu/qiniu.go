package qiniu

import (
	"strings"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

// const
const (
	DefaultExpire = 7200 // token超时时间：1小时
)

// Config 七牛云配置
type Config struct {
	DomainDefault string   // 默认域
	DomainHTTPS   string   // HTTPS域
	AccessKey     string   // 访问秘钥
	SecretKey     string   // 秘钥
	Bucket        string   // bucket
	NotifyURL     string   // 回调地址
	Access        string   // access 类型
	CachePath     string   // 缓存地址
	MimeType      []string // 允许的文件类型
	MaxSize       int64    // 允许最大文件
}

// Qiniu 七牛
type Qiniu struct {
	conf Config
}

// Upload 七牛上传
func (q *Qiniu) Upload(file string) error {

	// TODO
	return nil
}

// GetUploadToken 七牛上传凭证
func (q *Qiniu) GetUploadToken(path string, expire uint64) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope:           q.conf.Bucket + ":" + path,
		IsPrefixalScope: 1, // 允许用户上传以 scope 的 keyPrefix 为前缀的文件
		FsizeLimit:      q.conf.MaxSize,
		MimeLimit:       strings.Join(q.conf.MimeType, ";"),
		Expires:         expire,
	}

	mac := qbox.NewMac(q.conf.AccessKey, q.conf.SecretKey)

	upToken := putPolicy.UploadToken(mac)

	return upToken, nil
}

// GetDomain 获取访问域
func (q *Qiniu) GetDomain() string {
	return q.conf.DomainDefault
}

// GetMimeType 获取访问域
func (q *Qiniu) GetMimeType() []string {
	return q.conf.MimeType
}

// GetMaxSize 获取访问域
func (q *Qiniu) GetMaxSize() int64 {
	return q.conf.MaxSize
}

// New 返回七牛实例
func New(conf Config) *Qiniu {
	return &Qiniu{
		conf: conf,
	}
}
