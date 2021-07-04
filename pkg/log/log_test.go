package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevLoggernitiaIlisation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name        string
		In          string
		ExpectedErr bool
	}

	testcases := []testCase{
		{Name: "Set logger in dev mode", In: "DEV", ExpectedErr: false},
		{Name: "Set logger in normal mode", In: "", ExpectedErr: false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			_, err := NewLogger(tc.In)
			assert.Equal(t, tc.ExpectedErr, err != nil)
		})
	}
}
