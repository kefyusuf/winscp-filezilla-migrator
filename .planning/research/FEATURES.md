# Features Research — WinSCP2FileZilla v2

## Table Stakes (Must Have)

These features are expected — users won't tolerate their absence:

1. **INI File Selection** — Browse and select WinSCP.ini file
2. **Server List Parsing** — Parse all sessions from INI with folder hierarchy
3. **Password Decryption** — Decrypt WinSCP's custom XOR-based password encryption
4. **FileZilla Export** — Generate valid FileZilla sites.xml
5. **Server Preview** — Show parsed servers before migration
6. **Single Server Migration** — Migrate individual servers (not just batch)

## Differentiators (Competitive Advantage)

Features that make this tool stand out:

1. **Cross-Platform GUI** — Works on Windows, Linux, macOS (original is CLI-only)
2. **Selective Migration** — Choose specific servers to migrate
3. **Migration Preview** — Review server details before export
4. **RemoteDir/LocalDir** — Migrate default directory paths
5. **Progress Indicator** — Visual feedback during migration

## Anti-Features (Deliberately NOT Building)

Things we're explicitly not doing in v1:

1. **Reverse Migration** — FileZilla → WinSCP
2. **Cloud Sync** — No cloud storage integration
3. **Multi-File Import** — Single INI file at a time
4. **Advanced Protocol Options** — FTP/SFTP only, no TLS/SSH key management

## Feature Complexity Notes

- **Password Decryption**: High complexity — custom XOR algorithm with edge cases
- **INI Parsing**: Medium — nested folder structure requires recursive handling
- **FileZilla XML**: Low — well-documented format, straightforward generation

## Dependencies Between Features

- Password decryption depends on: HostName + UserName (key derivation)
- Server list depends on: INI parsing + password decryption
- Export depends on: Server list parsing

## Confidence Level: HIGH

FTP migration is well-understood domain with clear requirements.