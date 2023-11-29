package method

type ApiMethod int8

const (
	Undefined ApiMethod = iota
	Get
	GetList
	Post
	Put
	Delete
	Patch

	//useless stuff
	Options
	Connect
	Trace
	Head

	//todo
	PutList
	PatchList
	DeleteList
)
