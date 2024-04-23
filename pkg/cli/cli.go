// Package cli helps with building CLI services
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

var cmd *cobra.Command

// CommandHandlerFunc describes the header of functions that can be attached to a command
// All the functions passed to AddCommand must respect it
type CommandHandlerFunc func(ctx context.Context)

// Init initializes the CLI service
func Init(appID, appDescription string) {
	cmd = &cobra.Command{
		Use:   appID,
		Short: appDescription,
		Long:  appDescription,
	}
}

// AddCommand adds a new command to CLI service
func AddCommand(command, description string, handlerFunc CommandHandlerFunc) error {
	if cmd == nil {
		return ErrNotInitialized
	}

	cmd.AddCommand(&cobra.Command{
		Use:   command,
		Short: description,
		Long:  description,
		Run: func(cmd *cobra.Command, _ []string) {
			handlerFunc(cmd.Context())
		},
	})

	return nil
}

// AssignStringFlag set a string flag to CLI service
func AssignStringFlag(target *string, name, defaultValue, description string) {
	cmd.PersistentFlags().StringVar(target, name, defaultValue, description)
}

// Run runs the CLI service with a context attached
func Run(ctx context.Context) error {
	return cmd.ExecuteContext(ctx)
}
