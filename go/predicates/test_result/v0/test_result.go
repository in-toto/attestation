/*
Validator APIs for test-result v0.1 protos.
*/
package v0

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrResultRequired        = errors.New("result required")
	ErrInvalidResult         = errors.New("result must be one of PASSED, WARNED or FAILED")
	ErrConfigurationRequired = errors.New("configuration required")
)

// Results are the values the spec allows for result.
var Results = []string{"PASSED", "WARNED", "FAILED"}

// Validate checks the fields the test-result predicate requires: result is
// one of PASSED, WARNED or FAILED, and configuration lists at least one
// valid ResourceDescriptor.
func (t *TestResult) Validate() error {
	result := t.GetResult()
	if result == "" {
		return ErrResultRequired
	}
	if !slices.Contains(Results, result) {
		return fmt.Errorf("%w, got %q", ErrInvalidResult, result)
	}

	configuration := t.GetConfiguration()
	if len(configuration) == 0 {
		return ErrConfigurationRequired
	}
	for i, rd := range configuration {
		if err := rd.Validate(); err != nil {
			return fmt.Errorf("configuration[%d] error: %w", i, err)
		}
	}

	return nil
}
