package golden

import "os"

//go:generate mockery --name TestingTB
type TestingTB interface {
	Fatalf(format string, a ...interface{})
	Logf(format string, a ...interface{})
	Name() string
}

//go:generate mockery --name fsHandler --structname FSHandler
type fsHandler interface {
	MkdirAll(path string, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}
