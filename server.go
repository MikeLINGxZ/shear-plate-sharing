package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var connList map[string]*Tcp
var lock sync.Mutex

func runServer() error {
	connList = make(map[string]*Tcp)
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

			tcp := NewTcp(conn, config.Role)

			err = verifyPassword(tcp)
			if err != nil {
				log.Println("verify password failed:", err.Error())
				continue
			}

			lock.Lock()
			connList[tcp.GetTcpID()] = tcp
			lock.Unlock()

			go handlerConn(tcp)
		}
	}()
	return nil
}

func verifyPassword(tcp *Tcp) error {
	password, err := tcp.ReadMsg()
	if err != nil {
		return err
	}
	if password != "@"+config.Password+"@" {
		return errors.New("password not match")
	}
	return nil
}

func handlerConn(tcp *Tcp) {
	for {
		content, err := tcp.ReadMsg()
		if err != nil {
			log.Println(tcp.GetTcpID(), "read msg error:", err.Error())
			break
		}
		go notify(content)
	}
	lock.Lock()
	delete(connList, tcp.GetTcpID())
	lock.Unlock()
}

func notify(content string) {
	for _, tcp := range connList {
		tcp := tcp
		err := tcp.SendMsg(content)
		if err != nil {
			tcp.log("notify msg error:%s", err.Error())
		}
	}
}
