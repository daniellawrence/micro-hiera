package lib

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	VIOLATION_MISSING_INPUT_FILE       = "MISSING_INPUT_FILE"
	VIOLATION_INVALID_INPUT_FILE       = "INVALID_INPUT_FILE"
	VIOLATION_DUPLICATE_OVERRIDE_VALUE = "DUPLICATE_OVERRIDE_VALUE"
	VIOLATION_NON_MAP_MERGE            = "NON_MAP_MERGE"
)

var (
	Violations = map[string]log.Level{
		VIOLATION_MISSING_INPUT_FILE:       log.DebugLevel,
		VIOLATION_INVALID_INPUT_FILE:       log.DebugLevel,
		VIOLATION_DUPLICATE_OVERRIDE_VALUE: log.WarnLevel,
		VIOLATION_NON_MAP_MERGE:            log.ErrorLevel,
	}
)

type WrappedError struct {
	Violation string
	Level     log.Level
	Err       error
}

func (w *WrappedError) Error() string {
	return fmt.Sprintf("%s: %v", w.Violation, w.Err)
}

func NewWrappedError(violationName string, err error) *WrappedError {

	return &WrappedError{
		Violation: violationName,
		Level:     Violations[violationName],
		Err:       err,
	}
}
