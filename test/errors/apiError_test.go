package errors_test

import (
	"fmt"
	"testing"

	. "github.com/philiphil/restman/errors"
)

func TestError(t *testing.T) {
	err := ErrNotFound
	if err.Error() != "not found" {
		t.Error("Error() should return the message")
	}
	err2 := NotAllItemFound
	if err2.Error() != fmt.Sprintf("Error code: %d", err2.Code) {
		t.Error("Error() should return the message")
	}

}
