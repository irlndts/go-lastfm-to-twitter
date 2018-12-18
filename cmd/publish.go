package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish your chart to twitter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("publish called")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
