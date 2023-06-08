package websocket

import (
	bufio2 "bufio"
	"encoding/base64"
	"fmt"
	"github.com/xyhubl/yim/pkg/bufio"
	"log"
	"math/rand"
	"net"
	url2 "net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestServer(t *testing.T) {
	var (
		data = []byte{0, 1, 2}
	)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.FailNow()
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
		}
		rd := bufio.NewReader(conn)
		wr := bufio.NewWriter(conn)
		req, err := ReadRequest(rd)
		if err != nil {
			t.Error(err)
		}
		if req.RequestURI != "/sub" {
			t.Error(err)
		}
		ws, err := Upgrade(conn, rd, wr, req)
		if err != nil {
			t.Error(err)
		}
		if err = ws.WriteMessage(BinaryMessage, data); err != nil {
			t.Error(err)
		}
		if err = ws.Flush(); err != nil {
			t.Error(err)
		}
		op, b, err := ws.ReadMessage()
		if err != nil || op != BinaryMessage || !reflect.DeepEqual(b, data) {
			t.Error(err)
		}
	}()
	time.Sleep(time.Millisecond * 100)
	// ws client
	ws, err := websocket.Dial("ws://127.0.0.1:8080/sub", "", "*")
	if err != nil {
		t.FailNow()
	}
	// receive binary frame
	var b []byte
	if err = websocket.Message.Receive(ws, &b); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(b, data) {
		t.FailNow()
	}
	// send binary frame
	if err = websocket.Message.Send(ws, data); err != nil {
		t.FailNow()
	}
}

func TestWs(t *testing.T) {
	var (
		wg   sync.WaitGroup
		err  error
		conn *net.TCPConn
		req  *Request
		ws   *Conn
	)
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	wg.Add(1)

	conn, err = l.AcceptTCP()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	go func() {
		r := bufio.NewReaderSize(conn, 32)
		w := bufio.NewWriterSize(conn, 32)

		req, err = ReadRequest(r)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if req.RequestURI != "/sub" {
			t.Errorf("wrong uri: %s", req.RequestURI)
			t.FailNow()
		}
		ws, err = Upgrade(conn, r, w, req)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		var (
			op      int
			payload []byte
		)
		for {
			op, payload, err = ws.ReadMessage()
			if err != nil {
				t.Error(err, op, payload)
				t.Fatal(err)
			}

			// ws.WriteMessage(TextMessage, []byte("00000000000000000000000000000001"))

			fmt.Println(string(payload), err, op)
		}
	}()
	wg.Wait()
}

func TestWs2(t *testing.T) {

	serverURL := "ws://127.0.0.1:8080/sub"

	// 解析WebSocket服务器地址
	u, err := url2.Parse(serverURL)
	if err != nil {
		log.Fatal("Failed to parse server URL:", err)
	}
	// 建立TCP连接
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	// 发送WebSocket协议升级请求
	key := generateWebSocketKey()
	request := createWebSocketUpgradeRequest(u, key)
	fmt.Fprintf(conn, "%s\r\n", request)

	// 读取和验证升级响应
	response, err := bufio2.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read server response:", err)
	}
	if !strings.HasPrefix(response, "HTTP/1.1 101 Switching Protocols") {
		log.Fatal("Failed to upgrade connection:", response)
	}
	//r := bufio.NewReaderSize(conn, 32)
	//w := bufio.NewWriterSize(conn, 32)

	time.Sleep(1 * time.Second)
	f := createWebSocketFrame("000000000001")
	conn.Write(f)

	//for {
	//	n, payload, err := c.ReadMessage()
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println(n, string(payload), "=======")
	//}
}

func createWebSocketUpgradeRequest(u *url2.URL, key string) string {
	request := fmt.Sprintf("GET %s HTTP/1.1\r\n", u.Path)
	request += fmt.Sprintf("Host: %s\r\n", u.Host)
	request += "Upgrade: websocket\r\n"
	request += "Connection: Upgrade\r\n"
	request += "Sec-WebSocket-Version: 13\r\n"
	request += fmt.Sprintf("Sec-WebSocket-Key: %s\r\n", key)
	request += "\r\n"
	return request
}

func generateWebSocketKey() string {
	nonce := make([]byte, 16)
	rand.Read(nonce)
	return base64.StdEncoding.EncodeToString(nonce)
}

func createWebSocketFrame(message string) []byte {
	// 设置帧的FIN位、操作码和负载数据
	finBit := byte(1 << 7)
	opCode := byte(0x01) // 文本消息操作码
	payloadData := []byte(message)

	// 创建WebSocket帧
	frame := make([]byte, 2+len(payloadData))

	// 设置第一个字节：FIN位和操作码
	frame[0] = finBit | opCode

	// 设置第二个字节：Mask位和负载数据长度
	maskBit := byte(1 << 7)
	payloadLen := len(payloadData)

	if payloadLen <= 125 {
		frame[1] = byte(payloadLen) | maskBit
	} else if payloadLen <= 65535 {
		frame[1] = 126 | maskBit
		frame[2] = byte(payloadLen >> 8)
		frame[3] = byte(payloadLen)
	} else {
		frame[1] = 127 | maskBit
		frame[2] = byte(payloadLen >> 56)
		frame[3] = byte(payloadLen >> 48)
		frame[4] = byte(payloadLen >> 40)
		frame[5] = byte(payloadLen >> 32)
		frame[6] = byte(payloadLen >> 24)
		frame[7] = byte(payloadLen >> 16)
		frame[8] = byte(payloadLen >> 8)
		frame[9] = byte(payloadLen)
	}

	// 设置掩码和负载数据
	mask := []byte{0x01, 0x02, 0x03, 0x04} // 用于对负载数据进行掩码处理
	copy(frame[len(frame)-len(payloadData):], payloadData)
	for i := 0; i < len(payloadData); i++ {
		frame[i+len(frame)-len(payloadData)] ^= mask[i%4]
	}

	return frame
}
