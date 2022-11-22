package pipewriter

import (
	"context"
	"encoding/csv"
	"io"
	"log"
)

func Write[T, P any](ctx context.Context, dir, name string, w PipeWriter[T, P], page *P) (int, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pr, pw := io.Pipe()
	defer pw.Close()

	var cnt int
	ch := make(chan error)
	go func() {
		hasNext := true
		for hasNext {
			values, pg, next, err := w.ListWithPagination(ctx, page)
			if err != nil {
				ch <- err
				if err := pw.CloseWithError(err); err != nil {
					log.Println(err)
					cancel()
				}
				return
			}

			for _, v := range values {
				if _, err := pw.Write(w.Data(ctx, v)); err != nil {
					ch <- err
					if err := pw.CloseWithError(err); err != nil {
						log.Println(err)
						cancel()
					}
					return
				}
			}

			page = pg
			hasNext = next
			cnt += len(values)
		}

		if err := pw.Close(); err != nil {
			ch <- err
			log.Println(err)
			cancel()
		}
		ch <- nil
	}()

	select {
	case err := <-ch:
		if err != nil {
			return 0, "", err
		}
		break
	}

	if w.OverWriteFileName() != nil {
		name = w.OverWriteFileName()(ctx, name)
	}

	if err := w.Upload(ctx, dir, name, pr); err != nil {
		return 0, "", err
	}

	return cnt, name, nil
}

func WriteCSV[T, P any](ctx context.Context, dir, name string, w CsvWriter[T, P], page *P) (int, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pr, pw := io.Pipe()
	defer pw.Close()

	var cnt int
	ch := make(chan error)
	go func() {
		if _, err := pw.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			ch <- err
			if err := pw.CloseWithError(err); err != nil {
				log.Println(err)
				cancel()
			}
			return
		}

		cw := csv.NewWriter(pw)
		if err := cw.Write(w.HeaderRow(ctx)); err != nil {
			ch <- err
			if err := pw.CloseWithError(err); err != nil {
				log.Println(err)
				cancel()
			}
			return
		}
		cw.Flush()

		hasNext := true
		for hasNext {
			values, pg, next, err := w.ListWithPagination(ctx, page)
			if err != nil {
				ch <- err
				if err := pw.CloseWithError(err); err != nil {
					log.Println(err)
					cancel()
				}
				return
			}

			for _, v := range values {
				if err := cw.Write(w.ValueRow(ctx, v)); err != nil {
					ch <- err
					if err := pw.CloseWithError(err); err != nil {
						log.Println(err)
						cancel()
					}
					return
				}
			}
			cw.Flush()

			page = pg
			hasNext = next
			cnt += len(values)
		}

		if err := pw.Close(); err != nil {
			ch <- err
			log.Println(err)
			cancel()
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return 0, "", err
		}
		break
	}

	if w.OverWriteFileName() != nil {
		name = w.OverWriteFileName()(ctx, name)
	}

	if err := w.Upload(ctx, dir, name, pr); err != nil {
		return 0, "", err
	}

	return cnt, name, nil
}
