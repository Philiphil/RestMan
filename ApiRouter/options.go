package ApiRouter

import (
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Options(c *gin.Context) {
	c.Header("Allow", "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Length", "0")
	c.Status(200)
}
