package errd

import (
	"fmt"
	"runtime"
)

// error types to use throughout repo

type ErrMeta struct {
	Caller string
}

func (m *ErrMeta) GetCaller() {
	pc, _, _, _ := runtime.Caller(1)
	m.Caller = runtime.FuncForPC(pc).Name()
}

type CreateFileError struct {
	ErrMeta
	FName string
	Err   error
}

func (e *CreateFileError) Error() string {
	return fmt.Sprintf("failed to create file %s\n%v", e.FName, e.Err)
}

type JSONEncodeError struct {
	ErrMeta
	FName string
	Err   error
}

func (e *JSONEncodeError) Error() string {
	return fmt.Sprintf("failed to encode json to file %s\n%v", e.FName, e.Err)
}

type WriterError struct {
	ErrMeta
	WriterLoc string
	NumBytes  int
	Err       error
}

func (e *WriterError) Error() string {
	return fmt.Sprintf("failed to write %d bytes to writer at %s\n%v", e.NumBytes, e.WriterLoc, e.Err)
}
