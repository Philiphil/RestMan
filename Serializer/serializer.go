package serializer

import "github.com/philiphil/apiman/Serializer/Format"

// Serializer is the main serializer struct
type Serializer struct {
	Format Format.Format
}

// NewSerializer creates a new instance of Serializer
func NewSerializer(format Format.Format) *Serializer {
	return &Serializer{Format: format}
}
