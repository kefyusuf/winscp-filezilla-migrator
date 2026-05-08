package parser

import (
	"os"
	"testing"

	"github.com/kefyusuf/winscp-filezilla-migrator/domain/models"
)

func TestStripSchemePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http:erkatest.sm724pro.com", "erkatest.sm724pro.com"},
		{"https:secure.example.com", "secure.example.com"},
		{"ftp:ftp.example.com", "ftp.example.com"},
		{"ftps:ftps.example.com", "ftps.example.com"},
		{"sftp:sftp.example.com", "sftp.example.com"},
		{"http:", ""},
		{"https:", ""},
		{"NormalServer", "NormalServer"},
		{"folder\\server", "folder\\server"},
		{"", ""},
	}
	for _, tt := range tests {
		got := stripSchemePrefix(tt.input)
		if got != tt.expected {
			t.Errorf("stripSchemePrefix(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestCleanPathName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sm724\\http:erkatest.sm724pro.com", "sm724\\erkatest.sm724pro.com"},
		{"sm724\\http:", "sm724"},
		{"freelancer\\New site", "freelancer\\New site"},
		{"sm724-live\\adatepe.com", "sm724-live\\adatepe.com"},
		{"sm724\\http:\\erkatest.sm724pro.com", "sm724\\erkatest.sm724pro.com"},
		{"Production\\WebServer", "Production\\WebServer"},
		{"SimpleServer", "SimpleServer"},
		{"", ""},
		{"sm724/ase", "sm724\\ase"},
		{"sm724-live/adatepe.com", "sm724-live\\adatepe.com"},
		{"freelancer/New site", "freelancer\\New site"},
		{"sm724/http:erkatest.sm724pro.com", "sm724\\erkatest.sm724pro.com"},
		{"sm724/http:\\erkatest.sm724pro.com", "sm724\\erkatest.sm724pro.com"},
	}
	for _, tt := range tests {
		got := cleanPathName(tt.input)
		if got != tt.expected {
			t.Errorf("cleanPathName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestParseWinSCPIni(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions\Production\Server1]
HostName=web.example.com
UserName=admin
Password=A3FF0000
FSProtocol=2
PortNumber=22
RemoteDir=/var/www

[Sessions\Dev\Server2]
HostName=dev.example.com
UserName=test
Password=A3FF0001
FSProtocol=5
`

	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}

	expected := []models.Session{
		{Name: "Production\\Server1", HostName: "web.example.com", UserName: "admin", FSProtocol: "2", PortNumber: "22", RemoteDir: "/var/www"},
		{Name: "Dev\\Server2", HostName: "dev.example.com", UserName: "test", FSProtocol: "5"},
	}

	for i, exp := range expected {
		if i >= len(sessions) {
			break
		}
		if sessions[i].Name != exp.Name {
			t.Errorf("session[%d].Name = %q, want %q", i, sessions[i].Name, exp.Name)
		}
		if sessions[i].HostName != exp.HostName {
			t.Errorf("session[%d].HostName = %q, want %q", i, sessions[i].HostName, exp.HostName)
		}
	}
}

func TestParseWinSCPIni_SchemePrefix(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions\sm724\http:erkatest.sm724pro.com]
HostName=erkatest.sm724pro.com
UserName=admin
Password=A3FF0000
FSProtocol=2

[Sessions\sm724\https:secure.example.com]
HostName=secure.example.com
UserName=admin
Password=A3FF0001
FSProtocol=5
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	if sessions[0].Name != "sm724\\erkatest.sm724pro.com" {
		t.Errorf("session[0].Name = %q, want %q", sessions[0].Name, "sm724\\erkatest.sm724pro.com")
	}
	if sessions[1].Name != "sm724\\secure.example.com" {
		t.Errorf("session[1].Name = %q, want %q", sessions[1].Name, "sm724\\secure.example.com")
	}
}

func TestParseWinSCPIni_ForwardSlashSeparator(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions\freelancer/New site]
HostName=new.freelancer.com
UserName=user
Password=A3FF0000
FSProtocol=5

[Sessions\sm724-live/adatepe.com]
HostName=adatepe.com
UserName=admin
Password=A3FF0001
FSProtocol=2
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	if sessions[0].Name != "freelancer\\New site" {
		t.Errorf("session[0].Name = %q, want %q", sessions[0].Name, "freelancer\\New site")
	}
	if sessions[1].Name != "sm724-live\\adatepe.com" {
		t.Errorf("session[1].Name = %q, want %q", sessions[1].Name, "sm724-live\\adatepe.com")
	}
}

func TestParseWinSCPIni_AllForwardSlash(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions/freelancer/New site]
HostName=new.freelancer.com
UserName=user
Password=A3FF0000
FSProtocol=5

[Sessions/sm724-live/adatepe.com]
HostName=adatepe.com
UserName=admin
Password=A3FF0001
FSProtocol=2
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	if sessions[0].Name != "freelancer\\New site" {
		t.Errorf("session[0].Name = %q, want %q", sessions[0].Name, "freelancer\\New site")
	}
	if sessions[1].Name != "sm724-live\\adatepe.com" {
		t.Errorf("session[1].Name = %q, want %q", sessions[1].Name, "sm724-live\\adatepe.com")
	}
}

func TestParseWinSCPIni_SchemePrefixForwardSlash(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions\sm724/http:\erkatest.sm724pro.com]
HostName=erkatest.sm724pro.com
UserName=admin
Password=A3FF0000
FSProtocol=2
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	if sessions[0].Name != "sm724\\erkatest.sm724pro.com" {
		t.Errorf("session.Name = %q, want %q", sessions[0].Name, "sm724\\erkatest.sm724pro.com")
	}
}

func TestParseWinSCPIni_SchemeOnlySegment(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "winscp*.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `[Sessions\sm724\http:\erkatest.sm724pro.com]
HostName=erkatest.sm724pro.com
UserName=admin
Password=A3FF0000
FSProtocol=2
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	sessions, err := ParseWinSCPIni(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseWinSCPIni() error = %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	if sessions[0].Name != "sm724\\erkatest.sm724pro.com" {
		t.Errorf("session.Name = %q, want %q", sessions[0].Name, "sm724\\erkatest.sm724pro.com")
	}
}