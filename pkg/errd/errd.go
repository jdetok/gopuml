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

type FileCreateError struct {
	ErrMeta
	FName string
	Err   error
}

func (e *FileCreateError) Error() string {
	return fmt.Sprintf("failed to create file %s\n%v", e.FName, e.Err)
}

type FileOpenError struct {
	ErrMeta
	FName string
	Err   error
}

func (e *FileOpenError) Error() string {
	return fmt.Sprintf("failed to open file %s\n%v", e.FName, e.Err)
}

type FileRecursionError struct {
	ErrMeta
	Path string
	Ftyp string
	Err  error
}

func (e *FileRecursionError) Error() string {
	return fmt.Sprintf("recursive func failed to map %s files at %s\n%v",
		e.Ftyp, e.Path, e.Err)
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
