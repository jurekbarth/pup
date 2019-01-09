// Package validate provides config validation functions.
package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// RequiredString validation.
func RequiredString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("is required")
	}

	return nil
}

// RequiredStrings validation.
func RequiredStrings(s []string) error {
	for i, v := range s {
		if err := RequiredString(v); err != nil {
			return errors.Wrapf(err, "at index %d", i)
		}
	}
	return nil
}

// name regexp.
var name = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

// Name validation.
func Name(s string) error {
	if !name.MatchString(s) {
		return errors.Errorf("must contain only lowercase alphanumeric characters and '-'")
	}

	return nil
}
