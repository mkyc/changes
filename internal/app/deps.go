package app

import (
	"io"
	"os"
	"time"

	"github.com/spf13/afero"
)

// Deps holds the dependencies commands need for deterministic behavior.
// Left as a plain, extensible struct so later tasks (e.g. a Git adapter)
// can add fields without changing the DI pattern established here.
type Deps struct {
	Clock  func() time.Time
	FS     afero.Fs
	Stdout io.Writer
	Stderr io.Writer
}

// NewRealDeps builds a Deps wired to real time, the real OS filesystem,
// and the process's standard streams.
func NewRealDeps() *Deps {
	return &Deps{
		Clock:  time.Now,
		FS:     afero.NewOsFs(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
