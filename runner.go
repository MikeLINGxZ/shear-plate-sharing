package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func run() error {
	InitConfig()
	if config.Role == "server" {
		err := runServer()
		if err != nil {
			return err
		}
	}
	runClient()
	return nil
}

func sendBytes(conn net.Conn, content string) error {
	contentBytes := []byte(content)
	contentLen := int64(len(contentBytes))
	sizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBytes, uint32(contentLen))
	_, err := conn.Write(sizeBytes)
	if err != nil {
		return err
	}

	fmt.Println(conn.RemoteAddr().String(), "send contentLen:", contentLen)

	l, err := conn.Write(contentBytes)
	if err != nil {
		return err
	}
	fmt.Println(conn.RemoteAddr().String(), "send l:", l)
	fmt.Println(conn.RemoteAddr().String(), "send content:", string(contentBytes))
	return nil
}

func readBytes(conn net.Conn) (string, error) {
	var bytes [4]byte
	_, err := conn.Read(bytes[:])
	if err != nil {
		return "", err
	}
	contentLen := int32(binary.LittleEndian.Uint32(bytes[:]))
	fmt.Println(conn.RemoteAddr().String(), "read contentLen:", contentLen)

	msgBytes := make([]byte, contentLen)
	l, err := conn.Read(msgBytes[:])
	if err != nil {
		return "", err
	}
	fmt.Println(conn.RemoteAddr().String(), "read l:", l)
	fmt.Println(conn.RemoteAddr().String(), "read content:", string(msgBytes))
	return string(msgBytes), nil
}
