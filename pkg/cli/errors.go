package cli

import "errors"

// ErrNotInitialized is returned when cobra.Command is not initialized
var ErrNotInitialized = errors.New("not initialized")
