package pipewriter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

type (
	repo struct {
		iPersonRepo
		iUploadRepo
	}

	personRepo struct{}
	uploadRepo struct{}

	iPersonRepo interface {
		ListWithPagination(ctx context.Context, page int) ([]*person, int, bool, error)
		Data(ctx context.Context, value *person) []byte
		OverwriteFileName() func(ctx context.Context, origin string) string
	}

	iUploadRepo interface {
		Upload(ctx context.Context, dir, name string, file io.Reader) error
	}

	person struct {
		Firstname string `json:"first name"`
		Lastname  string `json:"last name"`
		Email     string
		Age       string
		CreatedAt string
	}
)

const baseURL = "https://maru44.github.io/maru44/pipewriter"

func get(url string) ([]*person, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	cli := new(http.Client)
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var out []*person
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (p *person) String() string {
	return p.Firstname + p.Lastname + p.Email
}

func (r *personRepo) ListWithPagination(ctx context.Context, page int) ([]*person, int, bool, error) {
	if page > 8 {
		return nil, 0, false, nil
	}
	url := fmt.Sprintf("%s/bench_%d.json", baseURL, page)
	out, err := get(url)

	if err != nil {
		return nil, 0, false, nil
	}
	np := page + 1
	return out, np, true, nil
}

func (r *personRepo) Data(ctx context.Context, value *person) []byte {
	return []byte(value.Firstname + value.Lastname + value.Email)
}

func (r *personRepo) OverwriteFileName() func(ctx context.Context, origin string) string {
	return nil
}

func (r *uploadRepo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return err
	}
	return nil
}

func BenchmarkWithPipe(b *testing.B) {
	repo := &repo{
		iPersonRepo: &personRepo{},
		iUploadRepo: &uploadRepo{},
	}

	if _, _, err := Write[*person, int](context.Background(), "", "bench_w_pipe", repo, 0); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkWithoutPipe(b *testing.B) {
	repo := &repo{
		iPersonRepo: &personRepo{},
		iUploadRepo: &uploadRepo{},
	}

	ctx := context.Background()
	var page int
	var out string
	for {
		ps, p, next, err := repo.ListWithPagination(ctx, page)
		if err != nil {
			b.Fatal(err)
		}
		if !next {
			break
		}
		for _, p := range ps {
			out += p.String()
		}
		page = p
	}
	if err := repo.Upload(ctx, "", "bench_wo_pipe", bytes.NewReader([]byte(out))); err != nil {
		b.Fatal(err)
	}
}
