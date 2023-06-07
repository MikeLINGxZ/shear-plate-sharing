package main

import (
	"encoding/binary"
	"log"
)

func Int64ToBytes(num int64) []byte {
	byteArray := make([]byte, 8)
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
