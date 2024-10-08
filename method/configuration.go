package method

import method_type "github.com/philiphil/restman/method/MethodType"

type ApiMethodConfiguration struct {
	Method              method_type.ApiMethod
	SerializationGroups []string
}
