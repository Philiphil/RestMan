package apiman

import (
	"fmt"
	"testing"
)

func TestParseAcceptHeader(t *testing.T) {
	fmt.Println(
		ParseAcceptHeader(
			"text/plain; q=0.5, text/html,text/x-dvi; q=0.8, text/x-c",
		),
	)
}
