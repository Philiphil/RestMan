package method_type

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
