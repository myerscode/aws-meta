package util

import (
	"github.com/pterm/pterm"
)

var logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

func LogInfo(msg string, args ...[]pterm.LoggerArgument) {
	logger.Info(msg, args...)
}

func LogTrace(msg string, args ...[]pterm.LoggerArgument) {
	logger.Trace(msg, args...)
}

func LogError(msg string, args ...[]pterm.LoggerArgument) {
	logger.Error(msg, args...)
}

func LogWarning(msg string, args ...[]pterm.LoggerArgument) {
	logger.Warn(msg, args...)
}
