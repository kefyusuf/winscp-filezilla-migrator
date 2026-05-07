# Architecture Research — WinSCP2FileZilla v2

## Component Boundaries

### Layer 1: Core / Domain
```
domain/
├── models/
│   ├── session.go      # Server session data
│   ├── folder.go       # Folder hierarchy
│   └── server.go       # Server config
├── parser/
│   ├── ini.go          # WinSCP INI parsing
│   └── decoder.go      # Password decryption
└── exporter/
    └── filezilla.go    # FileZilla XML generation
```

### Layer 2: Application / Logic
```
app/
├── service.go          # Orchestration layer
├── validator.go        # Input validation
└── converter.go        # Data transformation
```

### Layer 3: UI / Presentation
```
ui/
├── app.go              # Fyne app setup
├── windows/
│   ├── main.go         # Main window
│   └── preview.go      # Server preview
├── widgets/
│   ├── file_picker.go  # File selection
│   ├── server_list.go  # Server list table
│   └── migrate_btn.go # Migration controls
└── theme.go            # Fyne theme setup
```

## Data Flow

```
[WinSCP.ini file]
    → [INI Parser] → [Session models]
    → [Password Decoder] → [Decrypted sessions]
    → [Server List Widget] → [User preview/selection]
    → [Exporter] → [FileZilla XML]
    → [Save dialog] → [sites.xml file]
```

## Build Order (Dependency Graph)

```
Phase 1: Setup & Core
├── go mod init
├── directory structure
├── basic models (Session, Server, Folder)
└── INI parsing

Phase 2: Password Decryption
├── Decrypt() function
├── Test with known passwords
└── Edge case handling

Phase 3: Export Logic
├── FileZilla XML generation
├── Folder structure handling
└── Valid XML output

Phase 4: CLI Layer (optional)
├── cobra commands
└── CLI + GUI can coexist

Phase 5: Fyne UI
├── App scaffolding
├── File picker widget
├── Server list table
├── Preview window
└── Migration button

Phase 6: Polish
├── Error handling
├── Cross-platform builds
└── Release artifacts
```

## Component Dependencies

- **Password Decoder** depends on nothing (pure logic)
- **INI Parser** depends on domain models
- **FileZilla Exporter** depends on domain models
- **UI** depends on domain + application layers

## Confidence Level: HIGH

Standard layered architecture with clear separation. Similar to original CLI structure, just organized better.