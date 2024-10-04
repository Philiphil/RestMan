package errors_test

import (
	"testing"

	. "github.com/philiphil/apiman/errors"
)

func TestError(t *testing.T) {
	err := ErrNotFound
	if err.Error() != "not found" {
		t.Error("Error() should return the message")
	}
}
