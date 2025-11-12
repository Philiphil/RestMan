package router_test

import (
	"testing"

	. "github.com/philiphil/restman/router"
)

// TestGenerateETagFromResponseStable ensures same body yields identical ETag
func TestGenerateETagFromResponseStable(t *testing.T) {
	b := []byte("hello world")
	etag1 := GenerateETag(b)
	etag2 := GenerateETag(b)
	if etag1 != etag2 {
		t.Fatalf("expected identical etags, got %s and %s", etag1, etag2)
	}
	if len(etag1) == 0 || etag1[0] != '"' || etag1[len(etag1)-1] != '"' {
		t.Fatalf("etag should be quoted per RFC, got %s", etag1)
	}
}
