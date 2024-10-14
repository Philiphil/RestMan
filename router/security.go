package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/security"
)

func (r *ApiRouter[T]) FirewallCheck(c *gin.Context) (security.IUser, error) {
	var user security.IUser
	var err error
	for _, firewall := range r.Firewalls {
		user, err = firewall.GetUser(c)
		if err != nil { //problem
			if err.(errors.ApiError).Blocking {
				return user, err
			} else {
				continue
			}
		} else { //user found !
			return user, nil
		}
	}
	return user, nil
}

func (r *ApiRouter[T]) ReadingCheck(c *gin.Context, object *T) error {
	user, err := r.FirewallCheck(c)
	if err != nil {
		return err
	}
	rr, ok := security.HasReadingRights(*object)
	if ok {
		auth := rr.GetReadingRights()
		if !auth(user, *object) {
			return errors.ErrUnauthorized
		}
	}

	return nil
}

func (r *ApiRouter[T]) WritingCheck(c *gin.Context, object *T) error {
	user, err := r.FirewallCheck(c)
	if err != nil {
		return err
	}
	rr, ok := security.HasWritingRights(*object)
	if ok {
		auth := rr.GetWritingRights()
		if !auth(user, *object) {
			return errors.ErrUnauthorized
		}
	}

	return nil
}
