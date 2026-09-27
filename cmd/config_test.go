package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func TestToCIDRs(t *testing.T) {
	is := is.New(t)
	in := []string{"1.1.1.1", "2.2.2.0/24", "3.3.3.3/32"}
	is.Equal(toCIDRs(in), []string{"1.1.1.1/32", "2.2.2.0/24", "3.3.3.3/32"})
	is.Equal(in, []string{"1.1.1.1", "2.2.2.0/24", "3.3.3.3/32"}) // input is not modified
}
