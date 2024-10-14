package security_test

import (
	"testing"

	. "github.com/philiphil/restman/security"
)

func TestNoAuthorizationRequired(t *testing.T) {
	if !NoAuthorizationRequired(nil, nil) {
		t.Error("NoAuthorizationRequired() should return true")
	}
}
