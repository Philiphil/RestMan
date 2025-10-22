package entity

import (
	"strconv"
)

// ID is a type that represents an entity's primary identifier
type ID uint

// NullId represents a null or unset ID value.
const NullId ID = 0

// String converts the ID to its string representation.
func (e ID) String() string {
	return strconv.FormatUint(uint64(e), 10)
}

// CastId converts various types to an ID, returning NullId if conversion fails.
func CastId(id any) ID {
	switch v := id.(type) {
	case ID:
		return v
	case int:
		return ID(v)
	case uint:
		return ID(v)
	case string:
		convertedID, _ := strconv.ParseUint(v, 10, 64)
		return ID(convertedID)
	}

	return NullId
}
