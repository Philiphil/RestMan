package method

import (
	"github.com/philiphil/apiman/method/MethodType"
	"github.com/philiphil/apiman/security"
)

type ApiMethodConfiguration struct {
	Method                method_type.ApiMethod
	AuthorizationFunction security.AuthorizationFunction
	SerializationGroups   []string
}
