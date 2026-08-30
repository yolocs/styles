# Terminal Theme Format v1

This document defines `yolodev-terminal-theme` version 1, the canonical,
terminal-neutral format used in this repository. A canonical file is TOML and
contains metadata, six semantic terminal colors, and the normal and bright
ANSI palettes.

## Required Shape

Use [minimal.toml](examples/minimal.toml) as the complete example. Every key in
that file is required. Unknown and duplicate keys are errors.

- `format` is exactly `yolodev-terminal-theme`.
- `version` is the integer `1`.
- `theme.name`, `theme.variant`, and `theme.description` are non-empty strings.
- `theme.appearance` is `dark` or `light`.
- Every value under `colors`, `palette.normal`, and `palette.bright` is an sRGB
  color written as `#RRGGBB`.

Writers normalize hexadecimal digits to uppercase and emit fields in schema
order. Writers do not preserve comments. Readers reject format versions they
do not understand.

## Semantic Colors

- `background`: default terminal cell background.
- `foreground`: default terminal text.
- `cursor`: cursor surface.
- `cursor_text`: text underneath the cursor.
- `selection_background`: selected-cell background.
- `selection_foreground`: selected text.

## ANSI Palettes

Both `palette.normal` and `palette.bright` contain `black`, `red`, `green`,
`yellow`, `blue`, `magenta`, `cyan`, and `white`. Exporters map the normal
colors to ANSI indices 0–7 and bright colors to indices 8–15 in that order.

## Validation

Structural violations are errors and block save or export. The validator also
emits non-blocking quality warnings when WCAG relative-luminance contrast is:

- below 4.5:1 for foreground against background;
- below 3.0:1 for selection foreground against selection background; or
- below 3.0:1 for cursor against background.

Changing required fields or their meaning requires another format version.
Adding terminal-specific configuration does not: such configuration belongs
in an exporter rather than this format.
