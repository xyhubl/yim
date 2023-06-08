package websocket

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/xyhubl/yim/pkg/bufio"
)

const (
	finBit  = 1 << 7
	rsv1Bit = 1 << 6
	rsv2Bit = 1 << 5
	rsv3Bit = 1 << 4
	opBit   = 0x0f

	maskBit = 1 << 7
	lenBit  = 0x7f

	continuationFrame        = 0
	continuationFrameMaxRead = 100
)

const (
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
)

var (
	ErrMessageClose   = errors.New("close control message")
	ErrMessageMaxRead = errors.New("continuation frame max read")
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

func (c *Conn) Close() error {
	return c.rwc.Close()
}

// 解析websocket数据帧
func (c *Conn) readFrame() (fin bool, op int, payload []byte, err error) {
	var (
		b          byte
		mask       bool
		p          []byte
		maskKey    []byte
		payloadLen int64
	)

	// 1.First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}

	// is final frame
	fin = (b & finBit) != 0
	// rsv should be 0
	if rsv := b & (rsv1Bit | rsv2Bit | rsv3Bit); rsv != 0 {
		return false, 0, nil, fmt.Errorf("unexpected reserved bits rsv1=%d, rsv2=%d, rsv3=%d", b&rsv1Bit, b&rsv2Bit, b&rsv3Bit)
	}

	// op code
	op = int(b & opBit)

	// 2.Second byte. Mask/Payload len(7bits)
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}

	// is mask payload
	mask = (b & maskBit) != 0

	// payload length
	switch b & lenBit {
	case 126:
		// 16 bits
		if p, err = c.r.Pop(2); err != nil {
			return
		}
		payloadLen = int64(binary.BigEndian.Uint16(p))
	case 127:
		// 64 bits
		if p, err = c.r.Pop(8); err != nil {
			return
		}
		payloadLen = int64(binary.BigEndian.Uint64(p))
	default:
		// 7 bits
		payloadLen = int64(b & lenBit)
	}
	if mask {
		maskKey, err = c.r.Pop(4)
		if err != nil {
			return
		}
		if c.maskKey == nil {
			c.maskKey = make([]byte, 4)
		}
		copy(c.maskKey, maskKey)
	}
	// read payload
	if payloadLen > 0 {
		if payload, err = c.r.Pop(int(payloadLen)); err != nil {
			return
		}
		if mask {
			maskBytes(c.maskKey, 0, payload)
		}
	}
	return
}

func maskBytes(key []byte, pos int, b []byte) int {
	for i := range b {
		b[i] ^= key[pos&3]
		pos++
	}
	return pos & 3
}

// 读取消息
func (c *Conn) ReadMessage() (op int, payload []byte, err error) {
	var (
		fin         bool
		finOp, n    int
		partPayload []byte
	)
	for {
		if fin, op, partPayload, err = c.readFrame(); err != nil {
			return
		}
		switch op {
		case BinaryMessage, TextMessage, continuationFrame:
			if fin && len(payload) == 0 {
				return op, partPayload, nil
			}
			payload = append(payload, partPayload...)
			if op != continuationFrame {
				finOp = op
			}
			if fin {
				op = finOp
				return
			}
		case PingMessage:
		case PongMessage:
			if err = c.WriteMessage(PongMessage, partPayload); err != nil {
				return
			}
		case CloseMessage:
			err = ErrMessageClose
			return
		default:
			err = fmt.Errorf("unknown control message, fin=%t, op=%d", fin, op)
			return
		}
		if n > continuationFrameMaxRead {
			err = ErrMessageMaxRead
			return
		}
		n++
	}
}

func (c *Conn) WriteMessage(msgType int, msg []byte) (err error) {
	if err = c.WriteHeader(msgType, len(msg)); err != nil {
		return
	}
	err = c.WriteBody(msg)
	return
}

func (c *Conn) WriteHeader(msgType int, length int) (err error) {
	var h []byte
	if h, err = c.w.Peek(2); err != nil {
		return
	}
	// 1.First byte. FIN/RSV1/RSV2/RSV3/OpCode(4bits)
	h[0] = 0
	h[0] |= finBit | byte(msgType)
	// 2.Second byte. Mask/Payload len(7bits)
	h[1] = 0
	switch {
	case length <= 125:
		// 7 bits
		h[1] |= byte(length)
	case length < 65536:
		// 16 bits
		h[1] |= 126
		if h, err = c.w.Peek(2); err != nil {
			return
		}
		binary.BigEndian.PutUint16(h, uint16(length))
	default:
		// 64 bits
		h[1] |= 127
		if h, err = c.w.Peek(8); err != nil {
			return
		}
		binary.BigEndian.PutUint64(h, uint64(length))
	}
	return
}

func (c *Conn) WriteBody(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.w.Write(b)
	}
	return
}

func (c *Conn) Flush() error {
	return c.w.Flush()
}

func (c *Conn) Peek(n int) ([]byte, error) {
	return c.w.Peek(n)
}
