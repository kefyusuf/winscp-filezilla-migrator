package models

type Session struct {
	Name       string
	HostName   string
	UserName   string
	PortNumber string
	FSProtocol string
	Password   string
	RemoteDir  string
	LocalDir   string
}