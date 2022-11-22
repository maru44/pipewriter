package pipewriter

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type (
// 	writer[T, P any] interface {
// 		ListWithPagination(ctx context.Context, pagination *P) ([]*T, *P, bool, error)
// 		OverWriteFileName() func(ctx context.Context, origin string) string
// 		Upload(ctx context.Context, dir, name string, file io.Reader) error
// 	}

// 	PipeWriter[T, P any] interface {
// 		writer[T, P]
// 		Data(ctx context.Context, value *T) []byte
// 	}

// 	CsvWriter[T, P any] interface {
// 		writer[T, P]
// 		ValueRow(ctx context.Context, value *T) []string
// 		HeaderRow(ctx context.Context) []string
// 	}
// )

type (
	testWriter struct{}

	chara struct {
		name  string
		age   int
		color string
	}

	pg struct {
		offset int
		limit  int
	}
)

var charas = []*chara{
	{
		name:  "Ruby",
		age:   15,
		color: "red",
	},
	{
		name:  "Weiss",
		age:   17,
		color: "white",
	},
	{
		name:  "Break",
		age:   17,
		color: "black",
	},
	{
		name:  "Yang",
		age:   17,
		color: "yellow",
	},
	{
		name:  "penny",
		color: "light green",
	},
}

func (t *testWriter) ListWithPagination(ctx context.Context, page *pg) ([]*chara, *pg, bool, error) {
	if page == nil {
		return nil, nil, false, errors.New("no page")
	}

	return nil, nil, false, nil
}

func (t *testWriter) OverWriteFileName() func(ctx context.Context, origin string) string {
	return nil
}

func (t *testWriter) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	return nil
}

func (t *testWriter) Data(ctx context.Context, value *chara) []byte {
	return []byte("")
}

func TestWrite(t *testing.T) {
	ctx := context.Background()

	w := &testWriter{}

	tests := []struct {
		name         string
		page         *pg
		wantCnt      int
		wantFileName string
		wantErr      error
	}{
		{
			name:    "err in list",
			wantErr: errors.New("no page"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cnt, fileName, err := Write[chara, pg](ctx, "private", "test.csv", w, tt.page)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCnt, cnt)
			assert.Equal(t, tt.wantFileName, fileName)
		})
	}

}
