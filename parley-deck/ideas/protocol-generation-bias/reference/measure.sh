#!/bin/bash
# measure.sh — canonical, reproducible measurement of the quantities disputed in
# ideas/protocol-generation-bias/round-01.
#
# WHY THIS EXISTS
#   Three different numbers circulated in round-01 for the same quantities. Every
#   disagreement traced to an unstated definition, not to a counting error. This
#   script fixes each definition in code and prints the alternatives side by side,
#   so a reader can see which definition produces which number.
#
# TOOLING NOTE (load-bearing)
#   `grep`/`rg` on this machine resolve to ugrep/ripgrep, which honour .gitignore and
#   skip dot-files by default. A miss under those tools is NOT evidence of absence.
#   This script uses /usr/bin/grep (BSD grep, no ignore semantics) and `find` only.
#
# USAGE   bash measure.sh          # from anywhere; paths derive from script location
set -u

GREP=/usr/bin/grep
HERE=$(cd -- "$(dirname -- "$0")" && pwd)
IDEAS=$(cd -- "$HERE/../.." && pwd)
SELF_IDEA=protocol-generation-bias   # the idea this deliberation is about (confound)

if [ ! -d "$IDEAS" ]; then echo "FATAL: ideas dir not found at $IDEAS" >&2; exit 1; fi

hr()  { printf '%s\n' "------------------------------------------------------------------"; }
sec() { printf '\n'; hr; printf '%s\n' "$1"; hr; }

# frontmatter <file> : emit the YAML frontmatter block only.
# Definition: the lines strictly between the leading `---` on line 1 and the next
# line that is exactly `---`. A file whose line 1 is not `---` has NO frontmatter
# and emits nothing (so a prose mention can never be counted as a key).
frontmatter() {
  awk 'NR==1 { if ($0 != "---") exit; inf=1; next }
       inf && $0 == "---" { exit }
       inf { print }' "$1"
}

# sets_key <file> <key> : true if <key> is set as a top-level frontmatter key.
# Requires column 0 (so nested keys under `roles:` do not count) and a literal
# `key:` prefix (so a backticked prose mention of the key never counts).
sets_key() { frontmatter "$1" | $GREP -q "^$2:"; }

printf '%s\n' "=================================================================="
printf '%s\n' " CANONICAL MEASUREMENT REPORT — protocol-generation-bias"
printf '%s\n' "=================================================================="
printf 'ideas root : %s\n' "$IDEAS"
printf 'generated  : %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ') (UTC)"
printf 'grep binary: %s (%s)\n' "$GREP" "$($GREP --version 2>&1 | head -1)"

# ---------------------------------------------------------------------------
# Q3 — THE DENOMINATOR (answered first; everything else divides by it)
# ---------------------------------------------------------------------------
sec "Q3. DENOMINATOR — how many ideas are in ideas/ ?"

ALL_DIRS=$(find "$IDEAS" -mindepth 1 -maxdepth 1 -type d | sort)
N_DIRS=$(printf '%s\n' "$ALL_DIRS" | $GREP -c . )

NONDIR=$(find "$IDEAS" -mindepth 1 -maxdepth 1 ! -type d | wc -l | tr -d ' ')

WITH_PROMPT=""; WITHOUT_PROMPT=""
while IFS= read -r d; do
  [ -n "$d" ] || continue
  if [ -f "$d/00-prompt.md" ]; then WITH_PROMPT="$WITH_PROMPT$d
"; else WITHOUT_PROMPT="$WITHOUT_PROMPT$d
"; fi
done <<EOF
$ALL_DIRS
EOF

N_WITH=$(printf '%s\n' "$WITH_PROMPT" | $GREP -c . || true)
N_WITHOUT=$(printf '%s\n' "$WITHOUT_PROMPT" | $GREP -c . || true)

printf 'Definition A — "an idea is a directory":                     %s\n' "$N_DIRS"
printf 'Definition B — "a directory that has a 00-prompt.md":        %s\n' "$N_WITH"
printf 'Non-directory entries directly under ideas/:                 %s\n' "$NONDIR"
printf '\nDirectories WITHOUT a 00-prompt.md (the entire difference):\n'
if [ "$N_WITHOUT" -eq 0 ]; then printf '  (none)\n'; else
  printf '%s\n' "$WITHOUT_PROMPT" | while IFS= read -r d; do
    [ -n "$d" ] || continue
    printf '  %-40s contains: %s\n' "$(basename "$d")" "$(ls "$d" | tr '\n' ' ')"
  done
fi
printf '\nRULING: the denominator is %s for any question about a directory,\n' "$N_DIRS"
printf '        and %s for any question about a FRONTMATTER KEY (a key can only\n' "$N_WITH"
printf '        exist in a 00-prompt.md, so ideas lacking one are not eligible).\n'
printf '        Both "88" and "89" were right about different populations.\n'

# ---------------------------------------------------------------------------
# Q1 — require_model_diversity
# ---------------------------------------------------------------------------
sec "Q1. require_model_diversity — how many ideas SET it as a frontmatter key?"

FM_SET=""; PROSE_ONLY=""
while IFS= read -r d; do
  [ -n "$d" ] || continue
  f="$d/00-prompt.md"
  if sets_key "$f" require_model_diversity; then
    FM_SET="$FM_SET$(basename "$d")
"
  elif $GREP -q "require_model_diversity" "$f"; then
    PROSE_ONLY="$PROSE_ONLY$(basename "$d")
"
  fi
done <<EOF
$WITH_PROMPT
EOF

N_FM=$(printf '%s\n' "$FM_SET" | $GREP -c . || true)
N_PROSE=$(printf '%s\n' "$PROSE_ONLY" | $GREP -c . || true)
# NB: filter blank lines BEFORE the -v, otherwise the list terminator counts as a row.
N_FM_EXSELF=$(printf '%s\n' "$FM_SET" | $GREP . | $GREP -vc "^$SELF_IDEA$" || true)

printf 'Ideas SETTING the key in 00-prompt.md frontmatter:  %s of %s\n' "$N_FM" "$N_WITH"
printf '%s\n' "$FM_SET" | while IFS= read -r s; do [ -n "$s" ] || continue
  printf '  [KEY  ] %-45s created: %s\n' "$s" \
    "$(frontmatter "$IDEAS/$s/00-prompt.md" | $GREP '^created:' | head -1 | sed 's/^created: *//')"
done
printf '\nIdeas MENTIONING it in prose only (NOT a set key):  %s of %s\n' "$N_PROSE" "$N_WITH"
printf '%s\n' "$PROSE_ONLY" | while IFS= read -r s; do [ -n "$s" ] || continue
  printf '  [PROSE] %-45s %s\n' "$s" \
    "$($GREP -n 'require_model_diversity' "$IDEAS/$s/00-prompt.md" | head -1 | cut -c1-70)"
done

printf '\nCONFOUND — this idea (%s) sets the key and is the\n' "$SELF_IDEA"
printf 'subject of the measurement. Reported both ways, as required:\n'
printf '  INCLUDING this idea : %s of %s ideas set the key\n' "$N_FM" "$N_WITH"
printf '  EXCLUDING this idea : %s of %s ideas set the key\n' "$N_FM_EXSELF" "$((N_WITH - 1))"
printf '\nRECONCILIATION of the disputed numbers:\n'
printf '  "0 of 88"  = correct for the deck BEFORE this idea existed (excl. self).\n'
printf '  "2 of 88"  = key-set (%s) + prose-mention (%s); it counted a backticked\n' "$N_FM" "$N_PROSE"
printf '               prose reference as adoption. Adoption of a GATE means the\n'
printf '               key is SET; a sentence describing the flag gates nothing.\n'
printf '  CANONICAL  = %s of %s including self, %s of %s excluding self.\n' \
  "$N_FM" "$N_WITH" "$N_FM_EXSELF" "$((N_WITH - 1))"

# ---------------------------------------------------------------------------
# Q4 — the other opt-in gates
# ---------------------------------------------------------------------------
sec "Q4. Adoption of the other opt-in gates (frontmatter keys in 00-prompt.md)"

printf '%-28s %-12s %-12s %s\n' "KEY" "SET(n)" "PROSE-ONLY" "PERCENT of $N_WITH"
printf '%-28s %-12s %-12s %s\n' "---" "------" "----------" "-------------------"
for key in track checks strict_gate auto_implement require_model_diversity; do
  n=0; p=0
  while IFS= read -r d; do
    [ -n "$d" ] || continue
    f="$d/00-prompt.md"
    if sets_key "$f" "$key"; then n=$((n+1))
    elif $GREP -q "$key" "$f"; then p=$((p+1)); fi
  done <<EOF
$WITH_PROMPT
EOF
  pct=$(awk -v a="$n" -v b="$N_WITH" 'BEGIN{printf "%.1f%%", (b?100*a/b:0)}')
  printf '%-28s %-12s %-12s %s\n' "$key:" "$n" "$p" "$pct"
done
printf '\nNote: PROSE-ONLY counts ideas whose 00-prompt.md mentions the string\n'
printf 'anywhere but does NOT set it as a frontmatter key. It is shown to make\n'
printf 'the prose/key confound visible, and is NEVER added to the SET column.\n'

# Which values does track: actually take? (a set key with an unexpected value is
# still adoption, but the distribution is what makes the number interpretable)
printf '\ntrack: values actually set —\n'
{
  while IFS= read -r d; do
    [ -n "$d" ] || continue
    frontmatter "$d/00-prompt.md" | $GREP '^track:' | sed 's/^track: *//'
  done <<EOF
$WITH_PROMPT
EOF
} | sort | uniq -c | sed 's/^/  /'

# ---------------------------------------------------------------------------
# Q2 — "## Adversarial alternative"
# ---------------------------------------------------------------------------
sec "Q2. '## Adversarial alternative' — how many IDEAS carry the section?"

# Canonical artifacts = .md files not beginning with a dot. Agent transcripts
# (.codex.log, .hermes.log) and raw diffs (DIFF-fixups.txt) are process residue,
# not artifacts; they are counted separately, never in the headline number.
SECTION_FILES=$(find "$IDEAS" -type f -name '*.md' ! -name '.*' -print0 \
  | xargs -0 $GREP -l '^## Adversarial alternative' 2>/dev/null | sort)
N_SEC_FILES=$(printf '%s\n' "$SECTION_FILES" | $GREP -c . || true)

SECTION_IDEAS=$(printf '%s\n' "$SECTION_FILES" | sed "s|^$IDEAS/||" | cut -d/ -f1 | sort -u)
N_SEC_IDEAS=$(printf '%s\n' "$SECTION_IDEAS" | $GREP -c . || true)

printf 'Files containing a literal "## Adversarial alternative" heading: %s\n' "$N_SEC_FILES"
printf '%s\n' "$SECTION_FILES" | while IFS= read -r f; do [ -n "$f" ] || continue
  printf '  %s\n' "${f#"$IDEAS"/}"
done
printf '\nDistinct IDEAS carrying the section: %s of %s\n' "$N_SEC_IDEAS" "$N_DIRS"
printf '%s\n' "$SECTION_IDEAS" | while IFS= read -r s; do [ -n "$s" ] || continue
  printf '  %s\n' "$s"
done
pct=$(awk -v a="$N_SEC_IDEAS" -v b="$N_DIRS" 'BEGIN{printf "%.1f%%", (b?100*a/b:0)}')
printf '  => %s of %s ideas = %s\n' "$N_SEC_IDEAS" "$N_DIRS" "$pct"

# The over-count: any mention of the phrase, anywhere, in any file.
MENTION_FILES=$(find "$IDEAS" -type f -print0 \
  | xargs -0 $GREP -li 'adversarial alternative' 2>/dev/null | sort)
N_MEN_FILES=$(printf '%s\n' "$MENTION_FILES" | $GREP -c . || true)
MENTION_IDEAS=$(printf '%s\n' "$MENTION_FILES" | sed "s|^$IDEAS/||" | cut -d/ -f1 | sort -u)
N_MEN_IDEAS=$(printf '%s\n' "$MENTION_IDEAS" | $GREP -c . || true)

printf '\nCONTRAST — files MENTIONING the phrase anywhere (case-insensitive,\n'
printf 'including .log transcripts and prose quoting the protocol rule): %s files, %s ideas\n' \
  "$N_MEN_FILES" "$N_MEN_IDEAS"
printf '%s\n' "$MENTION_IDEAS" | while IFS= read -r s; do [ -n "$s" ] || continue
  printf '  %s\n' "$s"
done
printf '\nRECONCILIATION:\n'
printf '  "15 slugs" = the MENTION population counted by FILE (%s files, but only\n' "$N_MEN_FILES"
printf '               %s distinct ideas). It was a file count reported as slugs, and\n' "$N_MEN_IDEAS"
printf '               it includes files that merely QUOTE the protocol rule, name it\n'
printf '               in a finding, or are .log transcripts of an agent reading it.\n'
printf '  "5 of 89" / "~6" = intermediate mention-counts over .md files only; closer,\n'
printf '               but still counting mentions rather than section-bearing files.\n'
printf '  CANONICAL  = %s FILES carry the section, spanning %s IDEAS (%s of %s).\n' \
  "$N_SEC_FILES" "$N_SEC_IDEAS" "$N_SEC_IDEAS" "$N_DIRS"
printf '               Reported by idea, with the file count stated, as required.\n'

# ---------------------------------------------------------------------------
# Q5 — later-round reviewer rule compliance
# ---------------------------------------------------------------------------
sec "Q5. Later-round review compliance (review/round-NN, NN>=2)"

# Population: participant artifacts only.
#   - *.md under ideas/*/review/round-NN with NN >= 02
#   - EXCLUDING VOID.md, which self-declares "facilitator note, not a
#     participant artifact" and so cannot be held to a participant rule.
#   - EXCLUDING .log / .txt process residue by the *.md filter.
REVIEW_FILES=$(find "$IDEAS" -type f -name '*.md' ! -name '.*' ! -name 'VOID.md' \
  -path '*/review/round-*' \
  | $GREP -E '/review/round-(0[2-9]|[1-9][0-9])/[^/]+\.md$' | sort)
N_REV=$(printf '%s\n' "$REVIEW_FILES" | $GREP -c . || true)

N_VOID=$(find "$IDEAS" -type f -name 'VOID.md' -path '*/review/round-*' \
  | $GREP -cE '/review/round-(0[2-9]|[1-9][0-9])/' || true)

n_resp=0; n_head=0; n_head_other=0
MISSING_RESP=""; MISSING_HEAD=""
while IFS= read -r f; do
  [ -n "$f" ] || continue
  if frontmatter "$f" | $GREP -q '^responding-to:'; then
    n_resp=$((n_resp+1))
  else
    MISSING_RESP="$MISSING_RESP${f#"$IDEAS"/}
"
  fi
  if $GREP -q '^### @' "$f"; then
    n_head=$((n_head+1))
    # "other agent" = an @name that is not this file's own author (basename).
    self=$(basename "$f" .md)
    if $GREP -oh '^### @[A-Za-z0-9._-]*' "$f" | sed 's/^### @//' \
        | $GREP -qv "^$self$"; then
      n_head_other=$((n_head_other+1))
    fi
  else
    MISSING_HEAD="$MISSING_HEAD${f#"$IDEAS"/}
"
  fi
done <<EOF
$REVIEW_FILES
EOF

pr() { awk -v a="$1" -v b="$2" 'BEGIN{printf "%s / %s = %.1f%%\n", a, b, (b?100*a/b:0)}'; }

printf 'Denominator (participant .md artifacts in review/round-NN, NN>=2): %s\n' "$N_REV"
printf '  (VOID.md files excluded from the denominator: %s)\n' "$N_VOID"
printf '\ncarry a `responding-to:` frontmatter key : '; pr "$n_resp" "$N_REV"
printf 'contain any `### @` heading             : '; pr "$n_head" "$N_REV"
printf 'contain a `### @<OTHER agent>` heading  : '; pr "$n_head_other" "$N_REV"
printf '  ("other" = an @name different from the file'"'"'s own basename/author)\n'

# Guard against a false zero: does `responding-to:` ever appear OUTSIDE frontmatter?
# If it did, the frontmatter-only numerator would understate compliance.
n_resp_any=0
while IFS= read -r f; do
  [ -n "$f" ] || continue
  $GREP -q '^responding-to:' "$f" && n_resp_any=$((n_resp_any+1))
done <<EOF
$REVIEW_FILES
EOF
printf '\nGUARD — files with `^responding-to:` ANYWHERE in the file: %s\n' "$n_resp_any"
printf '        files with it inside frontmatter specifically  : %s\n' "$n_resp"
if [ "$n_resp_any" -eq "$n_resp" ]; then
  printf '        => identical, so the frontmatter restriction costs nothing here.\n'
else
  printf '        => DIFFER by %s: some files declare it outside frontmatter.\n' \
    "$((n_resp_any - n_resp))"
fi

printf '\nSample of files MISSING responding-to: (%s total, first 12)\n' "$((N_REV - n_resp))"
printf '%s\n' "$MISSING_RESP" | $GREP . | head -12 | sed 's/^/  /'

printf '\nSample of files MISSING any "### @" heading (%s total, first 12)\n' "$((N_REV - n_head))"
printf '%s\n' "$MISSING_HEAD" | $GREP . | head -12 | sed 's/^/  /'

# ---------------------------------------------------------------------------
sec "SUMMARY — one line per disputed quantity"
printf 'Q3 denominator      : %s idea directories; %s have a 00-prompt.md\n' "$N_DIRS" "$N_WITH"
printf 'Q1 require_model_div: %s/%s set the key (incl. self) | %s/%s (excl. self) | %s prose-only\n' \
  "$N_FM" "$N_WITH" "$N_FM_EXSELF" "$((N_WITH - 1))" "$N_PROSE"
printf 'Q2 adversarial alt  : %s ideas / %s files carry the section (of %s ideas)\n' \
  "$N_SEC_IDEAS" "$N_SEC_FILES" "$N_DIRS"
printf 'Q5 responding-to:   : %s/%s later-round review artifacts\n' "$n_resp" "$N_REV"
printf 'Q5 "### @other"     : %s/%s later-round review artifacts\n' "$n_head_other" "$N_REV"
hr
