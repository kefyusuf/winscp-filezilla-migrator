# winscp2filezilla

## What This Is

A cross-platform desktop application that migrates saved FTP/SFTP server configurations from WinSCP to FileZilla. Users select their WinSCP.ini file, preview their server list, and export to FileZilla's sites.xml format. Built with Go + Fyne for native performance on Windows, Linux, and macOS.

## Core Value

Enable users to seamlessly migrate their WinSCP server configurations to FileZilla without data loss, with a simple GUI that works across all desktop platforms.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Parse WinSCP.ini sessions and folder structure
- [ ] Decrypt WinSCP password encryption
- [ ] Export FileZilla sites.xml format
- [ ] Fyne-based cross-platform GUI
- [ ] File picker for WinSCP.ini selection
- [ ] Server list preview with selection
- [ ] Migration progress indicator
- [ ] File save dialog for output
- [ ] RemoteDir/LocalDir path migration
- [ ] Windows .exe build
- [ ] Linux AppImage/deb build
- [ ] macOS .dmg build

### Out of Scope

- [FileZilla → WinSCP reverse migration] — Out of scope for v1
- [Cloud sync feature] — Not in initial scope
- [Dark mode] — Deferred to v2

## Context

Original project (muety/winscp2filezilla) is a Go CLI tool that:
- Reads WinSCP.ini from registry or file
- Parses Sessions\ folders and server entries
- Decrypts passwords using custom XOR algorithm
- Exports FileZilla-compatible XML

The existing Go code provides the core logic. This project restructures it with:
- Proper modular architecture (core/logic/gui layers)
- Fyne GUI framework for native cross-platform UI
- Proper error handling and validation
- Comprehensive test coverage
- CI/CD for multi-platform builds

## Constraints

- **[Tech Stack]**: Go + Fyne GUI — Required for cross-platform native performance
- **[Compatibility]**: Windows, Linux, macOS — Must work on all three platforms
- **[Data Integrity]**: Passwords must decrypt correctly — Zero tolerance for data loss
- **[Existing Code]**: Must preserve WinSCP password decryption logic — Custom algorithm required

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Fyne over webview | Pure Go, no external dependencies, native look | — Pending |
| RemoteDir/LocalDir in v1 | Frequently requested, high user value | — Pending |
| Preview before migrate | Allows users to review/select servers | — Pending |

---

*Last updated: 2026-05-07 after initialization*