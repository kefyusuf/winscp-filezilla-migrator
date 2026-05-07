# Phase 3: GUI - Plan

**Phase:** 3 (GUI)
**Goal:** Build Fyne cross-platform UI with file picker, server list, and migration controls
**Status:** Planned

## Tasks

### T-01: Initialize Fyne App
**Description:** Set up Fyne application with main window
**Files:**
- `ui/app.go` — Fyne app initialization
- `ui/windows/main.go` — Main window setup
**Verification:** App window opens without errors

### T-02: Create File Picker Widget
**Description:** File selection dialog for WinSCP.ini
**File:** `ui/widgets/file_picker.go`
**Features:**
- Use fyne dialog for file open
- Filter for .ini files
- Display selected file path
**Verification:** Can browse and select INI file

### T-03: Create Server List Widget
**Description:** Tree view showing folder hierarchy and servers
**File:** `ui/widgets/server_list.go`
**Features:**
- Fyne Tree widget for folder structure
- Show server name, host, protocol badge
- Support selection for migration
**Verification:** Displays parsed sessions with folders

### T-04: Create Preview Panel Widget
**Description:** Details panel for selected server
**File:** `ui/widgets/preview.go`
**Features:**
- Show: Host, Username, Protocol, Port, RemoteDir
- Display in side panel when server selected
- Handle empty selection
**Verification:** Shows correct server details

### T-05: Create Migration Controls
**Description:** Migrate button and progress indicator
**File:** `ui/widgets/controls.go`
**Features:**
- Migrate button triggers export
- Progress indicator during migration
- Success/error message display
**Verification:** Migration completes with feedback

### T-06: Create Save Dialog
**Description:** File save dialog for output location
**File:** Part of `ui/widgets/controls.go`
**Features:**
- Use fyne dialog for file save
- Default name: sites.xml
- Filter for .xml files
**Verification:** Can select save location

### T-07: Integrate Core Logic with GUI
**Description:** Connect parser/exporter to GUI
**File:** `ui/app.go` (update)
**Features:**
- Call parser.ParseWinSCPIni() on file selection
- Call exporter.ExportToFileZilla() on migrate
- Handle errors from core logic
**Verification:** Full end-to-end works

### T-08: Add Error Handling
**Description:** Display errors in GUI
**File:** `ui/widgets/error_dialog.go`
**Features:**
- Dialog for parsing errors
- Dialog for export errors
- Specific error messages
**Verification:** Errors display correctly

---

## Success Criteria

| # | Criterion |
|---|-----------|
| 1 | Fyne app initializes and shows window |
| 2 | File picker can select .ini files |
| 3 | Server list shows folders and servers |
| 4 | Preview panel shows selected server details |
| 5 | Migrate button triggers export |
| 6 | Save dialog allows output location selection |
| 7 | Full migration works end-to-end |
| 8 | Errors display in dialogs |

---

## Dependencies

- T-01 (Fyne init) → T-02 (File picker)
- T-02 → T-03 (Server list, needs parsed sessions)
- T-03 → T-04 (Preview, needs selection)
- T-04 + T-06 → T-05 (Controls + save dialog)
- T-05 → T-07 (Integrate core logic)
- T-07 → T-08 (Error handling)

---

*Plan created: 2026-05-07*