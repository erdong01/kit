package route

import (
	"github.com/gin-gonic/gin"
	"rxt/cmd/rxstudent/app/handler/ads"
)

func adsGroup(group *gin.RouterGroup) {
	r := group.Group("/ads")
	r.GET("/:PositionId", ads.Find)
}
