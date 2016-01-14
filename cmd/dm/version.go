package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of this tool",
}

func init() {
	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %v\n", version)
	}
}
