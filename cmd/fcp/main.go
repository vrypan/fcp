package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vrypan/fcp/utils"
)

var FCP_VERSION string

func main() {
	Execute()
}

// cliCmd represents the cli command
var rootCmd = &cobra.Command{
	Use:   "fcp <source> <destination>",
	Short: "Copy user data from/to farcaster and the local filesystem",
	Long: `Copy user data from/to farcaster and the local filesystem.

Examples:
  Copy fname's casts, reactions and links from your hub to
  a local file:
  fcp fc://hubble.local:2283/fname localfile.data

  Upload the local file to your hub:
  fcp localfile.data fc://hubble.local:2283
  (Will probably fail because the messages already exist)

  Use fc+ssl if the hub is using SSL:
  fcp fc+ssl://hubble.local:2283/fname localfile.data

  When uploading data, you can use a new signer to resign all messages:
  fcp local.data fc+ssl://hubble.local:2283/fname --app-key=0x...

  Use "-" for stdin/stdout depending on the operation.
`,
	Run: fcpCmdMain,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func fcpCmdMain(cmd *cobra.Command, args []string) {
	version, _ := cmd.Flags().GetBool("version")
	if version {
		fmt.Printf("version %s\n", FCP_VERSION)
		os.Exit(0)
	}
	if len(args) != 2 {
		cmd.Help()
		return
	}
	hubAddress, useSsl, username, err := utils.ParseUrl(args[0])
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}

	if hubAddress == "" { // source is local file
		hubAddress, useSsl, username, err = utils.ParseUrl(args[1])
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			return
		}
		if hubAddress == "" {
			fmt.Printf("No hub address specified\n")
			return
		}
		opts := map[string]any{}
		opts["ssl"] = useSsl
		opts["pageSize"], _ = cmd.Flags().GetUint32("page-size")
		opts["signer"], _ = cmd.Flags().GetString("app-key")
		opts["reactions"], _ = cmd.Flags().GetBool("reactions")
		opts["links"], _ = cmd.Flags().GetBool("links")
		opts["casts"], _ = cmd.Flags().GetBool("casts")
		utils.Upload(hubAddress, args[0], opts)
	} else {
		if h, _, _, err := utils.ParseUrl(args[1]); h != "" || err != nil {
			fmt.Println("Invalid destination")
			return
		}
		opts := map[string]any{}
		opts["ssl"] = useSsl
		opts["pageSize"], _ = cmd.Flags().GetUint32("page-size")
		opts["signer"], _ = cmd.Flags().GetString("app-key")
		opts["reactions"], _ = cmd.Flags().GetBool("reactions")
		opts["links"], _ = cmd.Flags().GetBool("links")
		opts["casts"], _ = cmd.Flags().GetBool("casts")
		utils.Download(hubAddress, username[1:], args[1], opts)
	}

}
func init() {
	//rootCmd.AddCommand(fcpCmd)
	rootCmd.Flags().Uint32("page-size", 100, "Hub request page size")
	rootCmd.Flags().StringP("app-key", "k", "", "App key (signer)")
	rootCmd.Flags().BoolP("stats", "s", false, "Display stats")
	rootCmd.Flags().Bool("casts", true, "Read/write casts")
	rootCmd.Flags().Bool("reactions", true, "Read/write reactions (likes, recasts)")
	rootCmd.Flags().Bool("links", true, "Read/write links (follows)")
	rootCmd.Flags().BoolP("version", "v", false, "Display fcp version")
}
