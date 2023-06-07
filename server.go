package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

var connList map[string]*Tcp
var lock sync.Mutex
var serverLog *Log

func runServer() error {
	isServer = true
	// init log
	serverLog = NewLog(RTServer)
	// init conn list
	connList = make(map[string]*Tcp)
	// listen tcp
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Port))
	if err != nil {
		return err
	}
	// handle conn
	go func() {
		for {
			// accept conn
			conn, err := ln.Accept()
			if err != nil {
				serverLog.Log("accept conn error:", err.Error())
				continue
			}
			tcp := NewTcp(conn, config.Role, serverLog)
			// verify password
			err = verifyPassword(tcp)
			if err != nil {
				serverLog.Log("client password verify failed:", err.Error())
				continue
			}
			// add conn list
			lock.Lock()
			connList[tcp.GetTcpID()] = tcp
			lock.Unlock()
			// listen msg
			go listenMsg(tcp)
		}
	}()
	return nil
}

func verifyPassword(tcp *Tcp) error {
	contentType, contentBytes, err := tcp.Read()
	if err != nil {
		return err
	}
	if contentType != CTPassword {
		return errors.New("verify password fail")
	}
	if string(contentBytes) != config.Password {
		err := tcp.Send(CTSystem, []byte("403"))
		if err != nil {
			return err
		}
		return errors.New("verify password fail")
	}
	return nil
}

func listenMsg(tcp *Tcp) {
	defer tcp.Close()
	for {
		contentType, contentBytes, err := tcp.Read()
		if err != nil {
			continue
		}
		notifyMsg(tcp.GetTcpID(), contentType, contentBytes)
	}
	lock.Lock()
	delete(connList, tcp.GetTcpID())
	lock.Unlock()
}

func notifyMsg(id string, contentType ContentType, contentBytes []byte) {
	for _, tcp := range connList {
		tcp := tcp
		if tcp.GetTcpID() == id {
			continue
		}
		err := tcp.Send(contentType, contentBytes)
		if err != nil {
			tcp.log.Log("notifyMsg msg error:%s", err.Error())
		}
	}
}
