package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestRouter() http.Handler {
	staticFS, _ := fs.Sub(staticEmbed, "static")
	fileServer := http.FileServer(http.FS(staticFS))
	return router(fileServer)
}

func TestPing(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("status = %q, want %q", resp["status"], "ok")
	}
}

func TestUpload_NoFile(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/upload", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp ApiResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error == "" {
		t.Error("expected error message")
	}
}

func TestUpload_WrongMethod(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/upload", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpload_ValidFile(t *testing.T) {
	iniContent := `[Sessions\Production\Web]
HostName=web.example.com
UserName=admin
Password=A3FF0000
FSProtocol=2
PortNumber=22
RemoteDir=/var/www

[Sessions\Dev\Test]
HostName=dev.example.com
UserName=test
Password=A3FF0001
FSProtocol=5
`

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "test.ini")
	fw.Write([]byte(iniContent))
	w.Close()

	r := newTestRouter()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp ApiResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("success = false, want true; error = %q", resp.Error)
	}

	servers, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("data is not an array")
	}
	if len(servers) != 2 {
		t.Fatalf("got %d servers, want 2", len(servers))
	}

	s0 := servers[0].(map[string]interface{})
	if s0["host"] != "web.example.com" {
		t.Errorf("server[0].host = %q", s0["host"])
	}
	if s0["protocol"] != "SFTP" {
		t.Errorf("server[0].protocol = %q", s0["protocol"])
	}

	s1 := servers[1].(map[string]interface{})
	if s1["protocol"] != "FTP" {
		t.Errorf("server[1].protocol = %q", s1["protocol"])
	}
}

func TestMigrate_EmptyBody(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/migrate", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestMigrate_WrongMethod(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/migrate", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestMigrate_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "sites.xml")

	payload := map[string]interface{}{
		"servers": []map[string]string{
			{"name": "Test", "host": "example.com", "user": "admin", "port": "22", "protocol": "SFTP", "remoteDir": "/var/www"},
		},
		"output": outputPath,
	}
	body, _ := json.Marshal(payload)

	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/migrate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp ApiResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("success = false, want true; error = %q", resp.Error)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestStatic_ServesIndex(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "WinSCP to FileZilla Migrator") {
		t.Error("response doesn't contain expected title")
	}
}

func TestStatic_ServesIndexAtRoot(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/index.html", nil)
	r.ServeHTTP(w, req)

	// File server redirects /index.html to /
	if w.Code != http.StatusMovedPermanently && w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 or 301", w.Code)
	}
}

func TestRouter_NotFound(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpload_InvalidIni(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "bad.ini")
	fw.Write([]byte("not an ini file content"))
	w.Close()

	r := newTestRouter()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestMigrate_InvalidJson(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/migrate",
		io.NopCloser(strings.NewReader("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
