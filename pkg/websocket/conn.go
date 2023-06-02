package websocket

import (
	"io"

	"github.com/xyhubl/yim/pkg/bufio"
)

type Conn struct {
	rwc     io.ReadWriteCloser
	r       *bufio.Reader
	w       *bufio.Writer
	maskKey []byte
}

func newConn(rwc io.ReadWriteCloser, r *bufio.Reader, w *bufio.Writer) *Conn {
	return &Conn{rwc: rwc, r: r, w: w, maskKey: make([]byte, 4)}
}

// 解析websocket数据帧
