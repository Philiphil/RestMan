package route_test

import (
	"testing"

	. "github.com/philiphil/restman/route"
)

func TestDefaultApiRoutes(t *testing.T) {
	test := DefaultApiRoutes()
	if len(test) != 8 {
		t.Error("DefaultApiRoutes should return 8 methods")
	}
}

func TestBatchOperations(t *testing.T) {
	test := BatchOperations()
	if len(test) != 5 {
		t.Error("BatchOperations should return 5 methods")
	}
}

func TestAllApiRoutes(t *testing.T) {
	test := AllApiRoutes()
	if len(test) != 13 {
		t.Error("AllApiRoutes should return 13 methods")
	}
}
