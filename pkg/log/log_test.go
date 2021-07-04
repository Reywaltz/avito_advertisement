package log

import (
	"os"
	"testing"
)

func TestDevLoggernitiaIlisation(t *testing.T) {
	t.Parallel()

	os.Setenv("DEV", "true")
	defer os.Unsetenv("DEV")

	_, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoggernitiaIlisation(t *testing.T) {
	t.Parallel()

	_, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
}
