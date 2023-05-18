package main

import (
	"context"
	"golang.design/x/clipboard"
	"net"
)

func runClient() {
	// connect server
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		panic(err)
	}
	tcp := NewTcp(conn, "client")
	// send password
	err = tcp.SendMsg([]byte("@" + config.Password + "@"))
	if err != nil {
		panic(err)
	}
	// clipboard
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	clipboardCh := clipboard.Watch(context.Background(), clipboard.FmtText)
	networkCh := tcp.Watch()
	clipboardHandler(tcp, clipboardCh, networkCh)
}

func clipboardHandler(tcp *Tcp, clipboardCh, networkCh <-chan []byte) {
	lastContent := clipboard.Read(clipboard.FmtText)
	for {
		var content []byte
		select {
		case content = <-clipboardCh:
			if string(content) == string(lastContent) {
				continue
			}
			err := tcp.SendMsg(content)
			if err != nil {
				panic(err)
			}
		case content = <-networkCh:
			if string(content) == string(lastContent) {
				continue
			}
			clipboard.Write(clipboard.FmtText, content)
		}
		lastContent = content
	}
}
