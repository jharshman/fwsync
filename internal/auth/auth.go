package auth

import (
	"context"
	"sync"

	"github.com/spf13/cobra"
	"google.golang.org/api/compute/v1"
)

var (
	GoogleCloudAuthorizedClient *compute.Service
	once                        sync.Once
)

// Auth returns the authentication function.
func Auth() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return auth()
	}
}

// auth will authenticate to the Google APIs using default methods.
// https://cloud.google.com/docs/authentication/application-default-credentials#personal
func auth() (authErr error) {
	once.Do(func() {
		var cc *compute.Service
		cc, authErr = compute.NewService(context.Background())
		GoogleCloudAuthorizedClient = cc
	})
	return
}
