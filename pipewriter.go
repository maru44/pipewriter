package pipewriter

import (
	"context"
	"io"
)

type (
	pipeWriterTestKey struct{}

	writer[T, P any] interface {
		// ListWithPagination is method to load data.
		ListWithPagination(ctx context.Context, pagination *P) ([]*T, *P, bool, error)
		// OverwriteFileName returns function to overwrite file name before Upload method.
		// If this return nil, Upload method will be given the `name` argument given at Write or CsvWrite function.
		OverwriteFileName() func(ctx context.Context, origin string) string
		// Upload is to upload data.
		Upload(ctx context.Context, dir, name string, file io.Reader) error
	}

	// PipeWriter is writer interface for Write function.
	PipeWriter[T, P any] interface {
		writer[T, P]
		// Data returns data to write file made by value typed T.
		Data(ctx context.Context, value *T) []byte
	}

	// CsvWriter is writer interface for WriteCSV function.
	CsvWriter[T, P any] interface {
		writer[T, P]
		// ValueRow returns csv row made by value typed T.
		ValueRow(ctx context.Context, value *T) []string
		// HeaderRow returns csv header.
		HeaderRow(ctx context.Context) []string
	}
)

func isTest(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	test, ok := ctx.Value(pipeWriterTestKey{}).(bool)
	if ok {
		return test
	}
	return false
}
