package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
)

const headerLen int = 8

type Tcp struct {
	title string
	conn  net.Conn
	id    string
}

func NewTcp(conn net.Conn, title string) *Tcp {
	return &Tcp{
		title: title,
		conn:  conn,
		id:    conn.RemoteAddr().String() + "-" + uuid.New().String()[:5],
	}
}

func (t *Tcp) ReadMsg() ([]byte, error) {
	// read content len
	lenInfoBytes := make([]byte, headerLen)
	binLen, err := t.conn.Read(lenInfoBytes)
	if err != nil {
		return nil, err
	}
	if binLen != headerLen {
		return nil, errors.New("msg len not match")
	}
	contentLen := bytesToInt64(lenInfoBytes)
	// read content
	contentBytes := make([]byte, contentLen)
	binLen, err = t.conn.Read(contentBytes)
	if err != nil {
		return nil, err
	}
	if int64(binLen) != contentLen {
		return nil, errors.New("content len not match")
	}
	t.log("read msg: %s", string(contentBytes))
	return contentBytes, nil
}

func (t *Tcp) SendMsg(contentBytes []byte) error {
	contentLen := len(contentBytes)
	contentLenBytes := int64ToBytes(int64(contentLen))
	binLen, err := t.conn.Write(contentLenBytes)
	if err != nil {
		return err
	}
	if binLen != headerLen {
		return errors.New("msg len not match")
	}
	binLen, err = t.conn.Write(contentBytes)
	if err != nil {
		return err
	}
	if binLen != contentLen {
		return errors.New("content len not match")
	}
	t.log("send msg: %s", string(contentBytes))
	return nil
}

func (t *Tcp) GetTcpID() string {
	return t.id
}

func (t *Tcp) log(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	log.Printf("[%s-%s] %s \n", t.title, t.id, msg)
}

func int64ToBytes(num int64) []byte {
	byteArray := make([]byte, headerLen)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))

	return byteArray
}

func bytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}
