package main

import (
	"github.com/atotto/clipboard"
	"net"
)

var oldContent string

func runClient() {
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		panic(err)
	}

	tcp := NewTcp(conn, "client")
	// send password
	err = tcp.SendMsg("@" + config.Password + "@")
	if err != nil {
		panic(err)
	}

	// listen shear plate
	go func() {
		for {
			newContent, err := clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
			if newContent == oldContent {
				continue
			}
			err = tcp.SendMsg(newContent)
			if err != nil {
				panic(err)
			}
			oldContent = newContent
		}
	}()

	for {
		newContent, err := tcp.ReadMsg()
		if err != nil {
			panic(err)
		}
		if newContent == oldContent {
			continue
		}
		err = clipboard.WriteAll(newContent)
		if err != nil {
			panic(err)
		}
		oldContent = newContent
	}
}
