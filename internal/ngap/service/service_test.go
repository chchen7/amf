package service

import (
	"errors"
	"fmt"
	"io"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRetryableSCTPError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "EINTR", err: syscall.EINTR, want: true},
		{name: "wrapped EINTR", err: fmt.Errorf("accept failed: %w", syscall.EINTR), want: true},
		{name: "EAGAIN", err: syscall.EAGAIN, want: true},
		{name: "wrapped EAGAIN", err: fmt.Errorf("accept failed: %w", syscall.EAGAIN), want: true},
		{name: "other error", err: errors.New("accept failed"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isRetryableSCTPError(tt.err))
		})
	}
}

func TestClassifySCTPReadError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want sctpReadError
	}{
		{name: "EOF", err: io.EOF, want: sctpReadEOF},
		{name: "wrapped EOF", err: fmt.Errorf("read failed: %w", io.EOF), want: sctpReadEOF},
		{name: "unexpected EOF", err: io.ErrUnexpectedEOF, want: sctpReadEOF},
		{name: "wrapped unexpected EOF", err: fmt.Errorf("read failed: %w", io.ErrUnexpectedEOF), want: sctpReadEOF},
		{name: "timeout", err: syscall.EAGAIN, want: sctpReadTimeout},
		{name: "wrapped timeout", err: fmt.Errorf("read failed: %w", syscall.EAGAIN), want: sctpReadTimeout},
		{name: "interrupted", err: syscall.EINTR, want: sctpReadInterrupted},
		{name: "wrapped interrupted", err: fmt.Errorf("read failed: %w", syscall.EINTR), want: sctpReadInterrupted},
		{name: "closed", err: syscall.EBADF, want: sctpReadClosed},
		{name: "wrapped closed", err: fmt.Errorf("read failed: %w", syscall.EBADF), want: sctpReadClosed},
		{name: "other error", err: errors.New("read failed"), want: sctpReadFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, classifySCTPReadError(tt.err))
		})
	}
}
