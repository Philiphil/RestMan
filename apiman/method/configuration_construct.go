package method

import (
	"github.com/philiphil/apiman/apiman/method/MethodType"
)

func New() ApiMethodConfiguration {
	return ApiMethodConfiguration{
		Method:              method_type.Undefined,
		SerializationGroups: []string{},
	}

}

func Method(method method_type.ApiMethod, groups ...string) ApiMethodConfiguration {
	c := New()
	c.Method = method
	c.SerializationGroups = groups
	return c
}

func DefaultApiMethods() (d []ApiMethodConfiguration) {
	d = append(d,
		Method(method_type.Get),
		Method(method_type.GetList),
		Method(method_type.Head),
		Method(method_type.Put),
		Method(method_type.Patch),
		Method(method_type.Post),
		Method(method_type.Options),
		Method(method_type.Connect),
		Method(method_type.Trace),
	)
	return d
}
