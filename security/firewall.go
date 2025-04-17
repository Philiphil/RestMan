package security

import "github.com/gin-gonic/gin"

// Firewall is an interface for the firewall
// using the provided Header in gin.Context, it should return an User or an error
type Firewall interface {
	GetUser(c *gin.Context) (User, error)
}
