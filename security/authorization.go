// This package contains the authorization functions that are used to check if a user is allowed to access a resource
package security

import (
	"github.com/philiphil/restman/orm/entity"
)

// AuthorizationFunction is a signature for the authorization function
// it takes a user and an entity and returns a boolean
type AuthorizationFunction func(User, entity.Entity) bool

// AuthenticationRequired is a default implementation of the AuthorizationFunction
// it always returns true because it is reached only if the user is authenticated
// AuthenticationRequired validates that a user is authenticated before accessing an entity.
func AuthenticationRequired(user User, object entity.Entity) bool {
	return true
}
