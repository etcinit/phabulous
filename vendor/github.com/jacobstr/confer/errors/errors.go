package err

import (
	"fmt"
	"strings"
)

type UnsupportedConfigError string

// Returned when we're given a file we can't handle.
func (str UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(str))
}

type LoadError struct {
	Msg    string
	Errors []error
}

// Returned during multi-file load operations that
// encounter one or more failures.
func (m *LoadError) Error() string {
	merged := []string{}
	for _, err := range m.Errors {
		merged = append(merged, err.Error())
	}
	return m.Msg + " " + strings.Join(merged, ", ")
}
