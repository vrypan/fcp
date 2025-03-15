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
	currentDateTime := time.Now().Format("20060102150405")
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

	outfilePath := filepath.Join(dir, currentDateTime+"_reactions.backup")
	outfile, err := os.OpenFile(outfilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()
	pageToken := []byte{}

	count := 0
	for {
		// response, next, err := hub.GetLinksByFid(uint64(fid), next, 10)
		// response, next, err := hub.GetCastsByFid(280, next, 10)
		response, nextPageToken, err := hub.GetReactionsByFid(uint64(fid), pageToken, PAGE_SIZE)
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
		fmt.Printf("... Saving %d reactions\t%s", count, nextStr)
		if len(nextPageToken) == 0 {
			break
		}
		pageToken = nextPageToken
	}
	fmt.Println(" Done.")

	outfilePath = filepath.Join(dir, currentDateTime+"_casts.backup")
	outfile, err = os.OpenFile(outfilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()
	pageToken = []byte{}

	count = 0
	for {
		// response, next, err := hub.GetLinksByFid(uint64(fid), next, 10)
		response, nextPageToken, err := hub.GetCastsByFid(uint64(fid), pageToken, PAGE_SIZE)
		// response, nextPageToken, err := hub.GetReactionsByFid(uint64(fid), pageToken, PAGE_SIZE)
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
		fmt.Printf("... Saving %d casts\t%s", count, nextStr)
		if len(nextPageToken) == 0 {
			break
		}
		pageToken = nextPageToken
	}
	fmt.Println(" Done.")

	outfilePath = filepath.Join(dir, currentDateTime+"_links.backup")
	outfile, err = os.OpenFile(outfilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()
	pageToken = []byte{}

	count = 0
	for {
		response, nextPageToken, err := hub.GetLinksByFid(uint64(fid), pageToken, PAGE_SIZE)
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
		fmt.Printf("... Saving %d links\t%s", count, nextStr)
		if len(nextPageToken) == 0 {
			break
		}
		pageToken = nextPageToken
	}
	fmt.Println(" Done.")
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("dir", "d", ".", "Path to export dir")
	downloadCmd.Flags().String("hub", "hoyt.farcaster.xyz:2283", "Farcaster hub to use")
	downloadCmd.Flags().Bool("ssl", true, "Use SSL")
}
