package entity

import (
	"strconv"
)

type ID uint

func (e ID) String() string {
	return strconv.FormatUint(uint64(e), 10)
}

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

	return ID(0)
}
