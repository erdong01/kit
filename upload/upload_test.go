package upload

import (
	"github.com/erDong01/micro-kit/upload/qiniu"
	"log"
	"testing"
)

func TestUploadToken(t *testing.T) {
	uploadHandle := NewUpload()

	path := "test/test"

	upToken, _ := uploadHandle.GetUploadToken(path, qiniu.DefaultExpire)
	if upToken == "" {
		t.Error("GetUploadToken() is false")
	}
	log.Println("Get uptoken is : " + upToken)
}
