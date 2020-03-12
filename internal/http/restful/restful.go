package restful

import (
	"net/http"
	"rxt/internal/log"
	"rxt/internal/wrong"

	"github.com/gin-gonic/gin"
)

func Exception(c *gin.Context, err error) {
	errStruct := err.(*wrong.Err)
	s := errStruct.Format()
	errStruct.Trace = s
	log.Error(s)
	c.AbortWithStatusJSON(http.StatusExpectationFailed, errStruct)
}

func Success(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

func SuccessCreated(c *gin.Context, obj ...interface{}) {
	c.JSON(http.StatusCreated, obj)
}

func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
