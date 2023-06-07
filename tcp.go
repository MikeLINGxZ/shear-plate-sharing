package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net"
)

const headerLen int = 16

type ContentType int

const (
	CTUnknown ContentType = iota
	CTSystem
	CTPassword
	CTText
	CTImg
	CTFile
)

type SystemContent struct {
	Text string
	Code int
}

func (sc *SystemContent) Bytes() []byte {
	bytes, _ := json.Marshal(sc)
	return bytes
}

type Tcp struct {
	role   RoleType
	conn   net.Conn
	id     string
	dataCh chan clipboardData
	log    *Log
}

func NewTcp(conn net.Conn, role RoleType, log *Log) *Tcp {
	return &Tcp{
		role: role,
		conn: conn,
		id:   conn.RemoteAddr().String() + "-" + uuid.New().String()[:5],
		log:  log,
	}
}

func (t *Tcp) Read() (ContentType, []byte, error) {
	return t.read()
}

func (t *Tcp) Send(contentType ContentType, contentBytes []byte) error {
	return t.send(contentType, contentBytes)
}

func (t *Tcp) send(contentType ContentType, contentBytes []byte) error {
	header := t.GenHeader(len(contentBytes), contentType)
	_, err := t.conn.Write(header)
	if err != nil {
		return err
	}
	_, err = t.conn.Write(contentBytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tcp) read() (ContentType, []byte, error) {
	// read content len
	headerBytes := make([]byte, headerLen)
	_, err := t.conn.Read(headerBytes[:])
	if err != nil {
		return 0, nil, err
	}
	size, contentType := t.PauseHeader(headerBytes)
	contentBytes := make([]byte, size)
	_, err = t.conn.Read(contentBytes[:])
	if err != nil {
		return 0, nil, err
	}

	return contentType, contentBytes, nil
}

func (t *Tcp) GetTcpID() string {
	return t.id
}

func (t *Tcp) Close() {
	err := t.conn.Close()
	if err != nil {
		t.log.Log("close conn error: %s", err.Error())
	}
}

func (t *Tcp) ClipboardDataWatch() chan clipboardData {
	if t.dataCh != nil {
		return t.dataCh
	}
	t.dataCh = make(chan clipboardData, 1)
	go func() {
		for {
			data := clipboardData{}
			var err error
			data.ContentType, data.Content, err = t.Read()
			if err != nil {
				panic(err)
			}
			t.dataCh <- data
		}
	}()

	return t.dataCh
}

func (t *Tcp) GenHeader(size int, contentType ContentType) []byte {
	contentLenBytes := Int64ToBytes(int64(size))
	contentTypeBytes := Int64ToBytes(int64(contentType))
	contentLenBytes = append(contentLenBytes, contentTypeBytes...)
	return contentLenBytes
}

func (t *Tcp) PauseHeader(header []byte) (int, ContentType) {
	contentLenBytes := header[:8]
	contentTypeBytes := header[8:]
	size := BytesToInt64(contentLenBytes)
	contentType := ContentType(BytesToInt64(contentTypeBytes))
	return int(size), contentType
}
