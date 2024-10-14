package router_test

import (
	"testing"

	"github.com/philiphil/restman/format"
	. "github.com/philiphil/restman/router"
)

func TestSerializerRenderer(t *testing.T) {
	serializer := SerializerRenderer{
		Data:   "test",
		Format: format.JSON,
		Groups: []string{"test"},
	}

	_ = serializer

}
