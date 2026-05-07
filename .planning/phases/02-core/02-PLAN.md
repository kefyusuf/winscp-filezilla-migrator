# Phase 2: Core - Plan

**Phase:** 2 (Core)
**Goal:** Implement INI parsing, password decryption, FileZilla XML export
**Status:** Planned

## Tasks

### T-01: Create Domain Models
**Description:** Define data structures for sessions, servers, folders
**Files:**
- `domain/models/session.go` — Session struct (Name, HostName, UserName, Password, PortNumber, FSProtocol)
- `domain/models/server.go` — Server struct for FileZilla format
- `domain/models/folder.go` — Folder struct for hierarchy
**Verification:** All structs compile without errors

### T-02: Implement INI Parser
**Description:** Parse WinSCP.ini with folder structure extraction
**File:** `domain/parser/ini.go`
**Functions:**
- `ParseWinSCPIni(filepath string) ([]Session, error)` — Read and parse INI file
- Handle Sessions\ prefix in section names
- Extract folder hierarchy from session names
**Verification:** Parse test INI file successfully

### T-03: Implement Password Decryption
**Description:** Replicate WinSCP's XOR-based password decryption
**File:** `domain/parser/decrypt.go`
**Functions:**
- `Decrypt(host, username, password string) string` — Decrypt password
- Preserve original algorithm (PW_MAGIC = 0xA3, PW_FLAG = 0xFF)
**Verification:** Match known password decryption results

### T-04: Implement FileZilla XML Exporter
**Description:** Generate FileZilla sites.xml format
**File:** `domain/exporter/filezilla.go`
**Functions:**
- `ExportToFileZilla(sessions []Session, outputPath string) error` — Write XML
- Create proper FileZilla3 > Servers > Folder hierarchy
- Encode passwords as base64 with encoding="base64" attribute
**Verification:** Generated XML can be imported in FileZilla

### T-05: Implement Protocol Mapping
**Description:** Map WinSCP protocols to FileZilla protocol numbers
**Location:** Part of `domain/exporter/filezilla.go`
**Logic:**
- FSProtocol "2" → SFTP → Protocol "1"
- FSProtocol "5" or empty → FTP → Protocol "0"
- Port mapping: SFTP=22, FTP=21
**Verification:** Correct protocol numbers in output

### T-06: Create Service Layer
**Description:** Orchestrate parsing, decryption, and export
**File:** `app/service.go` (update existing)
**Functions:**
- `Migrate(inputPath, outputPath string) error` — Full migration flow
**Verification:** Full end-to-end migration works

### T-07: Add CLI Migration Command
**Description:** Add migrate command to CLI
**File:** Update `app/service.go`
**Command:** `winscp2filezilla migrate --in <ini> --out <xml>`
**Verification:** CLI command works end-to-end

### T-08: Create Test INI File
**Description:** Create sample WinSCP.ini for testing
**File:** `testdata/test.ini`
**Content:** Sample sessions with different protocols and folder structures
**Verification:** Can be parsed and migrated successfully

---

## Success Criteria

| # | Criterion |
|---|-----------|
| 1 | Domain models defined and compile |
| 2 | INI parser extracts all sessions with folder structure |
| 3 | Passwords decrypt correctly (match original tool) |
| 4 | FileZilla XML generates valid, importable output |
| 5 | Protocol mapping produces correct FileZilla values |
| 6 | CLI migrate command works end-to-end |
| 7 | Test INI file migrates successfully |

---

## Dependencies

- T-01 (Models) → T-02 (INI Parser)
- T-02 + T-03 (Parser + Decrypt) → T-04 (Exporter)
- T-04 → T-05 (Protocol Mapping, same file)
- T-04 → T-06 (Service Layer)
- T-06 → T-07 (CLI Command)
- T-07 + Test INI → T-08 (Testing)

---

*Plan created: 2026-05-07*