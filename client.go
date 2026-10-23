package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type ChatMessage struct {
	User string
	Body string
}

func main() {
	c, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Connection failed!")
		return
	}

	in := bufio.NewReader(os.Stdin)
	fmt.Print("Your name: ")
	myName, _ := in.ReadString('\n')
	myName = strings.TrimSpace(myName)

	var answer string
	c.Call("MyServer.Join", myName, &answer)
	if answer == "EXISTS" {
		fmt.Println("Name already taken.")
		return
	}

	go func() {
		seen := 0
		for {
			var list []ChatMessage
			c.Call("MyServer.GetMessages", struct{}{}, &list)

			if len(list) > seen {
				for i := seen; i < len(list); i++ {
					if list[i].User == "System" {
						fmt.Printf("\n[NOTIFICATION]: %s\n> ", list[i].Body)
					} else if list[i].User != myName {
						fmt.Printf("\n[%s]: %s\n> ", list[i].User, list[i].Body)
					}
				}
				seen = len(list)
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		fmt.Print("> ")
		line, _ := in.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		if line == "exit" {
			var ok bool
			c.Call("MyServer.Leave", myName, &ok)
			fmt.Println("Exiting chat...")
			break
		}

		m := ChatMessage{User: myName, Body: line}
		var success bool
		c.Call("MyServer.Send", m, &success)
	}
}
