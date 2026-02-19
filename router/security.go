package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/security"
)

// FirewallCheck executes all configured firewalls to authenticate and retrieve the current user.
func (r *ApiRouter[T]) FirewallCheck(c *gin.Context) (security.User, error) {
	var user security.User
	var err error
	for _, firewall := range r.Firewalls {
		user, err = firewall.GetUser(c)
		if err != nil { //problem
			if apiErr, ok := err.(errors.ApiError); ok {
				if apiErr.Blocking {
					return user, err
				} else {
					continue
				}
			} else {
				// If the error is not an ApiError, we consider it as a blocking error by default
				return user, apiErr
			}
		} else { //user found !
			return user, nil
		}
	}
	// no issue in the end, which means no firewall or no Authentication
	// no firewall is not an error, it just means we shouldnt block the request
	// no auth ... could be a problem but its the responsibility of the authorization function to check if the user is authenticated or not, so we just return nil and let the authorization function handle it
	//lets say theres no auth function cause despiste having firewall for some resources, an entity doesnt require auth
	// we need to make this non blocking to allow the request to go through, and let the authorization function handle the case where the user is not authenticated but the resource requires authentication
	return nil, nil
}

// ReadingCheck verifies that the authenticated user has permission to read the specified object.
func (r *ApiRouter[T]) ReadingCheck(c *gin.Context, object *T) error {
	user, err := r.FirewallCheck(c)
	if err != nil {
		return errors.Unauthorized
	}
	rr, ok := security.HasReadingRights(*object)
	if ok {
		auth := rr.GetReadingRights()
		if !auth(user, *object) {
			return errors.Unauthorized
		}
	}

	return nil
}

// WritingCheck verifies that the authenticated user has permission to modify the specified object.
func (r *ApiRouter[T]) WritingCheck(c *gin.Context, object *T) error {
	user, err := r.FirewallCheck(c)
	if err != nil {
		return err
	}
	rr, ok := security.HasWritingRights(*object)
	if ok {
		auth := rr.GetWritingRights()
		if !auth(user, *object) {
			return errors.Unauthorized
		}
	}

	return nil
}
