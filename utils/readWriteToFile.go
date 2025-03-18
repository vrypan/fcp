package utils

import (
	"encoding/binary"
	"os"

	"github.com/vrypan/fcp/farcaster"
	"google.golang.org/protobuf/proto"
)

func WriteData(f *os.File, messages *farcaster.MessagesResponse, opts map[string]any) error {
	data, err := proto.Marshal(messages)
	if err != nil {
		panic(err)
	}
	length := uint32(len(data))
	// Write length
	err = binary.Write(f, binary.LittleEndian, length)
	if err != nil {
		panic(err)
	}
	// Write data
	_, err = f.Write(data)
	return err
}

func ReadData(f *os.File, opts map[string]any) (*farcaster.MessagesResponse, error) {
	var length uint32
	var messages farcaster.MessagesResponse
	err := binary.Read(f, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	data := make([]byte, length)
	_, err = f.Read(data)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(data, &messages)
	if err != nil {
		panic(err)
	}
	return &messages, err
}
