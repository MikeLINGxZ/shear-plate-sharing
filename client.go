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

	// write or read
	go func() {
		for {
			content := <-clipboardCh
			err := tcp.SendMsg(content)
			if err != nil {
				panic(err)
			}
		}
	}()
	for {
		content, err := tcp.ReadMsg()
		if err != nil {
			panic(err)
		}
		changed := clipboard.Write(clipboard.FmtText, content)
		<-changed
	}
}
