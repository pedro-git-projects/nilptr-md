#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

CHARFILE="subset-chars.txt"

for SRC in MononokiNerdFont-*.ttf MononokiNerdFontMono-*.ttf; do
  BASENAME="${SRC%.ttf}"
  OUT="${BASENAME}.subset.woff2"

  echo "▶ Subsetting $SRC → $OUT"
  pyftsubset "$SRC" \
    --text-file="$CHARFILE"    \
    --flavor=woff2             \
    --output-file="$OUT"       \
    --layout-features="*"      \
    --glyph-names              \
    --symbol-cmap              \
    --drop-tables=PfEd         \
    --no-hinting               \
    --desubroutinize

  echo "✔ Generated $OUT"
done

