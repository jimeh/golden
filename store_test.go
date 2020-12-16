package golden

import (
	"errors"
	"os"
	"testing"

	"github.com/jimeh/golden/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStore_Filename(t *testing.T) {
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
			s := NewStore()
			tb := tt.tbFunc()

			got := s.Filename(tb)

			assert.Equal(t, got, tt.want)
			tb.AssertExpectations(t)
		})
	}
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name   string
		want   string
		fsFunc func() *mocks.FSHandler
		tbFunc func() *mocks.TestingTB
	}{
		{
			name: "hello world",
			want: "hello world!!!1 :)",
			fsFunc: func() *mocks.FSHandler {
				f := &mocks.FSHandler{}
				f.On("ReadFile",
					"testdata/MockTestStore_Get/hello_world.golden",
				).Return([]byte("hello world!!!1 :)"), nil)

				return f
			},
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Get/hello_world")

				return tb
			},
		},
		{
			name: "missing file",
			fsFunc: func() *mocks.FSHandler {
				f := &mocks.FSHandler{}
				f.On("ReadFile",
					"testdata/MockTestStore_Get/missing_file.golden",
				).Return(
					nil,
					errors.New(
						"open testdata/MockTestStore_Get/missing_file.golden: "+
							"no such file or directory",
					),
				)

				return f
			},
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
		},
		{
			name:   "empty name",
			fsFunc: func() *mocks.FSHandler { return &mocks.FSHandler{} },
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("")
				tb.On("Fatalf",
					"could not determine golden file path for: %+v", tb,
				).Return()

				return tb
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fsFunc()
			tb := tt.tbFunc()

			s := NewStore()
			s.fs = fs

			got := s.Get(tb)

			assert.Equal(t, tt.want, string(got))
			fs.AssertExpectations(t)
			tb.AssertExpectations(t)
		})
	}
}

func TestStore_Set(t *testing.T) {
	tests := []struct {
		name    string
		content string
		fsFunc  func() *mocks.FSHandler
		tbFunc  func() *mocks.TestingTB
	}{
		{
			name:    "hello world",
			content: "hello world!!!1 :)",
			fsFunc: func() *mocks.FSHandler {
				f := &mocks.FSHandler{}
				f.On("MkdirAll",
					"testdata/MockTestStore_Set", os.FileMode(0o755),
				).Return(nil)
				f.On("WriteFile",
					"testdata/MockTestStore_Set/hello_world.golden",
					[]byte("hello world!!!1 :)"), os.FileMode(0o644),
				).Return(nil)

				return f
			},
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Set/hello_world")
				tb.On(
					"Logf",
					"updating .golden file: %s",
					"testdata/MockTestStore_Set/hello_world.golden",
				)

				return tb
			},
		},
		{
			name:   "empty name",
			fsFunc: func() *mocks.FSHandler { return &mocks.FSHandler{} },
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("")
				tb.On("Fatalf",
					"could not determine golden file path for: %+v", tb,
				).Return()

				return tb
			},
		},
		{
			name: "mkdir failure",
			fsFunc: func() *mocks.FSHandler {
				f := &mocks.FSHandler{}
				f.On("MkdirAll",
					"testdata/MockTestStore_Set", os.FileMode(0o755),
				).Return(os.ErrPermission)

				return f
			},
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Set/mkdir_failure")
				tb.On("Logf", mock.Anything, mock.Anything)
				tb.On("Fatalf",
					"failed to create .golden directory: %s",
					os.ErrPermission.Error(),
				).Return()

				return tb
			},
		},
		{
			name:    "write failure",
			content: "hi",
			fsFunc: func() *mocks.FSHandler {
				f := &mocks.FSHandler{}
				f.On("MkdirAll",
					"testdata/MockTestStore_Set", os.FileMode(0o755),
				).Return(nil)
				f.On("WriteFile",
					"testdata/MockTestStore_Set/write_failure.golden",
					[]byte("hi"), os.FileMode(0o644),
				).Return(os.ErrDeadlineExceeded)

				return f
			},
			tbFunc: func() *mocks.TestingTB {
				tb := &mocks.TestingTB{}
				tb.On("Name").Return("MockTestStore_Set/write_failure")
				tb.On("Logf", mock.Anything, mock.Anything)
				tb.On("Fatalf",
					"failed to update .golden file: %s",
					os.ErrDeadlineExceeded.Error(),
				).Return()

				return tb
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fsFunc()
			tb := tt.tbFunc()

			s := NewStore()
			s.fs = fs

			s.Set(tb, []byte(tt.content))

			fs.AssertExpectations(t)
			tb.AssertExpectations(t)
		})
	}
}
