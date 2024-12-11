package util

import (
	"errors"
	"log/slog"
)

func LogWarn(commandName string, enableDebug bool, errorMessage string) {
	if !enableDebug {
		return
	}

	slog.Warn(commandName, errorMessage, "")
}

func LogErrorPanic(commandName string, errorMessage string) {
	slog.Error(commandName, errorMessage, "")
	panic(errors.New(errorMessage))
}
