package entity_test

import (
	"testing"

	"github.com/philiphil/restman/orm/entity"
)

// TestCastIdInvalidString ensures invalid numeric string returns NullId (current behavior).
func TestCastIdInvalidString(t *testing.T) {
	id := entity.CastId("abc")
	if id != entity.NullId {
		t.Fatalf("expected NullId for invalid numeric string, got %v", id)
	}
}
