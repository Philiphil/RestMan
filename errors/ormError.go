package errors

import (
	"fmt"
)

// ORM errors is for structuring errors that come from the ORM
// it allows for a more detailed error handling
// hence a better workflow
type OrmError struct {
	Code int
}

func (f OrmError) Error() string {
	return fmt.Sprintf("Error code: %d", f.Code)
}

var (
	ItemNotFound    = OrmError{1}
	NotAllItemFound = OrmError{2}
)
