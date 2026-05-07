# Pitfalls Research — WinSCP2FileZilla v2

## Domain-Specific Pitfalls

### 1. Password Decryption Edge Cases
**Warning signs:**
- Password decryption produces wrong results for certain host/user combinations
- Empty passwords not handled
- Unicode usernames cause issues

**Prevention:**
- Use the exact algorithm from original code (XOR-based with magic constants)
- Test with multiple real-world INI files
- Handle empty password case explicitly (return empty string)

**Phase to address:** Phase 2 (Core - Password Decryption)

### 2. INI Folder Hierarchy Parsing
**Warning signs:**
- Nested folders not preserved in output
- Folder names with special characters break

**Prevention:**
- Recursive parsing of Sessions\ prefix
- Handle "/" and "\" in folder names
- Test with deeply nested folder structures

**Phase to address:** Phase 2 (Core - INI Parsing)

### 3. Fyne Cross-Platform Issues
**Warning signs:**
- UI looks different on Windows vs Linux
- Fonts not loading correctly
- Window sizing issues

**Prevention:**
- Use Fyne's built-in theme (auto-detect system)
- Test on all three platforms
- Use flexible layouts that adapt to window size

**Phase to address:** Phase 4 (GUI Implementation)

### 4. FileZilla XML Compatibility
**Warning signs:**
- Generated XML doesn't import in FileZilla
- Protocol numbers wrong (FTP=0, SFTP=1)
- Base64 encoding missing for passwords

**Prevention:**
- Reference actual FileZilla XML structure
- Protocol mapping: WinSCP "2" = SFTP, "5" = FTP (actually depends on FSProtocol)
- Add proper encoding attributes
- Validate against FileZilla's import requirements

**Phase to address:** Phase 3 (Export Logic)

### 5. Cross-Platform Path Handling
**Warning signs:**
- Hardcoded Windows paths
- Forward vs backward slash issues
- Home directory detection fails on Linux/macOS

**Prevention:**
- Use filepath.Join() for all path operations
- Use os/user package for home directory
- Test path operations on all platforms

**Phase to address:** Phase 1 (Setup)

## Generic Pitfalls (Avoid)

- **Premature optimization**: Don't optimize before profiling
- **Over-engineering**: Simple domain, don't overcomplicate
- **Skipping tests**: Password logic is fragile, needs test coverage
- **Ignoring GoReleaser**: Multi-platform builds manual = pain

## Warning Signs to Watch For

1. Password decrypts differently than original tool → Review algorithm
2. Fyne app crashes on startup → Check Fyne version compatibility
3. XML doesn't import in FileZilla → Validate XML structure
4. Build fails on one platform → Check platform-specific code

## Phase Mapping

| Pitfall | Phase |
|---------|-------|
| Password decryption edge cases | Phase 2 |
| INI folder hierarchy | Phase 2 |
| Fyne cross-platform | Phase 4 |
| FileZilla XML compatibility | Phase 3 |
| Path handling | Phase 1 |

## Confidence Level: MEDIUM

These pitfalls are well-documented for similar migration tools. Original Go code already solves most of these.