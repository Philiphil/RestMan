package security

import "github.com/gin-gonic/gin"

type AuthProvider interface {
	GetUser(c *gin.Context) (IUser, error)
}
