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

func TestServerV1(t *testing.T) {
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
				t.Fatal(err)
			}
			fmt.Println(string(payload), err, op)

		}
	}()
	wg.Wait()
}

func TestClient(t *testing.T) {

	serverURL := "ws://127.0.0.1:8080/sub"

	// 解析WebSocket服务器地址
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8080")

	u, err := url2.Parse(serverURL)
	if err != nil {
		log.Fatal("Failed to parse server URL:", err)
	}
	// 建立TCP连接
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	// 发送WebSocket协议升级请求
	key := generateWebSocketKey()
	request := createWebSocketUpgradeRequest(u, key)
	fmt.Fprintf(conn, "%s", request)
	// 读取和验证升级响应
	b := bufio2.NewReader(conn)
	response, err := b.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read server response:", err)
	}
	if !strings.HasPrefix(response, "HTTP/1.1 101 Switching Protocols") {
		log.Fatal("Failed to upgrade connection:", response)
	}
	for {
		lineB, _, _ := b.ReadLine()
		if string(lineB) == "" {
			break
		}
	}

	r := bufio.NewReaderSize(conn, 32)
	w := bufio.NewWriterSize(conn, 32)

	c := newConn(conn, r, w)
	for {
		if err = c.WriteMessage(PingMessage, []byte("pong")); err != nil {
			t.Fatal(err)
		}
		if err = c.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	c.Close()
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
