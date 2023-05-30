package bufio

import (
	"bytes"
	"errors"
	"io"
)

const (
	defaultBufSize           = 4096
	minReadBufferSize        = 16
	maxConsecutiveEmptyReads = 100
)

var (
	ErrBufferFull        = errors.New("bufio: buffer full")
	ErrNegativeCount     = errors.New("bufio: negative count")
	ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
	ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
	errFillFullBuffer    = errors.New("bufio: try to fill full buffer")
	errNegativeRead      = errors.New("bufio: reader returned negative count from Read")
)

type Reader struct {
	buf  []byte
	rd   io.Reader
	r, w int
	err  error
}

func NewReaderSize(rd io.Reader, size int) *Reader {
	b, ok := rd.(*Reader)
	if ok && len(b.buf) >= size {
		return b
	}

	if size < minReadBufferSize {
		size = minReadBufferSize
	}

	r := new(Reader)
	r.reset(make([]byte, size), rd)
	return r
}

func NewReader(rd io.Reader) *Reader {
	return NewReaderSize(rd, defaultBufSize)
}

func (b *Reader) Reset(r io.Reader) {
	b.reset(b.buf, r)
}

func (b *Reader) reset(buf []byte, r io.Reader) {
	*b = Reader{
		buf: buf,
		rd:  r,
	}
}

func (b *Reader) ResetBuffer(r io.Reader, buf []byte) {
	b.reset(buf, r)
}

func (b *Reader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

func (b *Reader) fill() {
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		panic(errFillFullBuffer)
	}
	for i := 0; i < maxConsecutiveEmptyReads; i++ {
		n, err := b.rd.Read(b.buf[b.w:])
		if n < 0 {
			panic(errNegativeRead)
		}
		b.w += n
		if err != nil {
			b.err = err
			return
		}
		if n > 0 {
			return
		}
	}
	b.err = io.ErrNoProgress
}

func (b *Reader) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrNegativeCount
	}
	if n > len(b.buf) {
		return nil, ErrBufferFull
	}

	for b.w-b.r < n && b.err == nil {
		b.fill()
	}

	var err error
	if avail := b.w - b.r; avail < n {
		n = avail
		if err = b.readErr(); err == nil {
			err = ErrBufferFull
		}
	}
	return b.buf[b.r : b.r+n], err
}

func (b *Reader) Pop(n int) ([]byte, error) {
	d, err := b.Peek(n)
	if err == nil {
		b.r += n
		return d, err
	}
	return nil, err
}

func (b *Reader) Buffered() int { return b.w - b.r }

func (b *Reader) Discard(n int) (discarded int, err error) {
	if n < 0 {
		return 0, ErrNegativeCount
	}
	if n == 0 {
		return
	}

	remain := n

	for {
		skip := b.Buffered()
		if skip == 0 {
			b.fill()
			skip = b.Buffered()
		}

		if skip > remain {
			skip = remain
		}

		b.r += skip
		remain -= skip
		if remain == 0 {
			return n, nil
		}
		if b.err != nil {
			return n - remain, b.readErr()
		}
	}
}

func (b *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return 0, b.readErr()
	}
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			n, b.err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			return n, b.readErr()
		}
		b.fill()
		if b.r == b.w {
			return 0, b.readErr()
		}
	}
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *Reader) ReadByte() (c byte, err error) {
	for b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		b.fill()
	}
	c = b.buf[b.r]
	b.r++
	return c, nil
}

// 寻找某个 delim的前的数据 如果缓冲区满了 则不再寻找
func (b *Reader) ReadSlice(delim byte) (line []byte, err error) {
	for {
		if i := bytes.IndexByte(b.buf[b.r:b.w], delim); i >= 0 {
			line = b.buf[b.r : b.r+i+1]
			b.r += i + 1
			break
		}

		if b.err != nil {
			line = b.buf[b.r:b.w]
			b.r = b.w
			err = b.readErr()
			break
		}

		if b.Buffered() >= len(b.buf) {
			b.r = b.w
			line = b.buf
			err = ErrBufferFull
			break
		}

		b.fill()
	}
	return
}

func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error) {
	line, err = b.ReadSlice('\n')
	if err == ErrBufferFull {
		if len(line) > 0 && line[len(line)-1] == '\r' {
			if b.r == 0 {
				panic("bufio: tried to rewind past start of buffer")
			}
			b.r--
			line = line[:len(line)-1]
		}
		return line, true, nil
	}

	if len(line) == 0 {
		if err != nil {
			line = nil
		}
		return
	}
	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}
	return
}

type Writer struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}

func NewWriterSize(w io.Writer, size int) *Writer {
	b, ok := w.(*Writer)
	if ok && len(b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}
	return &Writer{
		buf: make([]byte, size),
		wr:  w,
	}
}

func NewWriter(w io.Writer) *Writer {
	return NewWriterSize(w, defaultBufSize)
}

func (b *Writer) Reset(w io.Writer) {
	b.err = nil
	b.n = 0
	b.wr = w
}

func (b *Writer) Available() int { return len(b.buf) - b.n }

func (b *Writer) Buffered() int { return b.n }

func (b *Writer) ResetBuffer(w io.Writer, buf []byte) {
	b.buf = buf
	b.err = nil
	b.n = 0
	b.wr = w
}

func (b *Writer) Flush() error {
	err := b.flush()
	return err
}

func (b *Writer) flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortBuffer
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

func (b *Writer) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

func (b *Writer) WriteRaw(p []byte) (nn int, err error) {
	if b.err != nil {
		return 0, b.err
	}
	if b.Buffered() == 0 {
		// if no buffer data, write raw writer
		nn, err = b.wr.Write(p)
		b.err = err
	} else {
		nn, err = b.Write(p)
	}
	return
}

func (b *Writer) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrNegativeCount
	}
	if n > len(b.buf) {
		return nil, ErrBufferFull
	}
	for b.Available() < n && b.err == nil {
		b.flush()
	}
	if b.err != nil {
		return nil, b.err
	}
	d := b.buf[b.n : b.n+n]
	b.n += n
	return d, nil
}

func (b *Writer) WriteString(s string) (int, error) {
	nn := 0
	for len(s) > b.Available() && b.err == nil {
		n := copy(b.buf[b.n:], s)
		b.n += n
		nn += n
		s = s[n:]
		b.flush()
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], s)
	b.n += n
	nn += n
	return nn, nil
}
