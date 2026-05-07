# Requirements — WinSCP2FileZilla v2

## v1 Requirements

### Core Logic
- [ ] **CORE-01**: Parse WinSCP.ini sessions and folder structure
- [ ] **CORE-02**: Decrypt WinSCP password encryption (XOR-based algorithm)
- [ ] **CORE-03**: Export FileZilla sites.xml with correct format
- [ ] **CORE-04**: Handle folder hierarchy in server names (Sessions\Folder\Server)
- [ ] **CORE-05**: Map WinSCP protocols to FileZilla protocol numbers (FTP=0, SFTP=1)

### GUI Components
- [ ] **GUI-01**: File picker for WinSCP.ini selection
- [ ] **GUI-02**: Server list display with folder hierarchy
- [ ] **GUI-03**: Server preview (show host, username, protocol, port)
- [ ] **GUI-04**: Migration button with progress indicator
- [ ] **GUI-05**: File save dialog for output (sites.xml)
- [ ] **GUI-06**: Error display for invalid INI or parsing failures

### Cross-Platform
- [ ] **CROSS-01**: Windows .exe build (main target)
- [ ] **CROSS-02**: Linux AppImage/.deb build
- [ ] **CROSS-03**: macOS .dmg build

### Directory Migration
- [ ] **DIR-01**: Migrate RemoteDir path from WinSCP
- [ ] **DIR-02**: Migrate LocalDir path from WinSCP

---

## v2 Requirements (Deferred)

- [ ] **UI-01**: Dark mode support
- [ ] **EXP-01**: Batch import multiple INI files
- [ ] **REV-01**: Reverse migration (FileZilla → WinSCP)
- [ ] **SYNC-01**: Cloud sync integration

---

## Out of Scope

- [FileZilla → WinSCP reverse migration] — Out of scope for v1, complex password encryption
- [SSH key management] — Beyond basic SFTP support
- [TLS/SSL configuration] — Use FileZilla defaults
- [Cloud storage sync] — Not in initial scope

---

## Traceability

| REQ-ID | Phase | Status |
|--------|-------|--------|
| CORE-01 | 2 | pending |
| CORE-02 | 2 | pending |
| CORE-03 | 3 | pending |
| CORE-04 | 2 | pending |
| CORE-05 | 3 | pending |
| GUI-01 | 4 | pending |
| GUI-02 | 4 | pending |
| GUI-03 | 4 | pending |
| GUI-04 | 4 | pending |
| GUI-05 | 4 | pending |
| GUI-06 | 4 | pending |
| CROSS-01 | 6 | pending |
| CROSS-02 | 6 | pending |
| CROSS-03 | 6 | pending |
| DIR-01 | 5 | pending |
| DIR-02 | 5 | pending |

---

*Requirements defined: 2026-05-07*