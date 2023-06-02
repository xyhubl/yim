package websocket

import (
	"fmt"
	"github.com/xyhubl/yim/pkg/bufio"
	"net"
	"reflect"
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
		conn net.Conn
		req  *Request
		ws   *Conn
	)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	wg.Add(1)

	conn, err = l.Accept()
	fmt.Println("11")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	go func() {
		r := bufio.NewReaderSize(conn, 1024)
		w := bufio.NewWriterSize(conn, 1024)

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
			}
			fmt.Println(string(payload), err, op)
		}
	}()
	wg.Wait()
}
