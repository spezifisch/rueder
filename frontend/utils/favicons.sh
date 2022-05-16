#!/bin/sh -e

INPUT="./public/favicon.svg"
OUTPUT="./public"

which inkscape > /dev/null || ( echo "inkscape not in PATH"; exit 1 )

[ -e "$INPUT" ]  || ( echo "favicon not found at $INPUT"; exit 1 )
[ -d "$OUTPUT" ] || ( echo "output dir not found at $OUTPUT"; exit 1 )

inkscape --export-png "$OUTPUT/favicon.png" "$INPUT"

for w in 16 32 150 180; do
    inkscape --export-png "$OUTPUT/favicon-$w.png" -w $w "$INPUT"
done
