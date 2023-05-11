package auth

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var (
	GoogleCloudAuthorizedClient *compute.Service
	once                        sync.Once
)

func Auth() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		serviceAccountJSONPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if serviceAccountJSONPath == "" {
			return errors.New("missing environment variable GOOGLE_APPLICATION_CREDENTIALS")
		}
		return auth(serviceAccountJSONPath)
	}
}

func auth(svcAcctFile string) (authErr error) {
	once.Do(func() {
		var cc *compute.Service
		cc, authErr = compute.NewService(context.Background(), option.WithCredentialsFile(svcAcctFile))
		GoogleCloudAuthorizedClient = cc
	})
	return
}
