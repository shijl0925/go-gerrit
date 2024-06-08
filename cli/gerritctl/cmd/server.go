package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// serverInfo Command
var serverInfo = &cobra.Command{
	Use:   "version",
	Short: "get server version",
	Run: func(cmd *cobra.Command, args []string) {
		version, _, err := gerritMod.Instance.Config.GetVersion(gerritMod.Context)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("✅ Connected with: %s\n", gerritMod.Username)
		fmt.Printf("✅ Server: %s\n", gerritMod.Url)
		fmt.Printf("✅ Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(serverInfo)
}
