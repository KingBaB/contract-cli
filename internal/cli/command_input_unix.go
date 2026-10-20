//go:build unix

package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func readCommandInput(ctx context.Context, path string) (body []byte, err error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fd, wasNonblocking, err := openCommandInputDescriptor(path)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, closeCommandInputDescriptor(fd, wasNonblocking)) }()
	body, err = readCommandInputDescriptor(ctx, fd)
	if err != nil {
		return nil, &os.PathError{Op: "read", Path: path, Err: err}
	}
	return body, nil
}

func openCommandInputDescriptor(path string) (int, bool, error) {
	// Nonblocking Open also covers a FIFO before a writer has connected.
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC|unix.O_NONBLOCK, 0)
	if err != nil {
		return -1, false, &os.PathError{Op: "open", Path: path, Err: err}
	}
	flags, err := unix.FcntlInt(uintptr(fd), unix.F_GETFL, 0)
	if err == nil {
		err = unix.SetNonblock(fd, true)
	}
	if err != nil {
		return -1, false, &os.PathError{Op: "set nonblocking", Path: path, Err: errors.Join(err, unix.Close(fd))}
	}
	return fd, flags&unix.O_NONBLOCK != 0, nil
}

func closeCommandInputDescriptor(fd int, wasNonblocking bool) error {
	var restoreErr error
	if !wasNonblocking {
		// Darwin /dev/fd duplicates share status flags with the original stdin.
		restoreErr = unix.SetNonblock(fd, false)
	}
	return errors.Join(restoreErr, unix.Close(fd))
}

func readCommandInputDescriptor(ctx context.Context, fd int) ([]byte, error) {
	var body []byte
	buffer := make([]byte, 32*1024)
	for {
		if err := waitCommandInputReady(ctx, fd); err != nil {
			return nil, err
		}
		count, err := unix.Read(fd, buffer)
		if errors.Is(err, unix.EINTR) || errors.Is(err, unix.EAGAIN) {
			continue
		}
		if err != nil {
			return nil, err
		}
		body = append(body, buffer[:count]...)
		if count == 0 {
			return body, ctx.Err()
		}
	}
}

func waitCommandInputReady(ctx context.Context, fd int) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		ready, err := pollCommandInputDescriptor(fd)
		if errors.Is(err, unix.EINTR) {
			continue
		}
		if err != nil {
			return fmt.Errorf("wait for command input: %w", err)
		}
		if ready {
			return nil
		}
	}
}
