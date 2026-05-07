package app

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/muety/winscp2filezilla/domain/exporter"
	"github.com/muety/winscp2filezilla/domain/parser"
)

var (
	inputPath  string
	outputPath string
)

func Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:     "winscp2filezilla",
		Short:   "Migrate WinSCP servers to FileZilla",
		Long:    "A tool that migrates saved FTP/SFTP server configurations from WinSCP to FileZilla.",
		Version: "2.0.0",
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate WinSCP.ini to FileZilla sites.xml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate()
		},
	}

	migrateCmd.Flags().StringVarP(&inputPath, "in", "i", getDefaultIniPath(), "Path to WinSCP.ini")
	migrateCmd.Flags().StringVarP(&outputPath, "out", "o", "sites.xml", "Output path for sites.xml")
	rootCmd.AddCommand(migrateCmd)

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

func runMigrate() error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}

	fmt.Printf("Reading WinSCP.ini from: %s\n", inputPath)

	sessions, err := parser.ParseWinSCPIni(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse INI: %w", err)
	}

	fmt.Printf("Found %d sessions\n", len(sessions))

	if outputPath == "" {
		outputPath = "sites.xml"
	}

	fmt.Printf("Writing FileZilla XML to: %s\n", outputPath)

	err = exporter.ExportToFileZilla(sessions, outputPath)
	if err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Println("Migration completed successfully!")
	return nil
}

func getDefaultIniPath() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir + "\\AppData\\Roaming\\WinSCP.ini"
}