package pipewriter

import (
	"context"
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCtx() context.Context {
	return context.WithValue(context.Background(), pipeWriterTestKey{}, true)
}

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
		name:  "Blake",
		age:   17,
		color: "black",
	},
	{
		name:  "Yang",
		age:   17,
		color: "yellow",
	},
	{
		name:  "Penny",
		color: "light green",
	},
}

func (t *testWriter) ListWithPagination(ctx context.Context, page *pg) ([]*chara, *pg, bool, error) {
	if page == nil {
		return nil, nil, false, errors.New("no page")
	}

	np := &pg{
		limit:  page.limit,
		offset: page.offset,
	}
	if page.offset > len(charas) {
		np.offset = 0
		return nil, np, false, nil
	}

	end := page.offset + page.limit
	next := true
	if end > len(charas) {
		end = len(charas)
		next = false
		np.offset = 0
	} else {
		np.offset += page.limit
	}

	return charas[page.offset:end], np, next, nil
}

func (t *testWriter) OverwriteFileName() func(ctx context.Context, origin string) string {
	return nil
}

func (t *testWriter) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	return nil
}

func (t *testWriter) Data(ctx context.Context, value *chara) []byte {
	return []byte(value.name)
}

func (t *testWriter) HeaderRow(ctx context.Context) []string {
	return []string{"name", "age", "color"}
}

func (t *testWriter) ValueRow(ctx context.Context, value *chara) []string {
	return []string{value.name, strconv.Itoa(value.age), value.color}
}

func TestWrite(t *testing.T) {
	ctx := testCtx()
	w := &testWriter{}

	tests := []struct {
		name         string
		page         *pg
		wantCnt      int
		wantFileName string
		wantErr      error
	}{
		{
			name: "OK",
			page: &pg{
				limit: 2,
			},
			wantCnt:      5,
			wantFileName: "test.csv",
		},
		{
			name:    "Err: at ListWithPagination",
			wantErr: errors.New("no page"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cnt, fileName, err := Write[*chara, *pg](ctx, "private", "test.csv", w, tt.page)

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

func TestWriteCSV(t *testing.T) {
	ctx := testCtx()
	w := &testWriter{}

	tests := []struct {
		name         string
		page         *pg
		wantCnt      int
		wantFileName string
		wantErr      error
	}{
		{
			name: "OK",
			page: &pg{
				limit: 2,
			},
			wantCnt:      5,
			wantFileName: "test.csv",
		},
		{
			name:    "Err: at ListWithPagination",
			wantErr: errors.New("no page"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cnt, fileName, err := WriteCSV[*chara, *pg](ctx, "private", "test.csv", w, tt.page)

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
