package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jimeh/golden/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		name   string
		tbFunc func() *mocks.TestingTB
		want   string
	}{
		{
			name: "test case",
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Path/test_case")

				return tb
			},
			want: "testdata/MockTestStore_Path/test_case.golden",
		},
		{
			name: "empty name",
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("")

				return tb
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := tt.tbFunc()

			got := Filename(tb)

			assert.Equal(t, got, tt.want)
			tb.AssertExpectations(t)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		tbFunc  func() *mocks.TestingTB
		file    string
		content string
		create  bool
	}{
		{
			name: "existing file",
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Get/existing_file")

				return tb
			},
			file:    "testdata/MockTestStore_Get/existing_file.golden",
			content: "hello world!!!1 :)",
			create:  true,
		},
		{
			name: "missing file",
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Get/missing_file")
				tb.On("Fatalf",
					"failed reading .golden file: %s",
					"open testdata/MockTestStore_Get/missing_file.golden: "+
						"no such file or directory",
				).Return()

				return tb
			},
			file:   "testdata/MockTestStore_Get/missing_file.golden",
			create: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := tt.tbFunc()

			if tt.create {
				err := os.MkdirAll(filepath.Dir(tt.file), 0o755)
				assert.NoError(t, err)

				err = ioutil.WriteFile(tt.file, []byte(tt.content), 0o644)
				assert.NoError(t, err)
			}

			got := Get(tb)

			if tt.create {
				assert.Equal(t, tt.content, string(got))
				err := os.Remove(tt.file)
				assert.NoErrorf(t, err,
					"filed to remove temporary test file: %s", tt.file,
				)
			}
			tb.AssertExpectations(t)
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name    string
		tbFunc  func() *mocks.TestingTB
		file    string
		content string
	}{
		{
			name: "hello world",
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Get/hello_world")
				tb.On(
					"Logf",
					"updating .golden file: %s",
					"testdata/MockTestStore_Get/hello_world.golden",
				)

				return tb
			},
			file:    "testdata/MockTestStore_Get/hello_world.golden",
			content: "hello world!!!1 :)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := tt.tbFunc()

			Set(tb, []byte(tt.content))

			require.FileExists(t, tt.file)
			content, err := ioutil.ReadFile(tt.file)
			require.NoError(t, err)

			assert.Equal(t, tt.content, string(content))
			tb.AssertExpectations(t)

			err = os.Remove(tt.file)
			assert.NoErrorf(t, err,
				"filed to remove temporary test file: %s", tt.file,
			)
		})
	}
}
