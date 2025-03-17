package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/vrypan/fcp/farcaster"
)

func Inspect(filename string, opts map[string]any) {
	stats := opts["stats"].(bool)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	count := 0
	dataTypes := make(map[int]int)
	for {
		messages, err := ReadData(f)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		count += len(messages.Messages)
		dataTypes[int(messages.Messages[0].Data.GetType())] += len(messages.Messages)
		if !stats {
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
