# PipeWriter

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/maru44/pipewriter/blob/master/LICENSE)
![ActionsCI](https://github.com/maru44/pipewriter/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/maru44/pipewriter)](https://goreportcard.com/report/github.com/maru44/pipewriter)

PipeWriter connect upload method and list method with io.pipe.
You can execute asynchronously closed method to upload data including io.Writer in it and closed method to list.

Godoc is [here](https://pkg.go.dev/github.com/maru44/pipewriter).

## Usage

You have to sutisfy `pipewriter.PipeWriter` interface then call `pipewriter.Write`.

**sample**

```go
package repository

import (
	"context"
	"io"

	"foo/bar/model"
)

type Bucket interface {
	Upload(ctx context.Context, dir, name string, file io.Reader) error
}

type User interface {
	ListWithPagination(ctx context.Context, page *model.Pagination) ([]*model.User, *model.Pagination, bool, error)
}

```

```go
package persistence

import (
	"context"
	"io"

	"foo/bar/model"
	"foo/bar/repository"
)

type (
	userRepo   struct{}
	bucketRepo struct{}
)

func NewUserRepo() repository.User {
	return &userRepo{}
}

func NewBucketRepo() repository.Bucket {
	return &bucketRepo{}
}

func (r *userRepo) ListWithPagination(ctx context.Context, page *model.Pagination) ([]*model.User, *model.Pagination, bool, error) {
	// ...
}

func (r *bucketRepo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	// ...
}

```

```go
package upload

import (
	"context"
	"io"

	"foo/bar/model"
	"foo/bar/repository"
)

type UploadRepo struct {
	User   repository.User
	Bucket repository.Bucket
}

func (r *UploadRepo) ListWithPagination(ctx context.Context, page *model.Pagination) ([]*model.User, *model.Pagination, bool, error) {
	return r.User.ListWithPagination(ctx, page)
}

func (r *UploadRepo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	return r.Bucket.Upload(ctx, dir, name, file)
}

func (r *UploadRepo) OverwriteFileName() func(ctx context.Context, origin string) string {
	return nil
}

func (r *UploadRepo) Data(ctx context.Context, value *model.User) []byte {
	return []byte(value.String)
}

```

```go
package main

import (
	"context"
	"fmt"

	"foo/bar/model"
	"foo/bar/persistence"
	"foo/bar/upload"
	"github.com/maru44/pipewriter"
)

func main() {
	repo := upload.UploadRepo{
		User: persistence.NewUserRepo(),
		Upload: persistence.NewBucketRepo(),
	}

	count, fileName, err := pipewriter.Write[*model.User, *model.Page](context.Background(), "private", "filename.txt", repo, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("uploaded:", count)
	fmt.Println("file name:", fileName)
}

```

If you want to upload csv file, similarly you only have to create something sutisfying `pipewriter.CsvWriter` interface then call `pipewriter.WriteCSV`.

There are some samples in test files or [e2e](https://github.com/maru44/pipewriter/tree/master/e2e).
