package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/vrypan/fckup/farcaster"
	"github.com/vrypan/fckup/fctools"
)

func Download(hubAddress string, username string, localFile string, opts map[string]any) {
	useSsl, _ := opts["ssl"].(bool)
	pageSize, ok := opts["pageSize"].(uint32)
	if !ok {
		pageSize = 100
	}

	hub := fctools.NewFarcasterHub(hubAddress, useSsl)
	defer hub.Close()

	var fid uint64
	if fidInt, err := strconv.Atoi(username); err == nil {
		fid = uint64(fidInt)
	} else if retrievedFid, err := hub.GetFidByUsername(username); err != nil {
		fmt.Printf("Error getting fid: %v\n", err)
		return
	} else {
		fid = retrievedFid
	}

	outfile, err := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()

	performDataFetch := map[string]func(uint64, []byte, uint32) (*farcaster.MessagesResponse, error){
		"reactions": hub.GetReactionsByFid,
		"casts":     hub.GetCastsByFid,
		"links":     hub.GetLinksByFid,
	}

	for messageType, hubFunction := range performDataFetch {
		pageToken := []byte{}
		count := 0
		for {
			response, err := hubFunction(fid, pageToken, pageSize)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := WriteData(outfile, response); err != nil {
				fmt.Println()
				fmt.Println(err)
				fmt.Println()
				return
			}
			count += len(response.Messages)

			nextStr := base64.StdEncoding.EncodeToString(response.NextPageToken)
			fmt.Printf("\033[2K\r... Saving %06d %-12s %s", count, messageType, nextStr)
			if len(response.NextPageToken) == 0 {
				break
			}
			pageToken = response.NextPageToken
		}
		fmt.Println("Done.")
	}
}

func Upload(hubAddress string, localFile string, opts map[string]any) {
	useSsl, _ := opts["ssl"].(bool)

	hub := fctools.NewFarcasterHub(hubAddress, useSsl)
	defer hub.Close()

	signerPrivateKey, _ := opts["signer"].(string)
	var signer []byte
	var err error
	if signerPrivateKey != "" {
		signerPrivateKey = strings.TrimPrefix(signerPrivateKey, "0x")
		if signer, err = hex.DecodeString(signerPrivateKey); err != nil {
			fmt.Println("Invalid private key")
			return
		}
	}

	f, err := os.Open(localFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	count := 0
	errorCount := 0
	successCount := 0
	for {
		messages, err := ReadData(f)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		count += len(messages.Messages)
		for _, m := range messages.Messages {
			if len(signer) != 0 {
				m = fctools.ResignMessage(m, signer)
			}
			_, err := hub.SubmitMessage(m)
			hash := base64.StdEncoding.EncodeToString(m.GetHash())
			if err != nil {
				errorCount += 1
				fmt.Printf("%s failed: %s\n", hash, err.Error())
			} else {
				successCount += 1
				fmt.Print("\033[2K\r")
				fmt.Printf("%s Uploaded", hash)
			}
		}
	}
	fmt.Println()
	fmt.Printf("Total:   %6d messages\n", count)
	fmt.Printf("Success: %6d messages\n", successCount)
	fmt.Printf("Error:   %6d messages\n", errorCount)
}
