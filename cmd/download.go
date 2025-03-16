/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrypan/fckup/farcaster"
	"github.com/vrypan/fckup/fctools"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Create a backup",
	Run:   downloadCmdMain,
}

const PAGE_SIZE = 1000

func downloadCmdMain(cmd *cobra.Command, args []string) {
	currentDateTime := time.Now().Format("2006-01-02-1504-05")
	hubAddress, _ := cmd.Flags().GetString("hub")
	useSsl, _ := cmd.Flags().GetBool("ssl")
	dir, _ := cmd.Flags().GetString("dir")

	if len(args) == 0 {
		cmd.Help()
		return
	}
	fid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	hub := fctools.NewFarcasterHub(hubAddress, useSsl)
	defer hub.Close()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	messageTypes := []string{"reactions", "casts", "links"}
	for _, t := range messageTypes {
		outfilePath := filepath.Join(dir, fmt.Sprintf("%06d_%s_%s.backup", fid, currentDateTime, t))
		outfile, err := os.OpenFile(outfilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer outfile.Close()
		pageToken := []byte{}
		count := 0
		for {
			var response []*farcaster.Message
			var nextPageToken []byte
			switch t {
			case "reactions":
				response, nextPageToken, err = hub.GetReactionsByFid(uint64(fid), pageToken, PAGE_SIZE)
			case "casts":
				response, nextPageToken, err = hub.GetCastsByFid(uint64(fid), pageToken, PAGE_SIZE)
			case "links":
				response, nextPageToken, err = hub.GetLinksByFid(uint64(fid), pageToken, PAGE_SIZE)
			default:
				err = fmt.Errorf("unknown type: %s", t)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			j, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
				return
			}
			if _, err = outfile.WriteString(string(j) + "\n"); err != nil {
				fmt.Println(err)
				return
			}
			count += len(response)

			nextStr := base64.StdEncoding.EncodeToString(nextPageToken)
			fmt.Print("\033[2K\r")
			fmt.Printf("... Saving %d %s\t%s", count, t, nextStr)
			if len(nextPageToken) == 0 {
				break
			}
			pageToken = nextPageToken
		}
		fmt.Println(" Done.")
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("dir", "d", ".", "Path to export dir")
	downloadCmd.Flags().String("hub", "hoyt.farcaster.xyz:2283", "Farcaster hub to use")
	downloadCmd.Flags().Bool("ssl", true, "Use SSL")
}
