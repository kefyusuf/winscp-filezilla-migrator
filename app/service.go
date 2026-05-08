package app

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/exporter"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/parser"
)

var (
	inputPath  string
	outputPath string
)

func Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:     "winscp-filezilla-migrator",
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

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List parsed sessions from WinSCP.ini",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd)
		},
	}

	listCmd.Flags().StringVarP(&inputPath, "in", "i", getDefaultIniPath(), "Path to WinSCP.ini")
	rootCmd.AddCommand(listCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("winscp-filezilla-migrator v1.0.0")
		},
	})

	rootCmd.SetArgs(args[1:])
	return rootCmd.Execute()
}

func runList(cmd *cobra.Command) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}

	fmt.Printf("Reading WinSCP.ini from: %s\n\n", inputPath)

	sessions, err := parser.ParseWinSCPIni(inputPath)
	if err != nil {
		return fmt.Errorf("failed to parse INI: %w", err)
	}

	fmt.Printf("Found %d sessions:\n", len(sessions))
	for i, s := range sessions {
		fmt.Printf("  [%d] Name=%q  Host=%s  User=%s  Port=%s  Proto=%s\n",
			i, s.Name, s.HostName, s.UserName, s.PortNumber, s.FSProtocol)
	}

	return nil
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