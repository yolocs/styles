# Agent Instructions

## Repository Purpose

This repository stores and develops personal style assets: terminal themes,
fonts, prompts, and related setup. Each style family belongs in its own
subdirectory. `yolodev` is the first-class application for creating and
editing these assets; terminal color themes are its first domain, not its
permanent limit.

Keep the current task narrow. Do not add speculative abstractions for future
style domains until a concrete domain needs them.

## Current Layout

- `yolodev/` contains the Go application and its tests.
- `themes/` contains terminal-theme assets outside the application.
- `standards/` contains app-independent format specifications and examples.
- `.agents/skills/` contains repository-local AI skills.

## Durable Decisions

- Canonical terminal themes use the versioned TOML format documented in
  `standards/terminal-theme-v1.md`.
- Theme variants are separate files. The initial placeholder lives under
  `themes/yolodev/`.
- Canonical theme data must remain terminal-neutral. Terminal-specific syntax
  belongs in exporters.
- `yolodev` uses noun-based commands such as `yolodev theme edit` so another
  style domain can be added without redesigning existing commands.
- The terminal-theme editor uses Bubble Tea v2 and Lip Gloss v2. Mouse input,
  including click-and-drag color picking, is an MVP requirement.
- Ghostty export is deterministic and limited to known color options. Never
  edit or install into the user's Ghostty configuration automatically.
- Existing files must not be overwritten without explicit confirmation.
- Writes to theme and export files must be atomic.

## Go Conventions

- Use Go 1.26 or newer.
- Prefer explicit, small packages and standard-library code over framework
  machinery.
- Keep the canonical theme model independent from the TUI and exporters.
- Keep the Ghostty renderer pure: validated theme in, bytes out.
- Model TUI interactions as testable state transitions; terminal rendering and
  filesystem effects stay at the edges.
- Return contextual errors. Do not panic for invalid input or ordinary I/O
  failures.

## Verification

Run Go commands from the repository root with Go's `-C` flag:

```sh
go -C yolodev test ./...
go -C yolodev vet ./...
```

Also validate the checked-in placeholder and compare Ghostty export against
its golden test. Mouse behavior requires a manual smoke test in a modern
terminal; Ghostty is the primary target.

## Decision Hygiene

Update this file when a durable repository-wide rule or boundary changes.
Put detailed rationale and feature behavior in a design specification instead
of turning this file into a changelog.
