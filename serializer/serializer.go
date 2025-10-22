// the serializer package is responsible for serializing and deserializing objects
// its main purpose is to provide a way to exclude fields from serialization or deserialization
// allowing restman to generate DTO using only tags on the struct
package serializer

import (
	"github.com/philiphil/restman/format"
)

// Serializer is responsible for serializing and deserializing objects
type Serializer struct {
	Format format.Format
}

// NewSerializer creates a new Serializer instance with the specified format.
func NewSerializer(format format.Format) *Serializer {
	return &Serializer{Format: format}
}
