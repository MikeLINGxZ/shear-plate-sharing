package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net"
)

const headerLen int = 8

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

type TcpMsg struct {
	Name    string
	Content []byte
	Type    ContentType
	To      string `json:"-"`
}

type Tcp struct {
	role    RoleType
	conn    net.Conn
	id      string
	watchCh chan *TcpMsg
	log     *Log
}

func NewTcp(conn net.Conn, role RoleType, log *Log) *Tcp {
	return &Tcp{
		role: role,
		conn: conn,
		id:   conn.RemoteAddr().String() + "-" + uuid.New().String()[:5],
		log:  log,
	}
}

func (t *Tcp) Watch() <-chan *TcpMsg {
	if t.watchCh != nil {
		return t.watchCh
	}
	t.watchCh = make(chan *TcpMsg, 1)
	go func() {
		for {
			msg, err := t.Read()
			if err != nil {
				panic(err)
			}
			t.watchCh <- msg
		}
	}()
	return t.watchCh
}

func (t *Tcp) Read() (*TcpMsg, error) {
	return t.read()
}

func (t *Tcp) Send(msg *TcpMsg) error {
	return t.send(msg.Name, msg.Content, msg.Type)
}

func (t *Tcp) send(name string, contentBytes []byte, contentType ContentType) error {
	msg := &TcpMsg{
		Name:    name,
		Content: contentBytes,
		Type:    contentType,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	contentLen := len(msgBytes)
	contentLenBytes := Int64ToBytes(int64(contentLen))
	binLen, err := t.conn.Write(contentLenBytes)
	if err != nil {
		return err
	}
	if binLen != headerLen {
		return errors.New("msg len not match")
	}
	binLen, err = t.conn.Write(msgBytes)
	if err != nil {
		return err
	}
	if binLen != contentLen {
		return errors.New("content len not match")
	}
	msg.To = t.conn.RemoteAddr().String()
	t.log.LogSendMsg(msg)
	return nil
}

func (t *Tcp) read() (*TcpMsg, error) {
	// read content len
	lenInfoBytes := make([]byte, headerLen)
	binLen, err := t.conn.Read(lenInfoBytes)
	if err != nil {
		return nil, err
	}
	if binLen != headerLen {
		return nil, errors.New("msg len not match")
	}
	msgLen := BytesToInt64(lenInfoBytes)
	// read content
	msgBytes := make([]byte, msgLen)
	binLen, err = t.conn.Read(msgBytes)
	if err != nil {
		return nil, err
	}
	if int64(binLen) != msgLen {
		return nil, errors.New("content len not match")
	}
	msg := &TcpMsg{}
	err = json.Unmarshal(msgBytes, msg)
	if err != nil {
		return nil, err
	}
	msg.To = t.conn.RemoteAddr().String()
	t.log.LogReadMsg(msg)
	return msg, nil
}

func (t *Tcp) GetTcpID() string {
	return t.id
}

func (t *Tcp) Close() {
	if t.watchCh != nil {
		close(t.watchCh)
	}
	err := t.conn.Close()
	if err != nil {
		t.log.Log("close conn error: %s", err.Error())
	}
}
