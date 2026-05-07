# Phase 1: Setup - Plan

**Phase:** 1 (Setup)
**Goal:** Initialize project structure, dependencies, and CI/CD
**Status:** Planned

## Tasks

### T-01: Initialize Go Module
**Description:** Initialize Go module with proper naming and version
**Command:** `go mod init github.com/muety/winscp2filezilla`
**Expected Output:** `go.mod` file created
**Verification:** Module name is `github.com/muety/winscp2filezilla`

### T-02: Create Directory Structure
**Description:** Create layered Go project structure
**Commands:**
```bash
mkdir -p domain/models domain/parser domain/exporter app cmd/cli ui/windows ui/widgets
```
**Expected Output:** All directories created
**Verification:** Run `ls -la` to confirm structure

### T-03: Add Dependencies
**Description:** Add required Go dependencies
**Commands:**
```bash
go get fyne.io/fyne/v2
go get github.com/spf13/cobra
go get github.com/go-ini/ini
go get github.com/beevik/etree
```
**Expected Output:** Dependencies added to go.mod
**Verification:** `go mod tidy` succeeds

### T-04: Create Basic App Structure
**Description:** Create minimal runnable application
**Files to create:**
- `main.go` — Entry point with Fyne app initialization
- `cmd/cli/main.go` — CLI entry point with cobra
**Verification:** `go build -o winscp2filezilla .` succeeds

### T-05: Create GitHub Actions Workflow
**Description:** CI/CD pipeline for multi-platform builds
**File:** `.github/workflows/build.yml`
**Content:** Build on push/PR for Windows, Linux, macOS
**Verification:** Workflow file created in `.github/workflows/`

### T-06: Verify Build Pipeline
**Description:** Ensure the project builds cleanly
**Command:** `go build -o winscp2filezilla .`
**Expected Output:** Binary created without errors
**Verification:** `./winscp2filezilla --help` shows CLI help

---

## Success Criteria

| # | Criterion |
|---|-----------|
| 1 | Go module initialized with correct name |
| 2 | Directory structure follows Go conventions |
| 3 | All dependencies resolve correctly |
| 4 | Basic application compiles without errors |
| 5 | GitHub Actions workflow exists |

---

*Plan created: 2026-05-07*