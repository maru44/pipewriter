package main

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/maru44/pipewriter"
)

type (
	repo struct {
		hunter iHunterRepo
	}

	hunterRepo  struct{}
	iHunterRepo interface {
		ListWithPagination(ctx context.Context, page *pg) ([]*hunter, *pg, bool, error)
	}

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

func (r *hunterRepo) ListWithPagination(ctx context.Context, page *pg) ([]*hunter, *pg, bool, error) {
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

func (r *repo) Data(ctx context.Context, value *hunter) []byte {
	return []byte(value.name)
}

func (r *repo) HeaderRow(ctx context.Context) []string {
	return []string{"name", "age", "color"}
}

func (r *repo) ValueRow(ctx context.Context, value *hunter) []string {
	return []string{value.name, strconv.Itoa(value.age), value.color}
}

func (r *repo) OverwriteFileName() func(ctx context.Context, origin string) string {
	return func(ctx context.Context, origin string) string {
		return "rwby" + origin
	}
}

func (r *repo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	f, err := os.OpenFile(dir+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return err
	}
	return nil
}

func (r *repo) ListWithPagination(ctx context.Context, page *pg) ([]*hunter, *pg, bool, error) {
	return r.hunter.ListWithPagination(ctx, page)
}

func main() {
	ctx := context.Background()
	w := &repo{
		hunter: &hunterRepo{},
	}

	_, _, err := pipewriter.Write[*hunter, *pg](ctx, "./", "test", w, &pg{limit: 1})
	if err != nil {
		panic(err)
	}

	_, _, err = pipewriter.WriteCSV[*hunter, *pg](ctx, "./", "test.csv", w, &pg{limit: 3})
	if err != nil {
		panic(err)
	}
}
