package pipewriter

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"sync"
	"time"
)

func Write[T, P any](ctx context.Context, dir, name string, w PipeWriter[T, P], page *P) (int, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if w.OverWriteFileName() != nil {
		name = w.OverWriteFileName()(ctx, name)
	}

	pr, pw := io.Pipe()
	defer pw.Close()

	var cnt int
	ch := make(chan error)
	testMode := isTest(ctx)
	wg := sync.WaitGroup{}
	go func() {
		defer pw.Close()
		defer wg.Done()
		wg.Add(1)
		for {
			values, pg, next, err := w.ListWithPagination(ctx, page)
			if err != nil {
				ch <- err
			}

			if !testMode {
				for _, v := range values {
					if _, err := pw.Write(w.Data(ctx, v)); err != nil {
						ch <- err
					}
				}
			}

			cnt += len(values)
			page = pg
			if !next {
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		time.Sleep(1 * time.Millisecond)
		close(ch)
	}()

	for err := range ch {
		return 0, "", err
	}

	if err := w.Upload(ctx, dir, name, pr); err != nil {
		return 0, "", err
	}

	return cnt, name, nil
}

func WriteCSV[T, P any](ctx context.Context, dir, name string, w CsvWriter[T, P], page *P) (int, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if w.OverWriteFileName() != nil {
		name = w.OverWriteFileName()(ctx, name)
	}

	pr, pw := io.Pipe()
	defer pw.Close()

	var cnt int
	ch := make(chan error)
	testMode := isTest(ctx)
	wg := sync.WaitGroup{}
	go func() {
		defer pw.Close()
		defer wg.Done()
		wg.Add(1)
		if !testMode {
			if _, err := pw.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
				ch <- err
				if err := pw.CloseWithError(err); err != nil {
					log.Println(err)
					cancel()
				}
				return
			}
		}

		cw := csv.NewWriter(pw)
		if !testMode {
			if err := cw.Write(w.HeaderRow(ctx)); err != nil {
				ch <- err
			}
			cw.Flush()
		}

		for {
			values, pg, next, err := w.ListWithPagination(ctx, page)
			if err != nil {
				ch <- err
			}

			if !testMode {
				for _, v := range values {
					if err := cw.Write(w.ValueRow(ctx, v)); err != nil {
						ch <- err
					}
				}
			}
			cw.Flush()

			page = pg
			cnt += len(values)
			if !next {
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		time.Sleep(1 * time.Millisecond)
		close(ch)
	}()

	for err := range ch {
		return 0, "", err
	}

	if err := w.Upload(ctx, dir, name, pr); err != nil {
		return 0, "", err
	}

	return cnt, name, nil
}
