package main

import (
	"context"
	"golang.design/x/clipboard"
	"net"
)

func runClient() {
	// connect server
	if config.Role == "server" {
		config.Host = "127.0.0.1"
	}
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		panic(err)
	}
	tcp := NewTcp(conn, "client")
	// send password
	err = tcp.Send(&TcpMsg{
		Name:    "",
		Content: []byte("@" + config.Password + "@"),
		Type:    CTPassword,
	})
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

func clipboardHandler(tcp *Tcp, clipboardCh <-chan []byte, msgCh <-chan *TcpMsg) {
	defer tcp.Close()
	lastContent := clipboard.Read(clipboard.FmtText)
	for {
		var content []byte
		select {
		case content = <-clipboardCh:
			if string(content) == string(lastContent) {
				continue
			}
			err := tcp.Send(&TcpMsg{
				Name:    "",
				Content: content,
				Type:    CTText,
			})
			if err != nil {
				panic(err)
			}
		case msg := <-msgCh:
			if msg.Type == CTText && string(msg.Content) == string(lastContent) {
				continue
			}
			if msg.Type != CTText {
				continue
			}
			clipboard.Write(clipboard.FmtText, content)
		}
		lastContent = content
	}
}
