#!/bin/zsh
# claude-1 round-04 independent rerun of original r3 probes (unchanged sources) on shared volume and /tmp.
set -u
W="/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude"
cd "$W"
export PATH="$W/.parley-runtime/quota-implementation/fixup-1/host/no-provider-bin:$PATH"
OUT="$W/.parley-runtime/claude1-r4/probes"
BIN=/tmp/claude1-r4-bin; mkdir -p $BIN
for p in create knoboff reinclude lease zadv; do
  go build -o $BIN/$p .parley-runtime/claude1-r3/$p/main.go 2>&1 | sed "s/^/[build $p] /"
done
for vol in shared local; do
  if [ $vol = shared ]; then B="$OUT/$vol-base"; else B="/tmp/claude1-r4-$vol-base"; fi
  rm -rf "$B"; mkdir -p "$B/create" "$B/tmp"
  {
    echo "### volume=$vol base=$B date=$(date -u +%FT%TZ)"
    echo "\$ create"; $BIN/create "$B/create"; echo "exit=$?"
    echo "\$ TMPDIR=$B/tmp knoboff"; TMPDIR="$B/tmp" $BIN/knoboff; echo "exit=$?"
    echo "\$ TMPDIR=$B/tmp reinclude"; TMPDIR="$B/tmp" $BIN/reinclude; echo "exit=$?"
    echo "\$ lease setup/hold/try"; $BIN/lease setup "$B/lease"; echo "exit=$?"
    ($BIN/lease hold "$B/lease" > "$B/hold.out" 2>&1 &) ; sleep 1.5
    $BIN/lease try "$B/lease"; echo "try exit=$?"
    sleep 6; cat "$B/hold.out"; $BIN/lease try "$B/lease"; echo "post-release try exit=$?"
    ls -la "$B/lease/parley-deck/ideas/demo" ; find "$B/lease" -name '*.lease*' -o -name '*lock*' | head
  } > "$OUT/$vol.log" 2>&1
done
$BIN/zadv > "$OUT/zadv.log" 2>&1; echo "zadv exit=$?" >> "$OUT/zadv.log"
echo done > "$OUT/DONE"
