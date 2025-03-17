package utils

import (
	"bufio"
	"encoding/binary"
	"os"
	"strings"

	"github.com/vrypan/fcp/farcaster"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteData(f *os.File, messages *farcaster.MessagesResponse, opts map[string]any) error {
	json := opts["json"].(bool)

	if json {
		data, err := protojson.Marshal(messages)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(string(data) + "\n")
		return nil
	} else {
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
}

func ReadBinaryData(f *os.File, opts map[string]any) (*farcaster.MessagesResponse, error) {
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

func ReadJsonData(f *os.File, opts map[string]any) (*farcaster.MessagesResponse, error) {
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\n")
		var messages farcaster.MessagesResponse
		err := protojson.Unmarshal([]byte(line), &messages)
		if err != nil {
			return nil, err
		}
		return &messages, nil
	} else {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
