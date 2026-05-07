# Stack Research — WinSCP2FileZilla v2

## Technology Choices

### Core Language
- **Go 1.21+** — Required for Fyne compatibility
- No external dependencies for core logic (keep it lightweight)

### GUI Framework
- **Fyne 2.x** — Pure Go cross-platform GUI
  - v2.4+ recommended for stability
  - Native look on Windows, Linux, macOS
  - No C bindings, no webview
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

Fyne is well-established for cross-platform Go apps. Stack is proven and stable.