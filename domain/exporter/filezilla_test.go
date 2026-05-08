package exporter

import (
	"os"
	"strings"
	"testing"

	"github.com/kefyusuf/winscp-filezilla-migrator/domain/models"
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

	if strings.Contains(content, "<Folder>Imported") {
		t.Error("root folder should not be 'Imported' anymore")
	}
	if !strings.Contains(content, "<Folder>WinSCP-") {
		t.Error("expected root folder to start with 'WinSCP-'")
	}
}

func TestExportToFileZilla_FolderStructure(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sites*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	sessions := []models.Session{
		{Name: "freelancer\\New site", HostName: "new.freelancer.com", UserName: "user1", Password: "p1", FSProtocol: "5"},
		{Name: "sm724-live\\adatepe.com", HostName: "adatepe.com", UserName: "user2", Password: "p2", FSProtocol: "2"},
		{Name: "sm724\\erkatest.sm724pro.com", HostName: "erkatest.sm724pro.com", UserName: "user3", Password: "p3", FSProtocol: "2"},
		{Name: "FlatServer", HostName: "flat.example.com", UserName: "user4", Password: "p4", FSProtocol: "5"},
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
	t.Logf("Output:\n%s", content)

	if !strings.Contains(content, "<Folder>freelancer") {
		t.Error("expected folder 'freelancer' in output")
	}
	if !strings.Contains(content, "<Name>New site</Name>") {
		t.Error("expected server 'New site' inside freelancer folder")
	}
	if !strings.Contains(content, "<Folder>sm724-live") {
		t.Error("expected folder 'sm724-live' in output")
	}
	if !strings.Contains(content, "<Name>adatepe.com</Name>") {
		t.Error("expected server 'adatepe.com' inside sm724-live folder")
	}
	if !strings.Contains(content, "<Folder>sm724") {
		t.Error("expected folder 'sm724' in output")
	}
	if !strings.Contains(content, "<Name>erkatest.sm724pro.com</Name>") {
		t.Error("expected server 'erkatest.sm724pro.com' inside sm724 folder")
	}
	if !strings.Contains(content, "<Name>FlatServer</Name>") {
		t.Error("expected flat server 'FlatServer'")
	}
}

func TestExportToFileZilla_NestedFolders(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sites*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	sessions := []models.Session{
		{Name: "Parent\\Child\\ServerName", HostName: "server.example.com", UserName: "user", Password: "p", FSProtocol: "5"},
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
	t.Logf("Nested folder output:\n%s", content)

	if strings.Contains(content, "<Name>Parent\\Child\\ServerName</Name>") {
		t.Error("full path should NOT be used as server name")
	}
	if !strings.Contains(content, "<Name>ServerName</Name>") {
		t.Error("expected server name to be only the last segment")
	}
	if !strings.Contains(content, "<Folder>Parent") {
		t.Error("expected parent folder 'Parent'")
	}
	if !strings.Contains(content, "<Folder>Child") {
		t.Error("expected child folder 'Child'")
	}
	if strings.Contains(content, "<Folder>Parent\\Child") {
		t.Error("Parent and Child should be separate nested folders, not one combined folder")
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