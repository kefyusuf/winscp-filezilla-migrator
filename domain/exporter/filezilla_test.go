package exporter

import (
	"os"
	"testing"

	"github.com/muety/winscp2filezilla/domain/models"
)

func TestExportToFileZilla(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sites*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	sessions := []models.Session{
		{Name: "Server1", HostName: "ftp.example.com", UserName: "admin", Password: "pass123", FSProtocol: "5", PortNumber: "21"},
		{Name: "Server2", HostName: "sftp.example.com", UserName: "user", Password: "pass456", FSProtocol: "2", PortNumber: "22"},
	}

	err = ExportToFileZilla(sessions, tmpFile.Name())
	if err != nil {
		t.Fatalf("ExportToFileZilla() error = %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("exported file is empty")
	}

	if len(content) < 100 {
		t.Errorf("exported file seems too short: %d bytes", len(content))
	}
}

func TestMapProtocol(t *testing.T) {
	tests := []struct {
		fsProtocol string
		want       string
	}{
		{"2", "1"},   // SFTP
		{"5", "0"},   // FTP
		{"", "0"},    // Default FTP
		{"1", "0"},   // Unknown -> FTP
	}

	for _, tt := range tests {
		got := mapProtocol(tt.fsProtocol)
		if got != tt.want {
			t.Errorf("mapProtocol(%q) = %q, want %q", tt.fsProtocol, got, tt.want)
		}
	}
}

func TestMapPort(t *testing.T) {
	tests := []struct {
		fsProtocol  string
		portNumber string
		want       string
	}{
		{"2", "", "22"},      // SFTP default
		{"5", "", "21"},      // FTP default
		{"2", "2222", "2222"}, // custom port
		{"", "", "21"},       // default FTP
	}

	for _, tt := range tests {
		got := mapPort(tt.fsProtocol, tt.portNumber)
		if got != tt.want {
			t.Errorf("mapPort(%q, %q) = %q, want %q", tt.fsProtocol, tt.portNumber, got, tt.want)
		}
	}
}