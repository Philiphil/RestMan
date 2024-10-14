package entity_test

import (
	"testing"

	. "github.com/philiphil/restman/orm/entity"
)

func TestEntity(t *testing.T) {
	e := Entity{}
	if e.GetId() != 0 {
		t.Error("GetId() should return the id")
	}
	e, ok := e.SetId(1).(Entity)
	if !ok {
		t.Error("SetId() should return the entity")
	}
}
