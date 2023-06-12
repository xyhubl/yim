package protocol

import (
	"errors"
	"github.com/xyhubl/yim/pkg/encoding/binary"
	"github.com/xyhubl/yim/pkg/websocket"
)

var (
	ErrProtoPackLen   = errors.New("protocol: default server codec pack length error")
	ErrProtoHeaderLen = errors.New("protocol: default server codec header length error")
)

var (
	ProtoReady  = &Proto{Op: OpProtoReady}
	ProtoFinish = &Proto{Op: OpProtoFinish}
)

const (
	MaxBodySize = int32(1 << 12)

	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)

	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf       []byte
		headerLen int16
		packLen   int32
		bodyLen   int
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])

	if packLen < 0 || packLen > _maxPackSize {
		return ErrProtoHeaderLen
	}

	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

func (p *Proto) WriteWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + len(p.Body)
	if err = ws.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = ws.Peek(_rawHeaderSize); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))          // 4
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize)) // 2
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))             // 2
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)                      // 4
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)                    // 4
	if p.Body != nil {
		err = ws.WriteBody(p.Body)
	}
	return
}

func (p *Proto) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + _heartSize
	if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		return
	}
	if buf, err = wr.Peek(packLen); err != nil {
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	return
}
