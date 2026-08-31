# yolodev

`yolodev` is the repository's Go workbench for personal style assets. The MVP
implements terminal themes through a noun-based command surface:

```text
yolodev theme edit [FILE]
yolodev theme validate FILE
yolodev theme export --format FORMAT [--output PATH] FILE
```

Ghostty and Codex TextMate themes are registered export formats. The CLI is
built with Cobra, and renderers are dispatched through an internal registry so
formats can be added without changing command parsing.

For example, export Signal Grid for Codex without installing it automatically:

```sh
yolodev theme export --format codex --output signal-grid.tmTheme ../themes/yolodev/signal-grid.toml
```

Move the result under `$CODEX_HOME/themes` and select it with `/theme`. The
export controls syntax highlighting, not the surrounding Codex interface.

From the repository root, use `make build` to create `bin/yolodev` and `make
run` to edit the placeholder theme.

## Editor Controls

- Click semantic rows or ANSI swatches to select a color.
- Click and drag the saturation/value plane or hue strip.
- Click the hex field, or press Tab, to enter an exact `#RRGGBB` value.
- Click Import, Save TOML, or Export for file operations.
- In the export popup, choose Ghostty or Codex with Left/Right or the mouse;
  press Tab to switch between the format selector and path.
- Press Ctrl+S to save and Ctrl+Q to quit.
- Press Escape to cancel the active edit or dialog.

The full layout requires at least 100 columns by 28 rows. Existing export
targets and unsaved edits require explicit confirmation in the editor.
