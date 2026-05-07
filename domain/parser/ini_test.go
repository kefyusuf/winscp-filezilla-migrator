package parser

import (
	"os"
	"testing"

	"github.com/muety/winscp2filezilla/domain/models"
)

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