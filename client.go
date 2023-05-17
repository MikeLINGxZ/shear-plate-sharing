package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"net"
	"time"
)

var lastContent string

func runClient() {
	content, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}
	lastContent = content
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		panic(err)
	}
	err = sendBytes(conn, "@"+config.Password+"@")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			content, err := clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
			if content != lastContent {
				err := sendBytes(conn, content)
				if err != nil {
					panic(err)
				}
				lastContent = content
			}
			time.Sleep(time.Second * 2)
		}
	}()
	for {
		content, err := readBytes(conn)
		fmt.Println("content:", content)
		if err != nil {
			panic(err)
		}
		err = clipboard.WriteAll(content)
		if err != nil {
			panic(err)
		}
		lastContent = content
		time.Sleep(time.Second * 2)
	}
}
