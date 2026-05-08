<!-- GSD:project-start source:PROJECT.md -->
## Project

**winscp-filezilla-migrator**

A cross-platform desktop application that migrates saved FTP/SFTP server configurations from WinSCP to FileZilla. Users select their WinSCP.ini file, preview their server list, and export to FileZilla's sites.xml format. Built with Go + Fyne for native performance on Windows, Linux, and macOS.

**Core Value:** Enable users to seamlessly migrate their WinSCP server configurations to FileZilla without data loss, with a simple GUI that works across all desktop platforms.

### Constraints

- **[Tech Stack]**: Go + Fyne GUI — Required for cross-platform native performance
- **[Compatibility]**: Windows, Linux, macOS — Must work on all three platforms
- **[Data Integrity]**: Passwords must decrypt correctly — Zero tolerance for data loss
- **[Existing Code]**: Must preserve WinSCP password decryption logic — Custom algorithm required
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Technology Choices
### Core Language
- **Go 1.21+** — Required for Fyne compatibility
- No external dependencies for core logic (keep it lightweight)
### GUI Framework
- **Fyne 2.x** — Pure Go cross-platform GUI
- **Why Fyne**: Single binary, no runtime deps, true cross-platform
### CLI Framework
- **Cobra** — For optional CLI mode alongside GUI
- Provides: subcommands, flags, help generation
### Data Parsing
- **go-ini/ini** — INI file parsing (same as original project)
- **etree** — XML generation for FileZilla export (same as original)
### Testing
- **testify** — Assertion library
- **Golden files** — For INI/XML parsing tests
### Build & Release
- **GoReleaser** — Multi-platform builds (.exe, .AppImage, .deb, .dmg)
- **GitHub Actions** — CI/CD pipeline
## Recommended Versions
| Component | Version | Rationale |
|-----------|---------|-----------|
| Go | 1.21+ | Fyne v2.4 requires Go 1.21+ |
| Fyne | 2.4.x | Stable, well-tested |
| Cobra | 1.7.x | Latest stable |
| go-ini | latest | Active maintenance |
| etree | latest | XML generation |
## What NOT to Use
- **Electron/loz** — Too heavy, requires Node.js
- **Qt (lorca)** — External Qt dependency
- **Wails** — Good but Go-only is simpler for this use case
- **gioui** — More complex API, less mature than Fyne for desktop
## Confidence Level: HIGH
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
