package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
)

const headerLen int = 8

type ContentType int

const (
	CTUnknown ContentType = iota
	CTPassword
	CTText
	CTImg
	CTFile
)

type TcpMsg struct {
	Name    string
	Content []byte
	Type    ContentType
}

type Tcp struct {
	title   string
	conn    net.Conn
	id      string
	watchCh chan *TcpMsg
}

func NewTcp(conn net.Conn, title string) *Tcp {
	return &Tcp{
		title: title,
		conn:  conn,
		id:    conn.RemoteAddr().String() + "-" + uuid.New().String()[:5],
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
	contentLenBytes := int64ToBytes(int64(contentLen))
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
	t.logMsg("send", msg)
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
	msgLen := bytesToInt64(lenInfoBytes)
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
	t.logMsg("read", msg)
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
		t.log("close conn error: %s", err.Error())
	}
}

func (t *Tcp) log(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	log.Printf("[%s-%s] %s \n", t.title, t.id, msg)
}

func (t *Tcp) logMsg(text string, msg *TcpMsg) {
	var content []byte
	if msg.Type == CTText || msg.Type == CTPassword {
		content = msg.Content
	}
	logStr := fmt.Sprintf("role:%s \nmsg type:%d \nmsg name: %s\nmsg content: %s\n", text, msg.Type, msg.Name, string(content))
	log.Printf("[%s-%s] \n%s \n", t.title, t.id, logStr)
}

func int64ToBytes(num int64) []byte {
	byteArray := make([]byte, headerLen)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))

	return byteArray
}

func bytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}
