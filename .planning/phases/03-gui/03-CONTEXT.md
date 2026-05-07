# Phase 3: GUI - Context

**Gathered:** 2026-05-07
**Status:** Ready for planning
**Source:** Auto mode — all gray areas selected

<domain>
## Phase Boundary

Build Fyne-based cross-platform GUI for winscp2filezilla. Phase 3 delivers the desktop application with file picker, server list preview, and migration controls.

</domain>

<decisions>
## Implementation Decisions

### GUI Framework (GUI-01)
- **D-01:** Use fyne.io/fyne/v2 for native cross-platform UI
- **D-02:** Support Windows, Linux, macOS with native look
- **D-03:** Use Fyne's built-in theme (auto-detect system theme)

### File Picker (GUI-01)
- **D-04:** Use fyne dialog for file selection
- **D-05:** Filter for .ini files only
- **D-06:** Remember last used directory

### Server List (GUI-02)
- **D-07:** Display sessions in Tree widget showing folder hierarchy
- **D-08:** Show server name, host, protocol badge in list
- **D-09:** Support selecting/deselecting servers for migration
- **D-10:** Show session count per folder

### Server Preview (GUI-03)
- **D-11:** Show details panel on server selection
- **D-12:** Display: Host, Username, Protocol, Port, RemoteDir
- **D-13:** Show password (masked) and edit option (future)

### Migration Controls (GUI-04, GUI-05)
- **D-14:** Migrate button triggers export process
- **D-15:** Progress indicator during migration
- **D-16:** Save dialog for output file location
- **D-17:** Success/error message after completion

### Error Handling (GUI-06)
- **D-18:** Show dialog for parsing errors
- **D-19:** Display specific error messages (file not found, invalid format)
- **D-20:** Allow retry after error

</decisions>

<canonical_refs>
## Canonical References

- `.planning/ROADMAP.md` — Phase 3 specifications
- `.planning/research/STACK.md` — Fyne framework
- `.planning/REQUIREMENTS.md` — GUI-01 to GUI-06 requirements

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- domain/parser/ini.go — INI parsing (reuse in GUI)
- domain/exporter/filezilla.go — XML export (reuse in GUI)
- app/service.go — has migrate logic, can be called from GUI

### Patterns
- Fyne app setup: fyne.App.NewWindow(), window.SetContent()
- Tree widget for hierarchical data
- Dialog for file open/save

</code_context>

<specifics>
## Specific Ideas

[auto] File picker: Use fyne dialog with .ini filter
[auto] Server list: Tree widget with folder hierarchy
[auto] Preview: Side panel with selected server details
[auto] Migration: Button + progress + save dialog
[auto] Theme: Built-in Fyne theme, auto-detect

</specifics>

<deferred>
## Deferred Ideas

- Dark mode toggle — Phase 5 (Polish)
- Multiple file selection — Phase 2+
- SSH key management — Out of scope

</deferred>

---

*Phase: 03-gui*
*Context gathered: 2026-05-07 via auto mode*