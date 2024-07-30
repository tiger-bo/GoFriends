package chat

import (
	pb "GoFriends/chat/chatmsg"
	"GoFriends/xstream"
	"errors"
	"fmt"
	"net"
	"sync"
)

type ChatRoomState struct {
	Users    map[int32]*pb.User
	Messages []*pb.Message
}

type Server struct {
	mu      sync.Mutex
	rooms   map[int32]*ChatRoomState
	streams map[int32]map[int32]net.Conn
}

func newServer() *Server {
	return &Server{
		rooms:   make(map[int32]*ChatRoomState),
		streams: make(map[int32]map[int32]net.Conn),
	}
}

func Init() {
	// Start listening for incoming TCP connections
	PORT := ":1234"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	fmt.Println("Listening on ", PORT)

	server := newServer()
	go func() {
		for {
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err.Error())
				return
			}
			// Handle the connection in a new goroutine
			go server.handleRequest(conn)
		}
	}()

}

func (s *Server) Login(conn net.Conn) (*pb.LoginRequest, error) {

	loginReq := &pb.LoginRequest{}
	err := xstream.XStreamRead(conn, loginReq)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[loginReq.RoomId]
	if !ok {
		room = &ChatRoomState{
			Users:    make(map[int32]*pb.User),
			Messages: make([]*pb.Message, 0),
		}
		s.rooms[loginReq.RoomId] = room
		room.Users[loginReq.User.Id] = loginReq.User
		s.streams[loginReq.RoomId] = make(map[int32]net.Conn)
		s.streams[loginReq.RoomId][loginReq.User.Id] = conn
	} else {
		room.Users[loginReq.User.Id] = loginReq.User
		s.streams[loginReq.RoomId][loginReq.User.Id] = conn
	}

	return loginReq, nil
}

func SendMessage(conn net.Conn, msg *pb.Message) error {
	return xstream.XStreamWrite(conn, msg)
}

func (s *Server) SendRoomState(conn net.Conn, roomState *pb.ChatRoomState) error {
	return xstream.XStreamWrite(conn, roomState)
}

func (s *Server) SendMessageToRoom(roomId int32, msg *pb.Message) error {
	for _, conn := range s.streams[roomId] {
		err := SendMessage(conn, msg)
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
		}
	}
	return nil
}

func (s *Server) SendRoomStateToRoom(roomId int32, roomState *pb.ChatRoomState) error {
	for _, conn := range s.streams[roomId] {
		err := s.SendRoomState(conn, roomState)
		if err != nil {
			fmt.Println("Error sending room state:", err.Error())
		}
	}
	return nil
}

func (s *Server) ListenUserMessages(conn net.Conn, roomId int32) {

	for {
		msg := &pb.Message{}
		err := xstream.XStreamRead(conn, msg)
		if err != nil {
			fmt.Println("Error reading message:", err.Error())
			return
		}
		s.SendMessageToRoom(roomId, msg)
	}

}

// Handle incoming requests
func (s *Server) handleRequest(conn net.Conn) {

	loginReq, err := s.Login(conn)
	if err != nil {
		fmt.Println("Error logging in:", err.Error())
		return
	}

	go s.ListenUserMessages(conn, loginReq.RoomId)
}

func LoginUser(conn net.Conn, loginReq *pb.LoginRequest) error {
	err := xstream.XStreamWrite(conn, loginReq)
	if err != nil {
		return err
	}

	loginRes := &pb.LoginResponse{}
	err = xstream.XStreamRead(conn, loginRes)
	if err != nil {
		return err
	}

	if !loginRes.Success {
		fmt.Println("Login failed:", loginRes.Message)
		return errors.New(loginRes.Message)
	}

	return nil
}
