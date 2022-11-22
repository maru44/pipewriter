package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/maru44/pipewriter"
)

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

func (t *testWriter) OverWriteFileName() func(ctx context.Context, origin string) string {
	return nil
}

func (t *testWriter) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	f, err := os.OpenFile(dir+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
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

func main() {
	ctx := context.Background()
	w := &testWriter{}

	_, _, err := pipewriter.Write[chara, pg](ctx, "./", "test", w, &pg{limit: 1})
	if err != nil {
		panic(err)
	}
}
