package ui

import (
	"io"
	"os"
)

func WithSpinner[T any](w io.Writer, msg string, fn func() (T, error)) (T, error) {
	if f, ok := w.(*os.File); !ok || !IsTTY(f) {
		return fn()
	}

	s := NewSpinner(w, msg)
	s.Start()
	defer s.Stop()

	return fn()
}
