//go:build !unix

package cli

import (
	"context"
	"io"
)

func readCommandInput(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	file, err := openCommandInput(ctx, path)
	if err != nil {
		return nil, err
	}
	closed := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { _ = file.Close(); close(closed) })
	defer func() {
		if !stop() {
			<-closed
		}
		_ = file.Close()
	}()
	body, err := io.ReadAll(file)
	if cancelled := ctx.Err(); cancelled != nil {
		return nil, cancelled
	}
	return body, err
}
