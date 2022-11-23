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
	repo struct{}

	hunter struct {
		name  string
		age   int
		color string
	}

	pg struct {
		offset int
		limit  int
	}
)

var rwby = []*hunter{
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

func (t *repo) ListWithPagination(ctx context.Context, page *pg) ([]*hunter, *pg, bool, error) {
	if page == nil {
		return nil, nil, false, errors.New("no page")
	}

	np := &pg{
		limit:  page.limit,
		offset: page.offset,
	}
	if page.offset > len(rwby) {
		np.offset = 0
		return nil, np, false, nil
	}

	end := page.offset + page.limit
	next := true
	if end > len(rwby) {
		end = len(rwby)
		next = false
		np.offset = 0
	} else {
		np.offset += page.limit
	}

	return rwby[page.offset:end], np, next, nil
}

func (t *repo) OverwriteFileName() func(ctx context.Context, origin string) string {
	return func(ctx context.Context, origin string) string {
		return "rwby" + origin
	}
}

func (t *repo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	f, err := os.OpenFile(dir+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (t *repo) Data(ctx context.Context, value *hunter) []byte {
	return []byte(value.name)
}

func (t *repo) HeaderRow(ctx context.Context) []string {
	return []string{"name", "age", "color"}
}

func (t *repo) ValueRow(ctx context.Context, value *hunter) []string {
	return []string{value.name, strconv.Itoa(value.age), value.color}
}

func main() {
	ctx := context.Background()
	w := &repo{}

	_, _, err := pipewriter.Write[*hunter, *pg](ctx, "./", "test", w, &pg{limit: 1})
	if err != nil {
		panic(err)
	}

	_, _, err = pipewriter.WriteCSV[*hunter, *pg](ctx, "./", "test.csv", w, &pg{limit: 3})
	if err != nil {
		panic(err)
	}
}
