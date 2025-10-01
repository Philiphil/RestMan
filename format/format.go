package format

type Format string

// Const of default formats handled by RestMan
const (
	Undefined Format = "undefined"
	Unknown          = "unknown"
	JSON             = "application/json"
	JSONLD           = "application/ld+json"
	XML              = "text/xml"
	CSV              = "application/csv"
)
