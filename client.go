package main

import (
	"context"
	"golang.design/x/clipboard"
	"net"
)

var clientLog *Log
var lastContent []byte
var msgHandler map[ContentType]func(msg *TcpMsg) error
var isServer bool

func runClient() {
	// init log
	clientLog = NewLog(RTClient)
	// init msg Handler
	initMsgHandler()
	// connect server
	if config.Role == "server" {
		config.Host = "127.0.0.1"
	}
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		clientLog.Log("connect server tcp error: %s", err.Error())
		panic(err)
	}
	tcp := NewTcp(conn, "client", clientLog)
	// send password
	err = tcp.Send(&TcpMsg{
		Name:    "",
		Content: []byte(config.Password),
		Type:    CTPassword,
	})
	if err != nil {
		clientLog.Log("send password error: %s", err.Error())
		panic(err)
	}
	// clipboard
	err = clipboard.Init()
	if err != nil {
		clientLog.Log("init clipboard error: %s", err.Error())
		panic(err)
	}

	clipboardCh := clipboard.Watch(context.Background(), clipboard.FmtText)
	msgCh := tcp.Watch()
	clipboardHandler(tcp, clipboardCh, msgCh)
}

func clipboardHandler(tcp *Tcp, clipboardCh <-chan []byte, msgCh <-chan *TcpMsg) {
	defer tcp.Close()
	lastContent = clipboard.Read(clipboard.FmtText)
	for {
		var content []byte
		select {
		case content = <-clipboardCh:
			if string(content) == string(lastContent) {
				continue
			}
			err := tcp.Send(&TcpMsg{
				Content: content,
				Type:    CTText,
			})
			if err != nil {
				clientLog.Log("send msg error: %s", err.Error())
				panic(err)
			}
			lastContent = content
		case msg := <-msgCh:
			f, ok := msgHandler[msg.Type]
			if !ok {
				continue
			}
			err := f(msg)
			if err != nil {
				clientLog.Log("handler msg error: %s", err.Error())
				panic(err)
			}
		}
	}
}

func initMsgHandler() {
	msgHandler = make(map[ContentType]func(msg *TcpMsg) error)
	msgHandler[CTText] = handlerText
}

func handlerText(msg *TcpMsg) error {
	if string(msg.Content) == string(lastContent) {
		return nil
	}
	clipboard.Write(clipboard.FmtText, msg.Content)
	lastContent = msg.Content
	return nil
}
