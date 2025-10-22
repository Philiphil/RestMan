package route

import "fmt"

// RouteType represents the type of HTTP route operation.
type RouteType int8

const (
	Undefined RouteType = iota
	//CRUD
	GetList
	Get
	Post
	Put
	Delete
	Patch

	//useless stuff
	Options
	Connect
	Trace
	Head

	//batch operations
	BatchGet
	BatchPost
	BatchPut
	BatchPatch
	BatchDelete
)

// String returns the HTTP method name for the RouteType.
func (e RouteType) String() string {
	switch e {
	case Patch:
		return "PATCH"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Get:
		return "GET"
	case Head:
		return "HEAD"
	case Delete:
		return "DELETE"
	case Options:
		return "OPTIONS"
	case Trace:
		return "TRACE"
	case Connect:
		return "CONNECT"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}
