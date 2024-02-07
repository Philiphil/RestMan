package apiman

import (
	"fmt"
	"github.com/philiphil/apiman/router"
	"testing"
)

func TestParseAcceptHeader(t *testing.T) {
	fmt.Println(
		router.ParseAcceptHeader(
			"text/plain; q=0.5, text/html,text/x-dvi; q=0.8, text/x-c",
		),
	)
}
