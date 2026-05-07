# Phase 2: Core - Context

**Gathered:** 2026-05-07
**Status:** Ready for planning
**Source:** Auto mode — all gray areas selected

<domain>
## Phase Boundary

Implement core functionality: WinSCP.ini parsing, password decryption, FileZilla XML export. Phase 2 delivers the fundamental migration logic that both CLI and GUI will use.

</domain>

<decisions>
## Implementation Decisions

### INI Parsing (CORE-01)
- **D-01:** Use `gopkg.in/ini.v1` package for INI parsing
- **D-02:** Parse Sessions\ prefix from section names to extract folder structure
- **D-03:** Handle nested folders with recursive path splitting

### Password Decryption (CORE-02)
- **D-04:** Preserve existing XOR algorithm from original project (muety/winscp2filezilla)
- **D-05:** Decrypt function requires: host, username, encrypted password
- **D-06:** Handle empty password edge case (return empty string)

### FileZilla XML Export (CORE-03)
- **D-07:** Use `beevik/etree` for XML generation
- **D-08:** Root element: FileZilla3 > Servers > Folder hierarchy
- **D-09:** Password stored as base64-encoded string with encoding="base64" attribute
- **D-10:** Default directories (LocalDir, RemoteDir) set to empty in v1

### Folder Hierarchy (CORE-04)
- **D-11:** Sessions stored as "Sessions\Folder\Server" — split by backslash
- **D-12:** Folder hierarchy recreated in FileZilla XML with Folder elements

### Protocol Mapping (CORE-05)
- **D-13:** WinSCP FSProtocol "2" = SFTP → FileZilla Protocol "1"
- **D-14:** WinSCP FSProtocol "5" or empty = FTP → FileZilla Protocol "0"
- **D-15:** Default port: 22 for SFTP, 21 for FTP

</decisions>

<canonical_refs>
## Canonical References

- `.planning/ROADMAP.md` — Phase 2 specifications
- `.planning/research/STACK.md` — Technology stack (go-ini, etree)
- `.planning/REQUIREMENTS.md` — CORE-01 to CORE-05 requirements

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Original project password.go has working decrypt algorithm — use as reference
- Original project main.go has INI parsing and XML export patterns

### Patterns
- go-ini loads entire INI file, iterate sections for Sessions\ prefix
- etree creates document, adds elements, WriteToFile for output

</code_context>

<specifics>
## Specific Ideas

[auto] INI parser: go-ini library, iterate sections with Sessions\ prefix
[auto] Password: Use existing XOR algorithm from original project
[auto] XML: etree with proper structure (FileZilla3 > Servers > Folders > Server)
[auto] Protocol: FSProtocol 2=SFTP(1), 5+=FTP(0), default=FTP(0)

</specifics>

<deferred>
## Deferred Ideas

- RemoteDir/LocalDir migration — Phase 4 (Advanced)
- Multiple protocol support beyond FTP/SFTP — out of scope

</deferred>

---

*Phase: 02-core*
*Context gathered: 2026-05-07 via auto mode*