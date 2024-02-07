package format

type Format string

const (
	Undefined Format = "undefined"
	Unknown          = "unknown"
	JSON             = "application/json"
	XML              = "text/xml"
	CSV              = "application/csv"
)
