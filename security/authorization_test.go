package security_test

import (
	"testing"

	. "github.com/philiphil/apiman/security"
)

func TestNoAuthorizationRequired(t *testing.T) {
	if !NoAuthorizationRequired(nil, nil) {
		t.Error("NoAuthorizationRequired() should return true")
	}
}
