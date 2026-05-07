# Roadmap — WinSCP2FileZilla v2

**[5] phases** | **[15] requirements** | All v1 requirements covered ✓

| # | Phase | Goal | Requirements | Success Criteria |
|---|-------|------|--------------|------------------|
| 1 | **Setup** | Initialize project structure, dependencies, CI/CD | — | Go module initialized, directory structure created, build pipeline works |
| 2 | **Core** | Implement INI parsing, password decryption, XML export | CORE-01, CORE-02, CORE-03, CORE-04, CORE-05 | Can parse real WinSCP.ini, decrypt passwords, generate valid FileZilla XML |
| 3 | **GUI** | Build Fyne cross-platform UI | GUI-01, GUI-02, GUI-03, GUI-04, GUI-05, GUI-06 | File picker works, server list displays, migration completes |
| 4 | **Advanced** | Add RemoteDir/LocalDir migration | DIR-01, DIR-02 | Directory paths correctly migrated in output XML |
| 5 | **Polish** | Finalize builds, tests, documentation | CROSS-01, CROSS-02, CROSS-03 | Windows/Linux/macOS builds work, basic tests pass |

---

## Phase Details

### Phase 1: Setup
**Goal:** Initialize project structure, dependencies, CI/CD

**Requirements:** None (infrastructure only)

**Success criteria:**
1. Go module initialized with `go mod init`
2. Directory structure created (domain/app/ui layers)
3. Basic Makefile for building
4. GitHub Actions workflow for CI
5. Dependencies (Fyne, Cobra, go-ini, etree) added

**Dependencies:** None (this is phase 1)

---

### Phase 2: Core
**Goal:** Implement INI parsing, password decryption, XML export

**Requirements:**
- [ ] **CORE-01**: Parse WinSCP.ini sessions and folder structure
- [ ] **CORE-02**: Decrypt WinSCP password encryption (XOR-based algorithm)
- [ ] **CORE-03**: Export FileZilla sites.xml with correct format
- [ ] **CORE-04**: Handle folder hierarchy in server names (Sessions\Folder\Server)
- [ ] **CORE-05**: Map WinSCP protocols to FileZilla protocol numbers (FTP=0, SFTP=1)

**Success criteria:**
1. Can parse real WinSCP.ini files with multiple sessions
2. Password decryption produces same results as original tool
3. Generated XML can be imported in FileZilla
4. Folder structure preserved in output
5. FTP and SFTP protocols correctly mapped

**Dependencies:** Phase 1 (Setup)

**Verification:**
- Parse test INI files
- Compare password decryption with known values
- Import generated XML in FileZilla

---

### Phase 3: GUI
**Goal:** Build Fyne cross-platform UI

**Requirements:**
- [ ] **GUI-01**: File picker for WinSCP.ini selection
- [ ] **GUI-02**: Server list display with folder hierarchy
- [ ] **GUI-03**: Server preview (show host, username, protocol, port)
- [ ] **GUI-04**: Migration button with progress indicator
- [ ] **GUI-05**: File save dialog for output (sites.xml)
- [ ] **GUI-06**: Error display for invalid INI or parsing failures

**Success criteria:**
1. Can browse and select WinSCP.ini file
2. Server list shows all sessions with folder structure
3. Clicking server shows details in preview panel
4. Migration button triggers export process
5. Can save output as sites.xml
6. Errors displayed clearly when parsing fails

**Dependencies:** Phase 2 (Core)

**UI Design:**
```
┌─────────────────────────────────────────────┐
│  [Open INI File]  selected.ini              │
├─────────────────────────────────────────────┤
│  Servers:                          Details:  │
│  ├─ Folder1                      Host: ...  │
│  │  └─ server1                    User: ...  │
│  ├─ Folder2                      Protocol: .. │
│  │  └─ server2                    Port: ...  │
│  └─ server3                                  │
├─────────────────────────────────────────────┤
│           [ Migrate to FileZilla ]          │
└─────────────────────────────────────────────┘
```

---

### Phase 4: Advanced
**Goal:** Add RemoteDir/LocalDir migration

**Requirements:**
- [ ] **DIR-01**: Migrate RemoteDir path from WinSCP
- [ ] **DIR-02**: Migrate LocalDir path from WinSCP

**Success criteria:**
1. RemoteDir from WinSCP.ini appears in FileZilla XML
2. LocalDir from WinSCP.ini appears in FileZilla XML
3. Paths correctly mapped to FileZilla's XML structure

**Dependencies:** Phase 3 (GUI)

---

### Phase 5: Polish
**Goal:** Finalize builds, tests, documentation

**Requirements:**
- [ ] **CROSS-01**: Windows .exe build
- [ ] **CROSS-02**: Linux AppImage/.deb build
- [ ] **CROSS-03**: macOS .dmg build

**Success criteria:**
1. Windows .exe runs without errors
2. Linux AppImage/deb builds available
3. macOS .dmg builds available
4. Basic unit tests for core logic pass
5. README.md updated with usage instructions
6. Version tag created for release

**Dependencies:** Phase 4 (Advanced)

---

## Dependencies Between Phases

```
Phase 1 (Setup)
    ↓
Phase 2 (Core) ← depends on Phase 1
    ↓
Phase 3 (GUI) ← depends on Phase 2
    ↓
Phase 4 (Advanced) ← depends on Phase 3
    ↓
Phase 5 (Polish) ← depends on Phase 4
```

---

## Key Decisions Log

| Phase | Decision | Rationale |
|-------|----------|-----------|
| 1 | Fyne for GUI | Pure Go, no external deps, cross-platform native |
| 2 | Preserve existing decrypt algorithm | Tested and works, no need to reinvent |
| 3 | Preview before migrate | User requested feature, good UX |

---

*Roadmap created: 2026-05-07*