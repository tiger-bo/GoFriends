package xstream

import (
	pb "GoFriends/chat/chatmsg"
	"bytes"
	"testing"
)

func TestXStreamReadAndWrite(t *testing.T) {

	// Test case 1: Normal case with a valid message
	message1 := &pb.Message{
		Id:      1,
		UserId:  1,
		Content: "hello",
	}
	var buffer1 bytes.Buffer

	XStreamWrite(&buffer1, message1)

	var message1Read pb.Message
	err1 := XStreamRead(&buffer1, &message1Read)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	if message1.Id != message1Read.Id || message1.UserId != message1Read.UserId || message1.Content != message1Read.Content {
		t.Errorf("Test case 1 failed: expected %v, got %v", message1, message1Read)

	}

}
