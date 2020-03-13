package ads

import (
	"github.com/gin-gonic/gin"
	ads "rxt/cmd/ads/proto"
	"rxt/internal/http/restful"
	"rxt/rpc/client/student/adsClient"
)

// 根据position_id获取广告
func Find(ctx *gin.Context) {
	var request ads.FindRequest
	if err := ctx.BindUri(&request); err != nil {
		restful.Exception(ctx, err)
		return
	}
	res, err := adsClient.New().Find(&request)
	if err != nil {
		restful.Exception(ctx, err)
		return
	}
	restful.Success(ctx, res)
}
