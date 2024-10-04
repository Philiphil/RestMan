package entity_test

import (
	"testing"

	. "github.com/philiphil/apiman/orm/entity"
)

func TestEntity(t *testing.T) {
	e := Entity{}
	if e.GetId() != 0 {
		t.Error("GetId() should return the id")
	}
}
