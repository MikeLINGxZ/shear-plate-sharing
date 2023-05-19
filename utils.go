package main

import (
	"encoding/binary"
	"log"
)

func Int64ToBytes(num int64) []byte {
	byteArray := make([]byte, headerLen)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))

	return byteArray
}

func BytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}

type Log struct {
	role RoleType
	id   string
}

func NewLog(role RoleType) *Log {
	return &Log{
		role: role,
	}
}

func (l *Log) SetId(id string) {
	l.id = id
}

func (l *Log) Log(format string, v ...any) {
	log.Printf(format+"\n", v...)
}

func (l *Log) LogSendMsg(msg *TcpMsg) {
	l.logMsg("send", msg)
}

func (l *Log) LogReadMsg(msg *TcpMsg) {
	l.logMsg("read", msg)
}

func (l *Log) logMsg(text string, msg *TcpMsg) {
	var content string
	if msg.Type == CTFile || msg.Type == CTImg {
		content = msg.Name
	} else {
		content = string(msg.Content)
	}
	log.Printf("[%s] %s msg: \nip: %s\ntype: %d\ncontent: %s\n\n", l.role, text, msg.To, msg.Type, content)
}
