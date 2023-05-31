package bytes

const (
	minThreshold = 1 << 12
)

type Writer struct {
	n   int
	buf []byte
}

func NewWriterSize(n int) *Writer {
	return &Writer{buf: make([]byte, n)}
}

func (w *Writer) Len() int {
	return w.n
}

func (w *Writer) Size() int {
	return len(w.buf)
}

func (w *Writer) Reset() {
	w.n = 0
}

func (w *Writer) Buffer() []byte {
	return w.buf[:w.n]
}

func (w *Writer) Peek(n int) []byte {
	var buf []byte
	w.grow(n)
	buf = w.buf[w.n : w.n+n]
	w.n += n
	return buf
}

func (w *Writer) Write(p []byte) {
	w.grow(len(p))
	w.n += copy(w.buf[w.n:], p)
}

func (w *Writer) grow(n int) {
	if w.n+n < len(w.buf) {
		return
	}
	var (
		buf    []byte
		oldCap = len(w.buf)
	)
	// 当需要的容量大于两倍扩容的容量,则直接按照新的切片需要的容量 + 原切片的容量
	if n > 2*oldCap {
		buf = make([]byte, oldCap+n)
	} else if len(w.buf) < minThreshold { // 到这里说明了 n < 2 * len(w.buf) 如果 len(w.buf) < minThreshold 那么扩容两倍即可
		buf = make([]byte, 2*oldCap)
	} else { // 到这里说明len(w.buf) > minThreshold 那么开始平滑扩容
		newCap := oldCap
		for 0 < newCap && newCap < oldCap+n {
			newCap += (newCap + 3*minThreshold) / 4
		}
		//zh: 如果精度超出了, 那么新的切片需要的容量 + 原切片的容量
		if newCap <= 0 {
			newCap = oldCap + n
		}
		buf = make([]byte, newCap)
	}
	copy(buf, w.buf[:w.n])
	w.buf = buf
}
