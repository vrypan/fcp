package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/vrypan/farcaster-go/farcaster"
)

func Inspect(filename string, opts map[string]any) {
	stats := opts["stats"].(bool)

	var f *os.File
	if filename == "-" {
		f = os.Stdin
	} else {
		var err error
		if f, err = os.Open(filename); err != nil {
			fmt.Println(err)
			return
		}
	}
	defer f.Close()

	count := 0
	dataTypes := make(map[int]int)

	allowedType := make(map[farcaster.MessageType]bool)
	allowedType[farcaster.MessageType_MESSAGE_TYPE_CAST_ADD] = opts["casts"].(bool)
	allowedType[farcaster.MessageType_MESSAGE_TYPE_REACTION_ADD] = opts["reactions"].(bool)
	allowedType[farcaster.MessageType_MESSAGE_TYPE_LINK_ADD] = opts["links"].(bool)

	for {
		messages, err := ReadData(f, nil)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		count += len(messages.Messages)
		dataTypes[int(messages.Messages[0].Data.GetType())] += len(messages.Messages)
		if !stats && allowedType[messages.Messages[0].Data.GetType()] {
			j, err := json.Marshal(messages.Messages)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(j))
		}
	}
	if stats {
		fmt.Printf("Total records: %d\n", count)
		for dataType, counts := range dataTypes {
			fmt.Printf("%s: %d\n", farcaster.MessageType_name[int32(dataType)], counts)
		}
	}
}
