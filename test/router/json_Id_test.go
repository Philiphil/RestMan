package router_test

import (
	"testing"

	. "github.com/philiphil/restman/router"
)

func TestMax(t *testing.T) {
	if Max(1, 2, 3, 4, 5) != 5 {
		t.Error("!")
	}
}
