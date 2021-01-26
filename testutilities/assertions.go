package testutilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckError is a helper that asserts a test's error matches if an expected error is passed in,
// otherwise it asserts that a tests's error is nil
func CheckError(t *testing.T, expectedErr, actualErr error) {
	if expectedErr != nil {
		assert.EqualError(t, actualErr, expectedErr.Error())
	} else {
		assert.NoError(t, actualErr)
	}
}
