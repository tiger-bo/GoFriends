package xstream

import (
	"bufio"
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
)

func XStreamRead(r io.Reader, x proto.Message) error {

	bytesReader := bufio.NewReader(r)
	length, err := binary.ReadUvarint(bytesReader)
	if err != nil {
		panic(err)
	}

	msg := make([]byte, length)
	_, err = io.ReadFull(bytesReader, msg)
	if err != nil && err != io.EOF {
		panic(err)
	}

	err = proto.Unmarshal(msg, x)
	if err != nil {
		panic(err)
	}
	return err
}

func XStreamWrite(w io.Writer, x proto.Message) error {

	msg, err := proto.Marshal(x)
	if err != nil {
		panic(err)
	}

	msgLengthBytes := make([]byte, 20)
	n := binary.PutUvarint(msgLengthBytes, uint64(len(msg)))

	_, err = w.Write(msgLengthBytes[:n])
	if err != nil {
		panic(err)
	}

	_, err = w.Write(msg)
	if err != nil {
		panic(err)
	}

	return err
}
