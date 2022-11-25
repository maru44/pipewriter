package pipewriter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testWriterI struct{}

	minarai interface {
		fmt.Stringer
		spell() string
	}

	witch struct {
		name  string
		color string
		fairy string
	}
)

var witches = []*witch{
	{
		name:  "Doremi",
		color: "red",
		fairy: "Dodo",
	},
	{
		name:  "Haduki",
		color: "orange",
		fairy: "Rere",
	},
	{
		name:  "Aiko",
		color: "blue",
		fairy: "Mimi",
	},
	{
		name:  "Onpu",
		color: "purple",
		fairy: "Roro",
	},
	{
		name:  "Momoko",
		color: "yellow",
		fairy: "Nini",
	},
	{
		name:  "Pop",
		color: "pink",
		fairy: "Fafa",
	},
	{
		name:  "Hana",
		color: "white",
		fairy: "Toto",
	},
	{
		name:  "Nozomi",
		color: "green",
	},
}

func (w *witch) String() string {
	return fmt.Sprintf("%s_%s", w.name, w.color)
}

func (w *witch) spell() string {
	switch w.name {
	case "Doremi":
		return "Pirika Pirilala Poporina Peperuto! "
	case "Haduki":
		return "Paipai Ponpoi Puwapuwa Puu!"
	case "Aiko":
		return "Pameruku Raruku Rarirori Poppun!"
	case "Onpu":
		return "Pururun Purun Famifami Faa!"
	case "Momoko":
		return "Perutan Petton Pararira Pon!"
	case "Pop":
		return "Pipitto Puritto Puritan Peperuto!"
	case "Hana":
		return "Pororin Pyuarin Hanahana Pii!"
	case "Nozomi":
		return "Potolila Potolala Puritan Pyantan"
	default:
		return ""
	}
}

func (t *testWriterI) ListWithPagination(ctx context.Context, page *pg) ([]minarai, *pg, bool, error) {
	out := make([]minarai, len(witches))
	for i, w := range witches {
		out[i] = w
	}
	return out, nil, false, nil
}

func (t *testWriterI) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	if name == "" {
		return errors.New("blank file name")
	}
	return nil
}

func (t *testWriterI) Data(ctx context.Context, value minarai) []byte {
	return []byte(value.spell())
}

func TestWrite_Interface(t *testing.T) {
	ctx := testCtx()
	w := &testWriterI{}

	tests := []struct {
		name         string
		page         *pg
		fileName     string
		wantCnt      int
		wantFileName string
		wantErr      error
	}{
		{
			name:         "OK",
			fileName:     "test.csv",
			wantCnt:      8,
			wantFileName: "test.csv",
		},
		{
			name: "Err: at Upload",
			page: &pg{
				limit: 1,
			},
			wantErr: errors.New("blank file name"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cnt, fileName, err := Write[minarai, *pg](ctx, "private", tt.fileName, w, tt.page)

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
