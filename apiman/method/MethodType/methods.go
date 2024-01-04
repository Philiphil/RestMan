package method_type

import "fmt"

const (
	Undefined ApiMethod = iota
	//CRUD
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

	//Lists
	GetList
	PutList
	PatchList
	DeleteList
)

func (e ApiMethod) String() string {
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
