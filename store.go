package golden

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dir = "testdata"
	ext = ".golden"
)

type Store struct {
	fs         fsHandler
	UpdateFunc func() bool
}

func NewStore() *Store {
	return &Store{
		fs:         &storage{},
		UpdateFunc: EnvVarUpdateFunc,
	}
}

func (s *Store) Update() bool {
	return (s.UpdateFunc)()
}

func (s *Store) Filename(tb TestingTB) string {
	if tb == nil || tb.Name() == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s%s", dir, filepath.FromSlash(tb.Name()), ext)
}

func (s *Store) Get(tb TestingTB) []byte {
	gp := s.Filename(tb)
	if gp == "" {
		tb.Fatalf("could not determine golden file path for: %+v", tb)
		return nil
	}

	g, err := s.fs.ReadFile(gp)
	if err != nil {
		tb.Fatalf("failed reading .golden file: %s", err.Error())
	}

	return g
}

func (s *Store) Set(tb TestingTB, input []byte) {
	gp := s.Filename(tb)
	if gp == "" {
		tb.Fatalf("could not determine golden file path for: %+v", tb)
		return
	}
	dir := filepath.Dir(gp)

	tb.Logf("updating .golden file: %s", gp)

	err := s.fs.MkdirAll(dir, os.FileMode(0o755))
	if err != nil {
		tb.Fatalf("failed to create .golden directory: %s", err.Error())
		return
	}

	err = s.fs.WriteFile(gp, input, os.FileMode(0o644)) //nolint:gosec
	if err != nil {
		tb.Fatalf("failed to update .golden file: %s", err.Error())
		return
	}
}
