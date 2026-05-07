package exporter

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/muety/winscp2filezilla/domain/models"
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
	rootFolder.SetText("Imported")

	rootLevelServers := []models.Session{}
	folderMap := buildFolderTree(sessions)

	for _, folder := range folderMap {
		addFolderToXML(rootFolder, folder)
	}

	for _, session := range sessions {
		hasFolder := false
		for _, folder := range folderMap {
			for _, s := range folder.Servers {
				if s.Name == session.Name {
					hasFolder = true
					break
				}
			}
			if hasFolder {
				break
			}
		}
		if !hasFolder {
			rootLevelServers = append(rootLevelServers, session)
		}
	}

	for _, server := range rootLevelServers {
		addServerToXML(rootFolder, server)
	}

	doc.Indent(2)
	err := doc.WriteToFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to write XML file: %w", err)
	}

	return nil
}

func buildFolderTree(sessions []models.Session) []*FolderNode {
	folderMap := make(map[string]*FolderNode)

	for _, session := range sessions {
		parts := strings.Split(session.Name, "\\")

		if len(parts) <= 1 {
			continue
		}

		folderPath := strings.Join(parts[:len(parts)-1], "\\")

		if _, exists := folderMap[folderPath]; !exists {
			folderMap[folderPath] = &FolderNode{Name: folderPath}
		}
		folderMap[folderPath].Servers = append(folderMap[folderPath].Servers, session)
	}

	result := make([]*FolderNode, 0)
	for _, f := range folderMap {
		result = append(result, f)
	}

	return result
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
	nameParts := strings.Split(session.Name, "/")
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

func mapPort(fsProtocol string, portNumber string) string {
	if portNumber != "" {
		return portNumber
	}
	if fsProtocol == "2" {
		return "22"
	}
	return "21"
}