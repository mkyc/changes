package app

import (
	"bytes"
	"time"

	"github.com/spf13/afero"
)

// NewFakeDeps builds a Deps suitable for tests: a fixed clock, an
// in-memory filesystem, and buffer-backed streams so tests can assert
// exact output without touching the real filesystem or process-global
// state.
func NewFakeDeps(fixedTime time.Time) (*Deps, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	deps := &Deps{
		Clock:  func() time.Time { return fixedTime },
		FS:     afero.NewMemMapFs(),
		Stdout: stdout,
		Stderr: stderr,
	}
	return deps, stdout, stderr
}
