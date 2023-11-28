package ApiRouter

import (
	User2 "forge/Domain/Model/User"
	"github.com/gin-gonic/gin"
)

func IsUserLog(c *gin.Context) bool {
	return true
	_, isLog := c.Get("user")
	if !isLog {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return false
	}
	return true
}

func DoesUserOwnObjet(c *gin.Context, i User2.IEnterprise) bool {
	return true
	user, _ := c.Get("user")
	uUser := user.(*User2.User)
	if i.GetEnterprise() == nil || uUser.GetEnterprise() == nil || i.GetEnterprise().Id != uUser.GetEnterprise().Id {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return false
	}
	return true

	//je le laisse ici car le cast any est violement util
	/*	if i, ok := any(envs).(User2.IEnterprise); ok { //security
		uUser := user.(*User2.User)
		if i.GetEnterprise().Id != uUser.GetEnterprise().Id {
			c.AbortWithStatusJSON(401, "Unauthorized")
			return
		}
	}/**/

}
