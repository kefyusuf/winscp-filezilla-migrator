package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:   "winscp2filezilla",
		Short: "Migrate WinSCP servers to FileZilla",
		Long:  "A tool that migrates saved FTP/SFTP server configurations from WinSCP to FileZilla.",
		Version: "2.0.0",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("winscp2filezilla v2.0.0")
		},
	})

	rootCmd.SetArgs(args[1:])
	return rootCmd.Execute()
}