package main

import (
	"GoFriends/chat"
	pb "GoFriends/chat/chatmsg"
	"GoFriends/xstream"
	"bufio"
	"fmt"
	"hash/fnv"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Print("Enter user name: ")
	userName, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	loginMsg := &pb.LoginRequest{RoomId: 1}
	loginMsg.User.Name = userName
	hash64 := fnv.New32()
	hash64.Write([]byte(userName))
	loginMsg.User.Id = int32(hash64.Sum32())

	err = chat.LoginUser(conn, loginMsg)
	if err != nil {
		fmt.Println("Error logging in:", err.Error())
		return
	}

	go func() {
		msgReq := &pb.Message{UserId: loginMsg.User.Id}
		for {
			msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			msgReq.Content = msg
			chat.SendMessage(conn, msgReq)
		}
	}()

	go func() {
		for {
			msgReq := &pb.Message{}
			xstream.XStreamRead(conn, msgReq)
			fmt.Println(msgReq.Content)
		}
	}()

	select {}
}
