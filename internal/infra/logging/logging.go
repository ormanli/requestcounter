package logging

import (
	"log/slog"

	"github.com/ormanli/requestcounter/internal/app/requestcounter"
)

// Setup setups logger configuration.
func Setup(cfg requestcounter.Config) {
	if cfg.InitDebug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
