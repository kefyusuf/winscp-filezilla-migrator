# WinSCP to FileZilla Migrator

A cross-platform desktop application that migrates saved FTP/SFTP server configurations from WinSCP to FileZilla. Select your WinSCP.ini file, preview your servers, and export them to FileZilla's format — all through a native GUI.

## Features

- **Native GUI** — Built with [Gio](https://gioui.org), runs on Windows, Linux, and macOS without a browser
- **CLI mode** — Headless migration for scripting and automation
- **Web server** — Optional web-based interface accessible from any browser
- **Password decryption** — Preserves the original WinSCP XOR-based password algorithm
- **Folder hierarchy** — WinSCP folder structure is preserved in the FileZilla export
- **Date-stamped export** — Each export creates a `WinSCP-YYYY-MM-DD_HHMM` folder for easy tracking
- **Turkish character support** — Automatically handles Windows-1254 encoded INI files
- **Data integrity** — Passwords are decrypted and re-encoded as base64 for FileZilla

## Downloads

Pre-built binaries are available on the [Releases](https://github.com/kefyusuf/winscp-filezilla-migrator/releases) page:

| Package | Contents |
|---------|----------|
| `winscp-filezilla-migrator_x.y.z_windows_amd64.zip` | GUI (double-click to run) |
| `winscp-filezilla-migrator-cli_x.y.z_windows_amd64.zip` | CLI executable |
| `winscp-filezilla-migrator-web_x.y.z_windows_amd64.zip` | Web server executable |
| `winscp-filezilla-migrator_x.y.z_linux_amd64.tar.gz` | Linux GUI |
| `winscp-filezilla-migrator_x.y.z_darwin_amd64.tar.gz` | macOS GUI |

## Usage

### GUI (recommended)

Run `winscp-filezilla-migrator` (or double-click `winscp-filezilla-migrator.exe` on Windows).

1. Click **Open WinSCP.ini** and select your WinSCP configuration file
2. Browse the parsed server list on the left panel
3. Click a server to view its details (host, port, protocol, credentials)
4. Click **Export to FileZilla XML** and choose a save location

### CLI

```sh
# List parsed sessions from WinSCP.ini
winscp-filezilla-migrator-cli list --in "C:\Users\You\AppData\Roaming\WinSCP.ini"

# Migrate from WinSCP.ini to FileZilla sites.xml
winscp-filezilla-migrator-cli migrate --in "C:\Users\You\AppData\Roaming\WinSCP.ini" --out sites.xml

# Use default WinSCP.ini path
winscp-filezilla-migrator-cli migrate --out sites.xml

# Quick start with defaults (reads %APPDATA%\WinSCP.ini, writes sites.xml)
winscp-filezilla-migrator-cli migrate
```

### Web Server

```sh
winscp-filezilla-migrator-web
# Server starts at http://localhost:9090
```

Open a browser to `http://localhost:9090` for a drag-and-drop web interface.

## Build from Source

Requires [Go 1.21+](https://go.dev/dl/).

```sh
# Clone
git clone https://github.com/kefyusuf/winscp-filezilla-migrator.git
cd winscp-filezilla-migrator

# GUI (Windows — hides console window)
go build -ldflags="-s -w -H=windowsgui" -o winscp-filezilla-migrator.exe .

# GUI (Linux / macOS)
go build -ldflags="-s -w" -o winscp-filezilla-migrator .

# CLI
go build -ldflags="-s -w" -o winscp-filezilla-migrator-cli ./cmd/cli/

# Web server
go build -ldflags="-s -w" -o winscp-filezilla-migrator-web ./web/

# Run all tests
go test ./...
```

## Project Structure

```
.
├── main.go              # GUI entry point (Gio)
├── cmd/
│   └── cli/main.go      # CLI entry point (Cobra)
├── web/
│   ├── main.go          # Web server entry point
│   ├── main_test.go     # HTTP integration tests
│   └── static/          # Embedded HTML/JS frontend
├── app/
│   └── service.go       # CLI service with Cobra commands
├── domain/
│   ├── models/          # Session data model
│   ├── parser/          # WinSCP.ini parser + password decryption
│   └── exporter/        # FileZilla XML exporter
└── testdata/            # Sample WinSCP.ini files for testing
```

## How It Works

1. **Parse** — Reads WinSCP.ini using `gopkg.in/ini.v1`, extracts `[Sessions\*]` sections
2. **Decrypt** — Applies WinSCP's XOR-based algorithm (PW_MAGIC=0xA3, PW_FLAG=0xFF) to recover passwords
3. **Export** — Generates a FileZilla-compatible `sites.xml` with `<Server>` entries, base64-encoded passwords, and preserved folder structure

The parser automatically handles:
- **Encoding detection** — UTF-8 vs Windows-1254 (Turkish), with BOM stripping
- **URL decoding** — Non-ASCII characters stored as `%XX` in WinSCP.ini are decoded
- **Scheme prefix filtering** — `http:`, `https:`, `ftp:` etc. prefixes are stripped from session names
- **Mixed path separators** — Both `\` and `/` are recognized as folder separators

## Technology Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.21+ |
| GUI | [Gio](https://gioui.org) (pure Go, no CGO) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| INI parsing | [go-ini/ini](https://gopkg.in/ini.v1) |
| XML generation | [etree](https://github.com/beevik/etree) |
| Build | [GoReleaser](https://goreleaser.com) |

## License

MIT
