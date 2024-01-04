package security

import "github.com/gin-gonic/gin"

type Firewall interface {
	GetUser(c *gin.Context) (IUser, error)
}
