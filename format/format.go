package format

type Format string

const (
	Undefined Format = "undefined"
	Unknown          = "unknown"
	JSON             = "JSON"
	XML              = "XML"
	CSV              = "CSV"
)
