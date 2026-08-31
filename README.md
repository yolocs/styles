# styles

Personal style assets and the tools that shape them: terminal themes today,
with room for fonts, prompts, and setup conventions when concrete needs arise.
Each style family lives outside the application that edits it.

## Repository Layout

```text
.agents/skills/       Repository-local AI workflows
standards/            App-independent file formats
themes/               Canonical terminal-theme assets
yolodev/              Go application for creating and editing styles
```

The first theme family is `themes/yolodev/`, including the Signal Grid dark
variant.

## Quick Start

Requirements: Go 1.26 or newer and a modern true-color terminal. Ghostty is the
primary target.

```sh
make build
make run
```

The editor supports mouse selection, click-and-drag saturation/value and hue
controls, direct hex input, canonical TOML import/save, and Ghostty or Codex
export. The CLI supports the same two formats.

Validate and export without opening the TUI:

```sh
./bin/yolodev theme validate themes/yolodev/placeholder.toml
./bin/yolodev theme export --format ghostty themes/yolodev/placeholder.toml
./bin/yolodev theme export --format ghostty --output /tmp/yolodev-ghostty themes/yolodev/placeholder.toml
./bin/yolodev theme export --format codex --output /tmp/signal-grid.tmTheme themes/yolodev/signal-grid.toml
```

Non-interactive export refuses to overwrite an existing destination. The app
never edits or installs into your Ghostty or Codex configuration.

To use a Codex export, place the `.tmTheme` file under `$CODEX_HOME/themes`
(normally `~/.codex/themes`) and select it with `/theme`. This changes syntax
highlighting in code; the surrounding Codex interface still follows its own
semantic colors and the terminal background.

## Canonical Format

[`standards/terminal-theme-v1.md`](standards/terminal-theme-v1.md) defines the
strict TOML contract. Canonical themes are terminal-neutral; exporter-specific
syntax stays in the app.

## AI-assisted Themes

The repository-local `$yolodev-theme` skill creates canonical themes from a
text prompt or attached image, validates them, and explains any intentional
contrast warnings. For example:

```text
Use $yolodev-theme to create a dark terminal palette inspired by this rainy
neon street photo. Save it as themes/yolodev/rainy-neon.toml.
```

## Development

```sh
make fmt-check
make test
make vet
```

See [`AGENTS.md`](AGENTS.md) for durable repository decisions and boundaries.
