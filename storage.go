package golden

import (
	"io/ioutil"
	"os"
)

// storage implements fsHandler.
type storage struct{}

func (s *storage) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (s *storage) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (s *storage) WriteFile(
	filename string,
	data []byte,
	perm os.FileMode,
) error {
	return ioutil.WriteFile(filename, data, perm)
}
