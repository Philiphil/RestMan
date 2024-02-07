package apiman

import (
	"github.com/gin-gonic/gin"
	method_type "github.com/philiphil/apiman/method/MethodType"
)

func (r *ApiRouter[T]) Options(c *gin.Context) {
	allowed := ""
	for _, method := range r.Methods {
		switch method.Method {
		case method_type.GetList:
			break
		case method_type.PutList:
			break
		case method_type.PatchList:
			break
		case method_type.DeleteList:
			break
		default:
			if len(allowed) > 0 {
				allowed += ","
			}
			allowed += method.Method.String()
		}
	}
	c.Header("Allow", allowed)
	c.Header("Content-Length", "0")
	c.Status(200)
}
