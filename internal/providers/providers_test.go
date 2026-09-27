package providers

import (
	"testing"

	"github.com/matryer/is"
)

func TestNames(t *testing.T) {
	is := is.New(t)
	is.Equal(Names(), []string{Google, Linode})
}

func TestNewInvalidProvider(t *testing.T) {
	is := is.New(t)
	p, err := New("nope", Settings{})
	is.True(err != nil)
	is.Equal(p, nil)
}

func TestNewGoogleRequiresProject(t *testing.T) {
	is := is.New(t)
	p, err := New(Google, Settings{})
	is.True(err != nil)
	is.Equal(p, nil) // must be an untyped nil, not a nil *gcp.Client
}
