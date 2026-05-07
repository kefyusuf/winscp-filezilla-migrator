# Research Summary — WinSCP2FileZilla v2

## Key Findings

### Stack
- **Language**: Go 1.21+
- **GUI**: Fyne 2.4.x (pure Go, cross-platform, native look)
- **CLI**: Cobra 1.7.x (optional CLI alongside GUI)
- **Parsing**: go-ini, etree (same as original project)
- **Build**: GoReleaser for multi-platform artifacts

### Table Stakes (Must Have)
1. INI file selection and parsing
2. Server session extraction with folder hierarchy
3. WinSCP password decryption (custom XOR algorithm)
4. FileZilla XML export with correct format
5. Server preview before migration

### Differentiators
1. Cross-platform GUI (original is CLI-only)
2. Selective server migration (choose specific servers)
3. RemoteDir/LocalDir migration
4. Visual progress indicator

### Architecture
- **Pattern**: Layered architecture (domain → app → ui)
- **Build order**: Core (Phase 2) → Export (Phase 3) → CLI (Phase 4) → GUI (Phase 5)

### Watch Out For
1. **Password decryption**: Edge cases with empty passwords, unicode usernames
2. **INI folder parsing**: Handle nested folders, special characters
3. **Fyne cross-platform**: Test on all three platforms, use built-in theme
4. **FileZilla XML**: Correct protocol mapping (FTP=0, SFTP=1), base64 encoding

### Out of Scope
- Reverse migration (FileZilla → WinSCP)
- Cloud sync features
- Advanced protocol options (TLS/SSH keys)

## Files
- `.planning/research/STACK.md` — Technology choices
- `.planning/research/FEATURES.md` — Feature categories
- `.planning/research/ARCHITECTURE.md` — Component structure
- `.planning/research/PITFALLS.md` — Common mistakes to avoid

---

*Research completed: 2026-05-07*