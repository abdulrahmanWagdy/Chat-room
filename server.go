package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type ChatMessage struct {
	User string
	Body string
}

type MyServer struct {
	allMsgs []ChatMessage
	online  []string
}

func (s *MyServer) Join(name string, res *string) error {
	for _, u := range s.online {
		if u == name {
			*res = "EXISTS"
			return nil
		}
	}
	s.online = append(s.online, name)

	joinMsg := ChatMessage{User: "System", Body: name + " has joined the chat."}
	s.allMsgs = append(s.allMsgs, joinMsg)

	*res = "WELCOME"
	fmt.Println("User joined:", name)
	return nil
}

func (s *MyServer) Leave(name string, ok *bool) error {
	for i, u := range s.online {
		if u == name {
			s.online = append(s.online[:i], s.online[i+1:]...)
			break
		}
	}

	leaveMsg := ChatMessage{User: "System", Body: name + " has left the chat."}
	s.allMsgs = append(s.allMsgs, leaveMsg)

	fmt.Println("User left:", name)
	*ok = true
	return nil
}

func (s *MyServer) Send(m ChatMessage, ok *bool) error {
	s.allMsgs = append(s.allMsgs, m)
	*ok = true
	return nil
}

func (s *MyServer) GetMessages(_ struct{}, res *[]ChatMessage) error {
	*res = s.allMsgs
	return nil
}

func main() {
	srv := new(MyServer)
	rpc.Register(srv)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Server error")
		return
	}

	fmt.Println("Server running on port 1234...")
	http.Serve(l, nil)
}
