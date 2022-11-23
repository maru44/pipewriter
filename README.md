# PipeWriter

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/maru44/pipewriter/blob/master/LICENSE)
![ActionsCI](https://github.com/maru44/pipewriter/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/maru44/pipewriter)](https://goreportcard.com/report/github.com/maru44/pipewriter)

You can read data and write file asynchronously using this package.

Godoc is [here](https://pkg.go.dev/github.com/maru44/pipewriter).

## Usage

You can load data and upload file asynchronously if you just create something sutisfying `pipewriter.PipeWriter` interface then call `pipewriter.Write`.

```go
package uploadrepository

import (
	"foo/bar/model"
)

type (
	repo struct{}
)

func NewUploadRepo() *repo {
	return &repo{}
}

func (r *repo) ListWithPagination(ctx context.Context, page *model.Pagination) ([]*model.User, *model.Pagination, bool, error) {
	// ...
}

func (r *repo) OverwriteFileName() func(ctx context.Context, origin string) string {
	// ...
}

func (r *repo) Upload(ctx context.Context, dir, name string, file io.Reader) error {
	// ...
}

func (r *repo) Data(ctx context.Context, value *model.User) []byte {
	// ...
}

```

```go
package upload

import (
	"fmt"

	"github.com/maru44/pipewriter"
	"foo/bar/model"
	"foo/bar/uploadrepository"
)

func Uploading(ctx context.Context) error {
	repo := uploadrepository.NewUploadRepo(ctx)
	count, fileName, err := pipewriter.Write[*model.User, *model.Page](ctx, "private", "filename.txt", repo, &model.Page{})
	if err != nil {
	    return err
	}
	fmt.Println("uploaded:", count)
	fmt.Println("file name:", fileName)
	return nil
}

```

If you want to upload csv file, similarly you only have to create something sutisfying `pipewriter.CsvWriter` interface then call `pipewriter.WriteCSV`.

There are some samples in tests or [e2e](https://github.com/maru44/pipewriter/tree/master/e2e).
