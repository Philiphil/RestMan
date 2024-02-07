package serializer

import (
	"github.com/philiphil/apiman/format"
)

// Serializer is the main serializer struct
type Serializer struct {
	Format format.Format
}

// NewSerializer creates a new instance of Serializer
func NewSerializer(format format.Format) *Serializer {
	return &Serializer{Format: format}
}
