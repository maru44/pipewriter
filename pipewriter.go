package pipewriter

import (
	"context"
	"io"
)

type (
	writer[T, P any] interface {
		ListWithPagination(ctx context.Context, pagination *P) ([]*T, *P, bool, error)
		OverWriteFileName() func(ctx context.Context, origin string) string
		Upload(ctx context.Context, dir, name string, file io.Reader) error
	}

	PipeWriter[T, P any] interface {
		writer[T, P]
		Data(ctx context.Context, value *T) []byte
	}

	CsvWriter[T, P any] interface {
		writer[T, P]
		ValueRow(ctx context.Context, value *T) []string
		HeaderRow(ctx context.Context) []string
	}
)
