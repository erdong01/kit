package upload

import (
	"log"
	"rxt/internal/upload/qiniu"
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
