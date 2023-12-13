package method

import (
	"github.com/philiphil/apiman/method/MethodType"
	"github.com/philiphil/apiman/security"
)

func New() ApiMethodConfiguration {
	return ApiMethodConfiguration{
		Method:                method_type.Undefined,
		AuthorizationFunction: security.NoAuthorizationRequired,
		SerializationGroups:   []string{},
	}

}

func Method(method method_type.ApiMethod, groups ...string) ApiMethodConfiguration {
	c := New()
	c.Method = method
	c.SerializationGroups = groups
	return c
}

func Security(c ApiMethodConfiguration, auth security.AuthorizationFunction) ApiMethodConfiguration {
	c.AuthorizationFunction = auth
	return c
}
