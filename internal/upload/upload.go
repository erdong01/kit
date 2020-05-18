package upload

import (
	"github.com/erDong01/gin-kit/internal/config"
	"github.com/erDong01/gin-kit/internal/upload/qiniu"
)

var fileConfig = map[string]interface{}{
	"mime_type": []string{
		"image/bmp", "image/gif", "image/jpeg", "image/png", "image/nd.adobe.photoshop", // psd
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",   // docx
		"application/vnd.ms-powerpoint",                                             // ppt
		"application/vnd.openxmlformats-officedocument.presentationml.presentation", // pptx
		"application/vnd.android.package-archive",                                   // android app
		"application/octet-stream",                                                  // ios app
		"application/msword",                                                        // doc
		"audio/mpeg",                                                                // mp3
		"audio/mp3",                                                                 // mp3
	},
	"max_size": 1024 * 1024 * 50, // (Bytes) = 50mb
}

// IUpload 上传接口
type IUpload interface {
	GetUploadToken(path string, expire uint64) (string, error)
	Upload(file string) error
	GetDomain() string
	GetMimeType() []string
	GetMaxSize() int64
}

// Upload 上传文件
type Upload struct {
}

// NewUpload 上传
func NewUpload() IUpload {
	qiniuConf := config.GetQiniu()
	qiniuConf.MimeType = fileConfig["mime_type"].([]string)
	qiniuConf.MaxSize = int64(fileConfig["max_size"].(int))

	return qiniu.New(qiniuConf)
}
