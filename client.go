package main

import (
	"context"
	"golang.design/x/clipboard"
	"net"
)

var clientLog *Log
var lastContent []byte
var msgHandler map[ContentType]func(contentType ContentType, contentBytes []byte) error
var isServer bool

func runClient() {
	// init log
	clientLog = NewLog(RTClient)
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
	err = tcp.Send(CTPassword, []byte(config.Password))
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

	clipboardCh := clipboardWatch()
	tcpCh := tcp.ClipboardDataWatch()
	dataCh := make(chan clipboardData, 1)
	go func() {
		for {
			select {
			case e := <-clipboardCh:
				dataCh <- e
			case e := <-tcpCh:
				dataCh <- e
			}
		}
	}()
	for {
		e := <-dataCh
		if string(e.Content) == string(lastContent) {
			continue
		}
		switch e.ContentType {
		case CTImg:
			clipboard.Write(clipboard.FmtImage, e.Content)
		case CTText:
			clipboard.Write(clipboard.FmtText, e.Content)
		}
		lastContent = e.Content
	}
}

type clipboardData struct {
	ContentType ContentType
	Content     []byte
}

func clipboardWatch() chan clipboardData {
	textCh := clipboard.Watch(context.Background(), clipboard.FmtText)
	imgCh := clipboard.Watch(context.Background(), clipboard.FmtImage)
	ch := make(chan clipboardData, 1)
	go func() {
		for {
			data := clipboardData{}
			select {
			case content := <-textCh:
				data.ContentType = CTText
				data.Content = content
			case content := <-imgCh:
				data.ContentType = CTImg
				data.Content = content
			}
		}
	}()
	return ch
}

func handler(clipboardData) {

}
