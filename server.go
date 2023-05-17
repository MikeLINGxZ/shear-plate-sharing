package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
)

var connList map[string]net.Conn
var lock sync.Mutex

func runServer() error {
	connList := make(map[string]net.Conn)
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Port))
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("accept conn error:", err.Error())
				continue
			}
			password, err := readBytes(conn)
			if err != nil {
				log.Println("read client password error:", err.Error())
				continue
			}
			if password != "@"+config.Password+"@" {
				log.Println("client password not match:", conn.RemoteAddr())
				continue
			}
			lock.Lock()
			id := conn.RemoteAddr().String() + "-" + uuid.New().String()[:5]
			connList[id] = conn
			lock.Unlock()
			go handlerConn(id, conn)
		}
	}()
	return nil
}

func handlerConn(id string, conn net.Conn) {
	for {
		content, err := readBytes(conn)
		if err != nil {
			log.Println("read client error:", err.Error())
			break
		}
		go sendContent(id, content)
	}
	lock.Lock()
	defer lock.Unlock()
	delete(connList, id)
}

func sendContent(id, content string) {
	for key, conn := range connList {
		if id == key {
			continue
		}
		_ = sendBytes(conn, content)
	}
}
