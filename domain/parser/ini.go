package parser

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/muety/winscp2filezilla/domain/models"
)

func ParseWinSCPIni(filepath string) ([]models.Session, error) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %w", err)
	}

	sessions := make([]models.Session, 0)

	for _, s := range cfg.Sections() {
		if !strings.HasPrefix(s.Name(), "Sessions\\") || !s.HasKey("HostName") {
			continue
		}

		session := models.Session{}
		session.Name = strings.TrimPrefix(s.Name(), "Sessions\\")
		session.HostName = s.Key("HostName").Value()

		if s.HasKey("UserName") {
			session.UserName = s.Key("UserName").Value()
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
			session.RemoteDir = s.Key("RemoteDir").Value()
		}

		if s.HasKey("LocalDir") {
			session.LocalDir = s.Key("LocalDir").Value()
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}