package app_test

import (
	"testing"
	"time"

	"github.com/mkyc/changes/internal/app"
)

func TestNewFakeDeps_IsolatedBuffers(t *testing.T) {
	fixed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	deps1, stdout1, stderr1 := app.NewFakeDeps(fixed)
	deps2, stdout2, stderr2 := app.NewFakeDeps(fixed)

	stdout1.WriteString("hello")
	stderr1.WriteString("oops")

	if stdout2.Len() != 0 {
		t.Fatalf("expected deps2 stdout to be empty, got %q", stdout2.String())
	}
	if stderr2.Len() != 0 {
		t.Fatalf("expected deps2 stderr to be empty, got %q", stderr2.String())
	}
	if deps1.Stdout == deps2.Stdout {
		t.Fatal("expected deps1 and deps2 to have distinct Stdout writers")
	}
	if deps1.FS == deps2.FS {
		t.Fatal("expected deps1 and deps2 to have distinct filesystems")
	}
}
