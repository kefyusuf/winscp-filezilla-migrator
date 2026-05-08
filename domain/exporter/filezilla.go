package exporter

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/models"
)

type FolderNode struct {
	Name     string
	Servers  []models.Session
	Children []*FolderNode
}

func ExportToFileZilla(sessions []models.Session, outputPath string) error {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	rootNode := doc.CreateElement("FileZilla3")
	serversNode := rootNode.CreateElement("Servers")
	rootFolder := serversNode.CreateElement("Folder")
	rootFolder.SetText("WinSCP-" + time.Now().Format("2006-01-02_1504"))

	rootTree := buildFolderTree(sessions)

	for _, folder := range rootTree.Children {
		addFolderToXML(rootFolder, folder)
	}

	for _, server := range rootTree.Servers {
		addServerToXML(rootFolder, server)
	}

	doc.Indent(2)
	err := doc.WriteToFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to write XML file: %w", err)
	}

	return nil
}

func buildFolderTree(sessions []models.Session) *FolderNode {
	root := &FolderNode{}

	for _, session := range sessions {
		parts := strings.Split(session.Name, "\\")

		if len(parts) <= 1 {
			root.Servers = append(root.Servers, session)
			continue
		}

		// Navigate or create folder hierarchy
		current := root
		for i := 0; i < len(parts)-1; i++ {
			name := parts[i]
			found := false
			for _, child := range current.Children {
				if child.Name == name {
					current = child
					found = true
					break
				}
			}
			if !found {
				node := &FolderNode{Name: name}
				current.Children = append(current.Children, node)
				current = node
			}
		}

		current.Servers = append(current.Servers, session)
	}

	return root
}

func addFolderToXML(parent *etree.Element, folder *FolderNode) {
	var folderElement *etree.Element

	for _, child := range parent.FindElements("Folder") {
		if child.Text() == folder.Name {
			folderElement = child
			break
		}
	}

	if folderElement == nil {
		folderElement = parent.CreateElement("Folder")
		folderElement.SetText(folder.Name)
	}

	for _, server := range folder.Servers {
		addServerToXML(folderElement, server)
	}

	for _, child := range folder.Children {
		addFolderToXML(folderElement, child)
	}
}

func addServerToXML(parent *etree.Element, session models.Session) {
	serverNode := parent.CreateElement("Server")

	nameNode := serverNode.CreateElement("Name")
	nameParts := strings.Split(session.Name, "\\")
	nameNode.SetText(nameParts[len(nameParts)-1])

	hostNode := serverNode.CreateElement("Host")
	hostNode.SetText(session.HostName)

	protocolNode := serverNode.CreateElement("Protocol")
	protocolNode.SetText(mapProtocol(session.FSProtocol))

	userNode := serverNode.CreateElement("User")
	userNode.SetText(session.UserName)

	passNode := serverNode.CreateElement("Pass")
	passNode.SetText(base64.StdEncoding.EncodeToString([]byte(session.Password)))
	passNode.CreateAttr("encoding", "base64")

	portNode := serverNode.CreateElement("Port")
	portNode.SetText(mapPort(session.FSProtocol, session.PortNumber))

	localDirNode := serverNode.CreateElement("LocalDir")
	localDirNode.SetText(session.LocalDir)

	remoteDirNode := serverNode.CreateElement("RemoteDir")
	remoteDirNode.SetText(session.RemoteDir)

	typeNode := serverNode.CreateElement("Type")
	typeNode.SetText("0")

	logonTypeNode := serverNode.CreateElement("Logontype")
	logonTypeNode.SetText("1")

	tzOffsetNode := serverNode.CreateElement("TimezoneOffset")
	tzOffsetNode.SetText("0")

	pasvModeNode := serverNode.CreateElement("PasvMode")
	pasvModeNode.SetText("MODE_DEFAULT")

	maxConnNode := serverNode.CreateElement("MaximumMultipleConnections")
	maxConnNode.SetText("0")

	encodingNode := serverNode.CreateElement("EncodingType")
	encodingNode.SetText("Auto")

	bypassNode := serverNode.CreateElement("BypassProxy")
	bypassNode.SetText("0")

	syncNode := serverNode.CreateElement("SyncBrowsing")
	syncNode.SetText("0")

	_ = serverNode.CreateElement("Comments")
}

func mapProtocol(fsProtocol string) string {
	if fsProtocol == "2" {
		return "1"
	}
	return "0"
}

func ExportToFileZillaRaw(sessions []interface{}, outputPath string) error {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	rootNode := doc.CreateElement("FileZilla3")
	serversNode := rootNode.CreateElement("Servers")
	rootFolder := serversNode.CreateElement("Folder")
	rootFolder.SetText("WinSCP-" + time.Now().Format("2006-01-02_1504"))

	for _, s := range sessions {
		m := s.(map[string]string)
		serverNode := rootFolder.CreateElement("Server")

		nameNode := serverNode.CreateElement("Name")
		nameNode.SetText(m["Name"])

		hostNode := serverNode.CreateElement("Host")
		hostNode.SetText(m["HostName"])

		protocolNode := serverNode.CreateElement("Protocol")
		if m["FSProtocol"] == "2" {
			protocolNode.SetText("1")
		} else {
			protocolNode.SetText("0")
		}

		userNode := serverNode.CreateElement("User")
		userNode.SetText(m["UserName"])

		passNode := serverNode.CreateElement("Pass")
		passNode.SetText(base64.StdEncoding.EncodeToString([]byte(m["Password"])))
		passNode.CreateAttr("encoding", "base64")

		portNode := serverNode.CreateElement("Port")
		portNode.SetText(m["PortNumber"])

		localDirNode := serverNode.CreateElement("LocalDir")
		_ = localDirNode
		remoteDirNode := serverNode.CreateElement("RemoteDir")
		remoteDirNode.SetText(m["RemoteDir"])

		typeNode := serverNode.CreateElement("Type")
		typeNode.SetText("0")
		logonNode := serverNode.CreateElement("Logontype")
		logonNode.SetText("1")
		tzNode := serverNode.CreateElement("TimezoneOffset")
		tzNode.SetText("0")
		pasvNode := serverNode.CreateElement("PasvMode")
		pasvNode.SetText("MODE_DEFAULT")
		maxNode := serverNode.CreateElement("MaximumMultipleConnections")
		maxNode.SetText("0")
		encNode := serverNode.CreateElement("EncodingType")
		encNode.SetText("Auto")
		bypassNode := serverNode.CreateElement("BypassProxy")
		bypassNode.SetText("0")
		syncNode := serverNode.CreateElement("SyncBrowsing")
		syncNode.SetText("0")
	}

	doc.Indent(2)
	return doc.WriteToFile(outputPath)
}

func mapPort(fsProtocol string, portNumber string) string {
	if portNumber != "" {
		return portNumber
	}
	if fsProtocol == "2" {
		return "22"
	}
	return "21"
}