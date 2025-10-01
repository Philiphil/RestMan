package entity_test

import (
	"testing"

	. "github.com/philiphil/restman/orm/entity"
)

func TestEntity(t *testing.T) {
	e := BaseEntity{}
	if e.GetId() != 0 {
		t.Error("GetId() should return the id")
	}
	_, ok := e.SetId(1).(BaseEntity)
	if !ok {
		t.Error("SetId() should return the entity")
	}
}
