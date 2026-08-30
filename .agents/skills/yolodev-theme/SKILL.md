---
name: yolodev-theme
description: Use when creating a canonical yolodev terminal color theme from a text prompt, mood, palette, screenshot, photograph, illustration, or attached image.
---

# Generate a Yolodev Terminal Theme

Create a complete terminal-theme TOML file whose colors express the source
idea while remaining usable in real terminal content. The image or prompt is
inspiration; terminal roles determine the final assignments.

## Required Workflow

1. Read `standards/terminal-theme-v1.md` from the repository root. Treat it as
   authoritative; do not invent keys or omit required fields.
2. Inspect an attached image with the available image-viewing capability. Do
   not invoke image generation. Identify useful anchors by hue, luminance, and
   emotional weight instead of mapping dominant pixels directly to ANSI slots.
3. Choose the default background and foreground first. Then choose cursor and
   selection pairs that remain visible. Build normal and bright ANSI palettes
   with recognizable red, green, yellow, blue, magenta, and cyan roles.
4. Write one complete variant. Use the path requested by the user, or default
   to `themes/yolodev/<variant>.toml` for a yolodev variant.
5. Inspect the destination before writing. Never replace an existing file
   unless the user explicitly requested replacement.
6. Validate with the repository application. From the repository root, a
   root-relative file uses:

   ```sh
   go -C yolodev run ./cmd/yolodev theme validate ../themes/yolodev/<variant>.toml
   ```

   Use an explicit absolute path when the destination is elsewhere.
7. Correct every error and validate again. Warnings may remain when they are an
   intentional aesthetic tradeoff; report the affected fields and ratios.

## Palette Decisions

- Use exact sRGB `#RRGGBB`; the writer will normalize letter case.
- Preserve the source's mood in neutrals and accents, not at the expense of
  making red, green, yellow, and blue unrecognizable in common CLI output.
- Make each bright color visibly related to its normal counterpart. A light
  theme may use stronger chroma or darker values rather than blindly lightening.
- Avoid identical foreground/background, cursor/background, and selection
  pairs. Treat validator contrast warnings as design feedback, not schema
  failures.

## Handoff

State the created path, summarize the visual choices in a few sentences, and
include validation results. Suggest `yolodev theme edit FILE` when the user
wants to tune the result interactively, or `yolodev theme export ghostty FILE`
when they want Ghostty output.

## Common Mistakes

- Returning a palette in prose instead of writing the complete TOML artifact.
- Copying image colors literally without assigning terminal semantics.
- Silently overwriting an existing variant.
- Calling a low-contrast warning a successful accessibility check.
- Adding Ghostty keys to canonical TOML instead of using the exporter.
