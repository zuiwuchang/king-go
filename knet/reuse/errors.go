package reuse

import (
	"errors"
)

var ErrorListenerClosed error = errors.New("listener already closed .")
