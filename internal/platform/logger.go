package platform

import (
	"log/slog"
	"os"
)

// InitLogger configures the global [slog] instance based on the application environment.
//
// Supported environments:
//   - "development": Uses a TextHandler with Debug level for human readability.
//   - "production" (or any other): Uses a JSONHandler with Info level for machine-readable logs.
func InitLogger(env string) {
	var handler slog.Handler

	switch env {
	case "development":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	slog.SetDefault(slog.New(handler))
}
