package common

import (
	"errors"
	"testing"
)

func TestNormalizeCLIErrorUnknownCommand(t *testing.T) {
	normalized := NormalizeCLIError(errors.New("unknown command \"foo\" for \"zephyr-mail\""))
	if normalized == nil {
		t.Fatal("expected error")
	}
	if normalized.Error() != "Unknown command" {
		t.Fatalf("expected normalized message, got: %s", normalized.Error())
	}
}

func TestNormalizeCLIErrorPreservesExitCode(t *testing.T) {
	err := WrapExitCode(errors.New("Unknown command"), 64)
	normalized := NormalizeCLIError(err)
	if ExitCode(normalized) != 64 {
		t.Fatalf("expected exit code 64, got %d", ExitCode(normalized))
	}
	if normalized.Error() != "Unknown command" {
		t.Fatalf("expected normalized message, got: %s", normalized.Error())
	}
}

func TestExitCodeFallbackIsOne(t *testing.T) {
	if ExitCode(errors.New("boom")) != 1 {
		t.Fatalf("expected exit code 1, got %d", ExitCode(errors.New("boom")))
	}
}
