package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kefyusuf/winscp-filezilla-migrator/domain/exporter"
	"github.com/kefyusuf/winscp-filezilla-migrator/domain/parser"
)

//go:embed static/*
var staticEmbed embed.FS

type Server struct {
	Name       string `json:"name"`
	HostName   string `json:"host"`
	UserName   string `json:"user"`
	PortNumber string `json:"port"`
	FSProtocol string `json:"protocol"`
	RemoteDir  string `json:"remoteDir"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func apiPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func router(fileServer http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/ping":
			apiPing(w, r)
		case "/api/upload":
			handleUpload(w, r)
		case "/api/migrate":
			handleMigrate(w, r)
		default:
			fileServer.ServeHTTP(w, r)
		}
	})
}

func main() {
	staticFS, _ := fs.Sub(staticEmbed, "static")
	fileServer := http.FileServer(http.FS(staticFS))

	http.ListenAndServe(":9090", router(fileServer))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, "Method not allowed")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendError(w, "No file uploaded")
		return
	}
	defer file.Close()

	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "winscp_"+header.Filename)
	out, err := os.Create(tmpPath)
	if err != nil {
		sendError(w, "Cannot create temp file")
		return
	}
	defer os.Remove(tmpPath)

	if _, err := io.Copy(out, file); err != nil {
		sendError(w, "Cannot save file")
		return
	}
	out.Close()

	sessions, err := parser.ParseWinSCPIni(tmpPath)
	if err != nil {
		sendError(w, fmt.Sprintf("Parse error: %v", err))
		return
	}

	servers := make([]Server, len(sessions))
	for i, s := range sessions {
		protocol := "FTP"
		if s.FSProtocol == "2" {
			protocol = "SFTP"
		}
		port := s.PortNumber
		if port == "" {
			if s.FSProtocol == "2" {
				port = "22"
			} else {
				port = "21"
			}
		}

		servers[i] = Server{
			Name:       s.Name,
			HostName:   s.HostName,
			UserName:   s.UserName,
			PortNumber: port,
			FSProtocol: protocol,
			RemoteDir:  s.RemoteDir,
		}
	}

	sendSuccess(w, servers)
}

func handleMigrate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, "Method not allowed")
		return
	}

	var req struct {
		Servers []Server `json:"servers"`
		Output  string   `json:"output"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request")
		return
	}

	sessions := make([]interface{}, len(req.Servers))
	for i, s := range req.Servers {
		fsProtocol := "5"
		if s.FSProtocol == "SFTP" {
			fsProtocol = "2"
		}
		sessions[i] = map[string]string{
			"Name":       s.Name,
			"HostName":   s.HostName,
			"UserName":   s.UserName,
			"PortNumber": s.PortNumber,
			"FSProtocol": fsProtocol,
			"Password":   "",
			"RemoteDir":  s.RemoteDir,
			"LocalDir":   "",
		}
	}

	err := exporter.ExportToFileZillaRaw(sessions, req.Output)
	if err != nil {
		sendError(w, fmt.Sprintf("Export error: %v", err))
		return
	}

	sendSuccess(w, map[string]string{"message": "Migration completed!"})
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{Success: true, Data: data})
}

func sendError(w http.ResponseWriter, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(ApiResponse{Success: false, Error: err})
}