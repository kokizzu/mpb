package cwriter

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
)

// https://github.com/dylanaraps/pure-sh-bible#cursor-movement
const (
	escOpen  = "\x1b["
	cuuAndEd = "A\x1b[J"
)

// ErrNotTTY not a TeleTYpewriter error.
var ErrNotTTY = errors.New("not a terminal")

// Writer is a buffered the writer that updates the terminal. The
// contents of writer will be flushed when Flush is called.
type Writer struct {
	*bytes.Buffer
	ew       escWriter
	out      io.Writer
	lines    int // used by writer_windows only
	fd       int
	terminal bool
	termSize func(int) (int, int, error)
}

// New returns a new Writer with defaults.
func New(out io.Writer) *Writer {
	w := &Writer{
		Buffer: new(bytes.Buffer),
		ew:     escWriter(make([]byte, 8, 16)),
		out:    out,
		termSize: func(_ int) (int, int, error) {
			return -1, -1, ErrNotTTY
		},
	}
	if f, ok := out.(*os.File); ok {
		w.fd = int(f.Fd())
		if IsTerminal(w.fd) {
			w.terminal = true
			w.termSize = func(fd int) (int, int, error) {
				return GetSize(fd)
			}
		}
	}
	return w
}

// GetTermSize returns WxH of underlying terminal.
func (w *Writer) GetTermSize() (width, height int, err error) {
	return w.termSize(w.fd)
}

type escWriter []byte

func (b escWriter) ansiCuuAndEd(out io.Writer, n int) error {
	b = strconv.AppendInt(b[:copy(b, escOpen)], int64(n), 10)
	_, err := out.Write(append(b, cuuAndEd...))
	return err
}
