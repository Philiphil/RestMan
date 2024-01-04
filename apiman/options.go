package apiman

import (
	"github.com/gin-gonic/gin"
)

func (r *ApiRouter[T]) Options(c *gin.Context) {
	allowed := ""
	for i, method := range r.Methods {
		if i > 0 {
			allowed += ","
		}
		allowed += method.Method.String()
	}
	c.Header("Allow", allowed)
	c.Header("Content-Length", "0")
	c.Status(200)
}
