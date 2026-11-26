// Package logging builds the application logger.
package logging

import (
	"log/slog"
	"os"
)

// New returns a JSON logger that writes to stdout.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
