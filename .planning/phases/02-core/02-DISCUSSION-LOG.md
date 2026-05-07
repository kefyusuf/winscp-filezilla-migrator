# Phase 2: Core - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-05-07
**Phase:** 2-core
**Areas discussed:** (auto mode - all gray areas selected by default)

---

## INI Parsing Approach

| Option | Description | Selected |
|--------|-------------|----------|
| go-ini library | Standard Go INI parsing | ✓ |
| Manual parsing | Custom regex-based parsing | |
| Custom struct | ini struct tags | |

**User's choice:** [auto] go-ini library (recommended default)

---

## Password Algorithm

| Option | Description | Selected |
|--------|-------------|----------|
| Original XOR algorithm | Preserve existing code from muety/winscp2filezilla | ✓ |
| New implementation | Rewrite from scratch | |
| Third-party library | Use existing Go crypto packages | |

**User's choice:** [auto] Original XOR algorithm (recommended default)

---

## XML Structure

| Option | Description | Selected |
|--------|-------------|----------|
| Preserve folder hierarchy | Replicate WinSCP folder structure in FZ | ✓ |
| Flat list | All servers in root folder | |

**User's choice:** [auto] Preserve folder hierarchy (recommended default)

---

## Protocol Mapping

| Option | Description | Selected |
|--------|-------------|----------|
| FSProtocol 2=SFTP | Standard mapping | ✓ |
| Different mapping | Custom logic | |

**User's choice:** [auto] FSProtocol 2=SFTP (recommended default)

---

## the agent's Discretion

All decisions made via auto-select in auto mode — no user interaction required.