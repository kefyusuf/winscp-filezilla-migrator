# Phase 1: Setup - Context

**Gathered:** 2026-05-07
**Status:** Ready for planning

<domain>
## Phase Boundary

Initialize project structure, dependencies, and CI/CD for winscp2filezilla v2. This phase sets up the foundation for all subsequent phases.

</domain>

<decisions>
## Implementation Decisions

### Project Setup
- **D-01:** Use `github.com/muety/winscp2filezilla` as module name (preserves original project reference)
- **D-02:** Use standard Go project layout (root-level packages, no src/ directory)
- **D-03:** Use Go 1.21+ as minimum version (Fyne v2.4 requirement)
- **D-04:** No Makefile needed — Go build commands are simple enough

### Dependencies
- **D-05:** `fyne.io/fyne/v2` for GUI
- `github.com/spf13/cobra` for CLI
- `github.com/go-ini/ini` for INI parsing
- `github.com/beevik/etree` for XML generation

### CI/CD
- **D-06:** GitHub Actions for CI/CD (standard for Go projects)
- **D-07:** Workflow should run on push to main/master and PRs
- **D-08:** Build for Windows, Linux, macOS on release

### Directory Structure
- `domain/models/` — Data structures (Session, Server, Folder)
- `domain/parser/` — INI parsing logic
- `domain/exporter/` — FileZilla XML generation
- `app/` — Application orchestration
- `cmd/cli/` — CLI entry point
- `ui/` — Fyne GUI components

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/ROADMAP.md` — Phase 1 specifications
- `.planning/research/STACK.md` — Technology stack choices

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — this is greenfield refactor

### Established Patterns
- Original project (muety/winscp2filezilla) used flat structure — new structure organizes better

### Integration Points
- After setup, Phase 2 will use domain/models, domain/parser

</code_context>

<specifics>
## Specific Ideas

[auto] Module naming: Keep original project reference for compatibility
[auto] Directory structure: Standard Go layout (not src/)
[auto] CI: GitHub Actions
[auto] Build: Direct Go commands, no Makefile

</specifics>

<deferred>
## Deferred Ideas

[If none: "None — discussion stayed within phase scope"]

</deferred>

---

*Phase: 01-setup*
*Context gathered: 2026-05-07*