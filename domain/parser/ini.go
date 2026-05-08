package parser

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"gopkg.in/ini.v1"

	"github.com/kefyusuf/winscp-filezilla-migrator/domain/models"
)

func ParseWinSCPIni(filepath string) ([]models.Session, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read INI file: %w", err)
	}

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	if !utf8.Valid(data) {
		decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), charmap.Windows1254.NewDecoder()))
		if err != nil {
			return nil, fmt.Errorf("failed to decode INI from Windows-1254: %w", err)
		}
		data = decoded
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %w", err)
	}

	sessions := make([]models.Session, 0)

	for _, s := range cfg.Sections() {
		sectionName := strings.ReplaceAll(s.Name(), "/", "\\")
		if !strings.HasPrefix(sectionName, "Sessions\\") || !s.HasKey("HostName") {
			continue
		}

		session := models.Session{}
		session.Name = cleanPathName(urlDecode(strings.TrimPrefix(sectionName, "Sessions\\")))
		session.HostName = urlDecode(s.Key("HostName").Value())

		if s.HasKey("UserName") {
			session.UserName = urlDecode(s.Key("UserName").Value())
		}

		encryptedPassword := s.Key("Password").Value()
		session.Password = Decrypt(session.HostName, session.UserName, encryptedPassword)

		if s.HasKey("FSProtocol") {
			session.FSProtocol = s.Key("FSProtocol").Value()
		}

		if s.HasKey("PortNumber") {
			session.PortNumber = s.Key("PortNumber").Value()
		}

		if s.HasKey("RemoteDir") {
			session.RemoteDir = urlDecode(s.Key("RemoteDir").Value())
		}

		if s.HasKey("LocalDir") {
			session.LocalDir = urlDecode(s.Key("LocalDir").Value())
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func cleanPathName(name string) string {
	name = strings.ReplaceAll(name, "/", "\\")
	name = strings.TrimSuffix(name, "\\")
	parts := strings.Split(name, "\\")
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		p = stripSchemePrefix(p)
		if p == "" {
			continue
		}
		cleaned = append(cleaned, p)
	}
	return strings.Join(cleaned, "\\")
}

var schemes = []string{"http:", "https:", "ftp:", "ftps:", "sftp:"}

func stripSchemePrefix(s string) string {
	for _, prefix := range schemes {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimPrefix(s, prefix)
		}
	}
	return s
}

func urlDecode(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return strings.TrimLeft(decoded, "\uFEFF")
}