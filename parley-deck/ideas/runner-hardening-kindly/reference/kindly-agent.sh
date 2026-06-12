#!/usr/bin/env bash
# Cross-model review and consult via the opposite agent CLI. See
# references/audit.md for the gate protocol and references/reviewers.md for the
# verified CLI mechanics, launch guidance, and snapshot behavior.
set -euo pipefail
trap '' PIPE
# Optional index-refresh locks would write the audited repo's .git during
# status and diff probes.
export GIT_OPTIONAL_LOCKS=0

usage() {
  local code="${1:-2}"
  local usage_fd=2
  [ "$code" = 0 ] && usage_fd=1
  cat >&"$usage_fd" <<'EOF' || true
Usage: agent.sh --slug <slug> [options] [BRIEF|-]

Options:
  --slug <slug>        Unit identifier; required. Recorded in the report header and filename.
  --mode <mode>        audit | consult  (default: audit). audit writes its report to
                       <repo>/.local/audits/; consult writes to .local/consults/.
  --scope <scope>      audit only: uncommitted | commit:<sha> | base:<ref>  (default: uncommitted)
  --resume <session>   Follow-up pass on an existing reviewer session (audit: verify findings
                       closed; a resumed pass converges a gate but never closes it)
  --reviewer <cli>     codex | claude (default: the opposite host agent when
                       CODEX_*, CLAUDE* or AI_AGENT markers are visible; otherwise required)
  --exec               audit only: behavioral review — the reviewer may run tests and start
                       local infrastructure (codex: workspace-write sandbox; claude: Bash inside
                       an OS sandbox — repo and temp dirs writable, network constrained).
                       Runs in the live tree (tests need ignored build artifacts); tracked and
                       untracked content is hashed and any mutation is flagged loudly.
  --repo <dir>         Target repository (default: git toplevel of the current directory;
                       consult also accepts a plain non-git directory)
  --model <model>      Pin a specific top-tier reviewer model (default: the CLI's default model)
  BRIEF                Trailing argument(s): the consult question, or audit intent, constraints,
                       validation output and rebuttals. Use '-' to read from stdin.
                       Required for consult.

Read-only audits review a disposable snapshot checkout — uncommitted work is captured as a
snapshot commit, committed scopes check out HEAD — so the live tree stays free for the caller
while the review runs. Every run appends a provenance line to the report directory's index.jsonl.
Prints the report path and the reviewer session id; exits nonzero only when no usable report is written.
Full-depth reviewer runs are intentionally slow; supervise them until the report line lands.
EOF
  exit "$code"
}

fail() {
  echo "agent.sh: $1" >&2 || true
  exit 1
}

# Signal a process and all its descendants (pgrep -P works on macOS and Linux);
# a bare kill can orphan the reviewer CLI's children and leak a running review.
signal_tree() {
  local sig="$1" pid="$2" child
  kill "-${sig}" -- "-$pid" 2>/dev/null || true
  for child in $(pgrep -P "$pid" 2>/dev/null); do
    signal_tree "$sig" "$child"
  done
  kill "-${sig}" "$pid" 2>/dev/null || true
}

tree_pids() {
  local pid="$1" child
  printf '%s\n' "$pid"
  for child in $(pgrep -P "$pid" 2>/dev/null); do
    tree_pids "$child"
  done
}

# TERM the tree, give every member up to five seconds to exit cleanly, then
# KILL the stragglers. Membership is collected pre-TERM so a root that exits
# first still leaves its children the grace window; zombies pass kill -0, so
# the state probe keeps them from burning the window.
kill_tree() {
  local pid="$1" waited=0 alive member members
  members="$(tree_pids "$pid")"
  signal_tree TERM "$pid"
  # Reparented members drop out of the live-tree walk; signal recorded pids too.
  for member in $members; do
    kill -TERM "$member" 2>/dev/null || true
  done
  while [ "$waited" -lt 50 ]; do
    alive=0
    for member in $members; do
      kill -0 "$member" 2>/dev/null || continue
      case "$(ps -o state= -p "$member" 2>/dev/null || true)" in (*Z*) continue ;; esac
      alive=1
      break
    done
    [ "$alive" = 1 ] || break
    sleep 0.1
    waited=$((waited + 1))
  done
  signal_tree KILL "$pid"
  for member in $members; do
    kill -KILL "$member" 2>/dev/null || true
  done
}

exec_isolated() {
  if command -v setsid >/dev/null 2>&1; then
    exec setsid "$@"
  elif command -v perl >/dev/null 2>&1; then
    exec perl -MPOSIX=setsid -e 'setsid() or die "setsid: $!"; exec @ARGV or die "exec: $!"' -- "$@"
  else
    exec "$@"
  fi
}

# Timing seams: env overrides exist for the test harness; production keeps the defaults.
positive_int_env() {
  local value="${!1:-}"
  case "$value" in
    (''|*[!0-9]*|0) printf '%s' "$2" ;;
    (*) printf '%s' "$value" ;;
  esac
}

event_metric() {
  local flag="$1" file="$2"
  if [ -e "$file" ]; then
    wc "$flag" < "$file" | tr -d '[:space:]'
  else
    printf '0'
  fi
}

event_summary() {
  local summary
  [ -s "$1" ] || return 0
  summary="$(
    {
      tail -100 "$1" 2>/dev/null \
        | jq -Rr 'fromjson? | select(type == "object" and (.type? != null)) | if .subtype? != null then "\(.type)/\(.subtype)" else .type end' 2>/dev/null \
        | tail -1
    } || true
  )"
  printf '%s' "$summary"
}

event_final_result_subtype() {
  { jq -Rr 'fromjson? | select(type == "object" and .type == "result") | .subtype // empty' "$1" 2>/dev/null || true; } | tail -1
}

event_final_result_failed() {
  { jq -Rr '
    def has_errors:
      ((.errors? | type) == "array" and (.errors | length > 0)) or
      ((.errors? | type) == "string" and (.errors | length > 0));
    fromjson? |
    select(type == "object" and .type == "result") |
    select(
      (.is_error? == true) or
      (.api_error_status? != null) or
      (.apiErrorStatus? != null) or
      (.error_type? != null) or
      (.error? != null) or
      has_errors or
      ((.subtype? | type) == "string" and (.subtype | test("(^|_)error")))
    ) | "yes"
  ' "$1" 2>/dev/null || true; } | tail -1
}

event_metadata_tail() {
  local line parsed
  tail -50 "$1" 2>/dev/null | while IFS= read -r line; do
    parsed="$(printf '%s\n' "$line" | jq -Rr '
      (fromjson? // null) as $event |
      if $event == null then "unparseable event"
      elif ($event | type) == "object" then
        "type=\($event.type // "unknown")" +
        (if $event.subtype? != null then " subtype=\($event.subtype)" else "" end)
      else "non-object event"
      end
    ' 2>/dev/null || true)"
    printf '%s\n' "${parsed:-unparseable event}"
  done
}

event_failure_details() {
  {
    jq -Rr '
      def has_errors:
        ((.errors? | type) == "array" and (.errors | length > 0)) or
        ((.errors? | type) == "string" and (.errors | length > 0));
      def text(v):
        if v == null then ""
        elif (v | type) == "string" then v
        else (v | tostring)
        end;
      fromjson? |
      select(type == "object") |
      select(
        (.type? == "error") or
        (.type? == "turn.failed") or
        (.type? == "rate_limit_event") or
        (.level? == "error") or
        (.is_error? == true) or
        (.api_error_status? != null) or
        (.apiErrorStatus? != null) or
        (.error_type? != null) or
        (.error? != null) or
        has_errors or
        ((.subtype? | type) == "string" and (.subtype | test("(^|_)error")))
      ) |
      [
        "type=\(.type // "unknown")",
        (if .subtype? != null then "subtype=\(.subtype)" else empty end),
        (if .code? != null then "code=\(.code)" else empty end),
        (if ((.error? | type) == "object" and .error.code? != null) then "code=\(.error.code)" else empty end),
        (if .status? != null then "status=\(.status)" else empty end),
        (if .api_error_status? != null then "api_error_status=\(.api_error_status)" else empty end),
        (if .apiErrorStatus? != null then "apiErrorStatus=\(.apiErrorStatus)" else empty end),
        (if .error_type? != null then "error_type=\(.error_type)" else empty end),
        (text(.message?) as $message | if $message != "" then "message=\($message)" else empty end),
        (text(.error?) as $error | if $error != "" then "error=\($error)" else empty end),
        (if ((.error? | type) == "object" and .error.message? != null) then "message=\(.error.message)" else empty end),
        (text(.details?) as $details | if $details != "" then "details=\($details)" else empty end),
        (text(.errors?) as $errors | if $errors != "" then "errors=\($errors)" else empty end),
        (text(.result?) as $result | if $result != "" then "result=\($result)" else empty end)
      ] | join(" ")
    ' "$1" 2>/dev/null
    grep -E '^ERROR[: ]' "$1" 2>/dev/null || true
  } | tail -20
}

stderr_failure_details() {
  local pattern
  pattern='rate[-_ ]?limit|rate_limit_error|usageLimitExceeded|usage limit|session limit|too many requests|(^|[^0-9])429([^0-9]|$)|overloaded|overloaded_error|serverOverloaded|server_error|server error|internalServerError|api_error|api error|timeout_error|request timed out|timed out|deadline exceeded|(status|http|code|response|error)[^0-9]{0,40}5[0-9][0-9]|5[0-9][0-9][^0-9]{0,40}(status|http|code|response|error)|authentication_failed|authentication_error|authentication (failed|error|required)|unauthorized|forbidden|permission_error|oauth_org_not_allowed|oauth (org|error|failed|not allowed)|billing_error|billing (error|failed|required)|invalid_request|invalid_request_error|invalid request|badRequest|bad request|model_not_found|model not found|not_found_error|request_too_large|max_output_tokens|max output tokens|contextWindowExceeded|context length|cyberPolicy|sandboxError|Exceeded USD budget|exceeded[[:space:][:alpha:]]*budget'
  grep -Eai "$pattern" "$1" 2>/dev/null | tail -20 | sed 's/^/stderr: /' || true
}

event_recovery_hint() {
  if jq -Re '
    fromjson? |
    select(type == "object") |
    select(
      (.type? == "rate_limit_event") or
      (.code? == "rate_limit") or
      ((.error? | type) == "object" and .error.code? == "rate_limit") or
      (.error? == "rate_limit") or
      ((.error? | type) == "object" and .error.type? == "rate_limit_error") or
      (.error? == "usageLimitExceeded") or
      (.error_type? == "rate_limit") or
      (.api_error_status? == 429) or
      (.apiErrorStatus? == 429) or
      ((.message? | type) == "string" and (.message | test("rate.?limit|usage limit|session limit|too many requests"; "i"))) or
      ((.result? | type) == "string" and (.result | test("rate.?limit|usage limit|session limit|too many requests"; "i")))
    )
  ' "$1" >/dev/null 2>&1; then
    echo "agent.sh: reviewer appears rate-limited; retry after the provider reset or rerun with another explicit --reviewer"
  elif jq -Re '
    fromjson? |
    select(type == "object") |
    select(
      (.error_type? == "overloaded") or
      (.error_type? == "server_error") or
      (.error? == "serverOverloaded") or
      (.error? == "internalServerError") or
      ((.error? | type) == "object" and .error.type? == "overloaded_error") or
      ((.error? | type) == "object" and .error.type? == "api_error") or
      ((.error? | type) == "object" and .error.type? == "timeout_error") or
      (.api_error_status? == 529) or
      (.api_error_status? == 504) or
      (.apiErrorStatus? == 529) or
      (.apiErrorStatus? == 504)
    )
  ' "$1" >/dev/null 2>&1; then
    echo "agent.sh: reviewer provider appears overloaded or unavailable; retry after the service recovers or choose another explicit --reviewer"
  elif jq -Re '
    fromjson? |
    select(type == "object") |
    select((.api_error_status? != null) or (.apiErrorStatus? != null) or (.error? != null) or (.is_error? == true) or (.error_type? != null))
  ' "$1" >/dev/null 2>&1; then
    echo "agent.sh: reviewer reported an API error; inspect the saved failed log and retry or choose another explicit --reviewer"
  fi
}

stderr_recovery_hint() {
  if grep -Eiq 'rate[-_ ]?limit|rate_limit_error|usageLimitExceeded|usage limit|session limit|too many requests|(^|[^0-9])429([^0-9]|$)' "$1" 2>/dev/null; then
    echo "agent.sh: reviewer appears rate-limited; retry after the provider reset or rerun with another explicit --reviewer"
  elif grep -Eiq 'overloaded|overloaded_error|serverOverloaded|server_error|server error|internalServerError|api_error|api error|timeout_error|request timed out|timed out|deadline exceeded|(status|http|code|response|error)[^0-9]{0,40}5[0-9][0-9]|5[0-9][0-9][^0-9]{0,40}(status|http|code|response|error)' "$1" 2>/dev/null; then
    echo "agent.sh: reviewer provider appears overloaded or unavailable; retry after the service recovers or choose another explicit --reviewer"
  elif grep -Eiq 'authentication_failed|authentication_error|authentication (failed|error|required)|unauthorized|forbidden|permission_error|oauth_org_not_allowed|oauth (org|error|failed|not allowed)|billing_error|billing (error|failed|required)|invalid_request|invalid_request_error|invalid request|badRequest|bad request|model_not_found|model not found|not_found_error|request_too_large|max_output_tokens|max output tokens|contextWindowExceeded|context length|cyberPolicy|sandboxError|Exceeded USD budget|exceeded[[:space:][:alpha:]]*budget' "$1" 2>/dev/null; then
    echo "agent.sh: reviewer reported a provider or configuration error; inspect the saved failed log before retrying"
  fi
}

event_error_summary() {
  { jq -Rr '
    def has_errors:
      ((.errors? | type) == "array" and (.errors | length > 0)) or
      ((.errors? | type) == "string" and (.errors | length > 0));
    fromjson? |
    select(type == "object") |
    select(
      (.type? == "error") or
      (.type? == "turn.failed") or
      (.type? == "rate_limit_event") or
      (.level? == "error") or
      (.is_error? == true) or
      (.api_error_status? != null) or
      (.apiErrorStatus? != null) or
      (.error_type? != null) or
      (.error? != null) or
      has_errors or
      ((.subtype? | type) == "string" and (.subtype | test("(^|_)error")))
    ) |
    "agent.sh: reviewer emitted failure event: type=\(.type // "unknown")" +
    (if .subtype? != null then " subtype=\(.subtype)" else "" end) +
    (if .api_error_status? != null then " api_error_status=\(.api_error_status)" else "" end) +
    (if .apiErrorStatus? != null then " apiErrorStatus=\(.apiErrorStatus)" else "" end) +
    (if .error_type? != null then " error_type=\(.error_type)" else "" end) +
    (if (.error? | type) == "string" then " error=\(.error)" else "" end) +
    (if .code? != null then " code=\(.code)" else "" end) +
    (if ((.error? | type) == "object" and .error.code? != null) then " code=\(.error.code)" else "" end)
  ' "$1" 2>/dev/null || true; } | tail -1
}

# Supervises a running reviewer: content-free heartbeats on stderr, plus a stall
# guard — a healthy deep review keeps emitting events, so a long silence after the
# first event means a wedged process (denial loop, dead transport), not depth.
stalled=0
wait_with_progress() {
  local pid="$1" event_file="$2" label="$3"
  local started="$4" last_notice now elapsed lines bytes summary detail
  local heartbeat_secs stall_secs last_bytes=0 last_change
  heartbeat_secs="$(positive_int_env KINDLY_HEARTBEAT_SECS 30)"
  stall_secs="$(positive_int_env KINDLY_STALL_SECS 1800)"
  last_notice="$started"
  last_change="$started"
  while kill -0 "$pid" 2>/dev/null; do
    sleep 1
    kill -0 "$pid" 2>/dev/null || break
    now="$(date +%s)"
    bytes="$(event_metric -c "$event_file")"
    if [ "$bytes" != "$last_bytes" ]; then
      last_bytes="$bytes"
      last_change="$now"
    elif [ $((now - last_change)) -ge "$stall_secs" ]; then
      stalled=1
      echo "agent.sh: no new ${label} reviewer events for $((now - last_change))s; stopping the stalled run" >&2 || true
      echo "stalled: no new ${label} reviewer events for $((now - last_change))s" >> "$errlog" 2>/dev/null || true
      kill_tree "$pid"
      break
    fi
    if [ $((now - last_notice)) -ge "$heartbeat_secs" ]; then
      elapsed=$((now - started))
      lines="$(event_metric -l "$event_file")"
      summary="$(event_summary "$event_file")"
      detail="${summary:+, last event: ${summary}}"
      echo "agent.sh: ${label} reviewer still running (${elapsed}s elapsed, ${lines} events, ${bytes} bytes${detail})" >&2 || true
      last_notice="$now"
    fi
  done
  wait "$pid"
}

slug="" mode="audit" scope="uncommitted" resume="" reviewer="" repo="" model="" brief="" exec_mode=0

while [ $# -gt 0 ]; do
  case "$1" in
    --slug) slug="${2:?--slug needs a value}"; shift 2 ;;
    --mode) mode="${2:?--mode needs a value}"; shift 2 ;;
    --scope) scope="${2:?--scope needs a value}"; shift 2 ;;
    --resume) resume="${2:?--resume needs a value}"; shift 2 ;;
    --reviewer) reviewer="${2:?--reviewer needs a value}"; shift 2 ;;
    --exec) exec_mode=1; shift ;;
    --repo) repo="${2:?--repo needs a value}"; shift 2 ;;
    --model) model="${2:?--model needs a value}"; shift 2 ;;
    --effort) fail "reasoning depth is pinned to xhigh for both reviewers" ;;
    -h|--help) usage 0 ;;
    -) brief="${brief:+${brief} }$(cat)"; shift ;;
    --*) fail "unknown option: $1" ;;
    *) brief="${brief:+${brief} }$1"; shift ;;
  esac
done

# Codex markers win over CLAUDECODE: nested processes inherit the outer env.
agent_marker="$(printf '%s' "${AI_AGENT:-}" | tr '[:upper:]' '[:lower:]')"
detected_host=""
if [ -n "${CODEX_THREAD_ID:-}" ] || [ -n "${CODEX_SANDBOX:-}" ] || [[ "$agent_marker" == codex* ]]; then
  detected_host="codex"
elif [ -n "${CLAUDECODE:-}" ] || [ -n "${CLAUDE_CODE_SESSION_ID:-}" ] || [ -n "${CLAUDE_CODE_ENTRYPOINT:-}" ] || [[ "$agent_marker" == claude* ]]; then
  detected_host="claude"
fi
if [ -z "$reviewer" ]; then
  case "$detected_host" in
    codex) reviewer="claude" ;;
    claude) reviewer="codex" ;;
    *) fail "could not detect host agent; pass --reviewer claude from Codex or --reviewer codex from Claude" ;;
  esac
fi
same_reviewer=0
if [ -n "$detected_host" ] && [ "$reviewer" = "$detected_host" ]; then
  same_reviewer=1
  echo "agent.sh: warning: reviewer matches the detected host agent (${detected_host}); pass the opposite reviewer for a cross-model check unless this pairing is deliberate" >&2 || true
fi

[ -n "$slug" ] || usage
case "$slug" in (*[!a-zA-Z0-9._-]*) fail "slug must be filename-safe: $slug" ;; esac
case "$mode" in (audit|consult) ;; (*) fail "unknown mode: $mode" ;; esac
case "$reviewer" in (codex|claude) ;; (*) fail "unknown reviewer: $reviewer" ;; esac
command -v jq >/dev/null || fail "jq is required to parse reviewer output (brew install jq)"
command -v "$reviewer" >/dev/null 2>&1 || fail "reviewer CLI not found on PATH: $reviewer"
if [ "$exec_mode" = 1 ] && [ "$mode" = "consult" ]; then
  fail "--exec runs tests and infrastructure for an audit; a consult is advisory and has nothing to run"
fi
if [ "$mode" = "consult" ] && [ "$scope" != "uncommitted" ]; then
  fail "--scope is audit-only; consult reads repository context directly"
fi

if [ -z "$repo" ]; then
  if [ "$mode" = "consult" ]; then
    # A consult only reads; any directory works, git or not.
    repo="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
  else
    repo="$(git rev-parse --show-toplevel 2>/dev/null)" \
      || fail "not inside a git repository; pass --repo"
  fi
else
  [ -d "$repo" ] || fail "no such repo directory: $repo"
  if [ "$mode" != "consult" ]; then
    # A non-git --repo would hand the reviewer an empty scope behind silently
    # failing git calls.
    repo="$(git -C "$repo" rev-parse --show-toplevel 2>/dev/null)" \
      || fail "not a git work tree: $repo"
  fi
fi

classification='Number every finding (F1, F2, ...) and classify each by severity (High, Medium or Low) and kind:
- defect: wrong behavior, broken edge-case or contract violation introduced by this change
- nitpick: style, naming or polish; report every one, folding repeats of the same pattern into a single finding that lists its locations
- decision: needs an operator call (spec ambiguity, posture tradeoff); lay out the options and leave the choice open
- pre-existing: a real issue that predates this change
Report every concrete finding at every severity and of every kind — the verdict counts them all, and clean means zero findings remain. Ground each finding in code you actually read; speculation about unlikely preconditions, or extra defense-in-depth where the primary defense holds, is not a finding. When the brief records a disposition for a known issue (a rebuttal, an accepted tradeoff), weigh it openly in the report — state whether you concur and why — rather than silently honoring or re-raising it. Lead the report with a one-line verdict (clean, or counts per severity), cite a repository-relative file and line for each finding, and write web sources as plain markdown links. The report is the entire deliverable: end it after the last finding.'

web_note="Check claims about external interfaces (APIs, CLI flags, versions) against upstream documentation via web search when the repository alone cannot settle them."
if [ "$exec_mode" = 1 ]; then
  conduct="Verify behavior empirically where it sharpens the review: run the test suite, linters, or local infrastructure, and clean up whatever you start. ${web_note} Leave every tracked file untouched — report findings rather than fixing them."
else
  conduct="Rely on the repository state and any validation output provided rather than attempting builds, tests or other write-requiring commands. ${web_note}"
fi
classification="${classification} ${conduct}"

consult_prompt='You are a fresh outside advisor with read-only access to this repository. Read whatever bears on the question below, then give a candid, actionable answer — surface risks, alternatives, simplifications, and concrete next steps, and ground each point in what you read. Check external facts against upstream documentation via web search when the repository cannot settle them. Lead with your recommendation and keep it practical. This is advisory thinking rather than a code gate — favor useful guidance over severity tables or a pass/fail verdict.'

# Scope validation runs before anything lands on disk; the phrases that depend on
# the snapshot are composed after it exists.
if [ "$mode" = "audit" ]; then
  case "$scope" in
    uncommitted)
      if [ -z "$resume" ] && [ -z "$(git -C "$repo" status --porcelain --untracked-files=all -- . ':(exclude).local' 2>/dev/null)" ]; then
        fail "nothing to review: the working tree is clean"
      fi ;;
    commit:?*)
      target="${scope#commit:}"
      git -C "$repo" rev-parse -q --verify "${target}^{commit}" >/dev/null 2>&1 \
        || fail "unknown commit: ${target}"
      # Fix commits stack on top of a reviewed commit; a gate must widen to cover them.
      if [ "$(git -C "$repo" rev-parse "${target}^{commit}" 2>/dev/null)" != "$(git -C "$repo" rev-parse HEAD 2>/dev/null)" ] \
        && git -C "$repo" merge-base --is-ancestor "$target" HEAD 2>/dev/null; then
        echo "agent.sh: warning: HEAD has moved past ${target}; follow-up commits are outside this scope — base:${target}^ covers them" >&2 || true
      elif ! git -C "$repo" diff --quiet "$target" HEAD 2>/dev/null; then
        echo "agent.sh: warning: the checkout differs from ${target}; context reads follow HEAD" >&2 || true
      fi ;;
    base:?*)
      merge_base="$(git -C "$repo" merge-base "${scope#base:}" HEAD 2>/dev/null)" \
        || fail "no merge-base with: ${scope#base:}"
      if [ -z "$resume" ] && [ "$merge_base" = "$(git -C "$repo" rev-parse HEAD)" ]; then
        fail "nothing to review: HEAD is the merge-base with ${scope#base:}"
      fi ;;
    *) fail "unknown scope: $scope" ;;
  esac
fi
if [ "$mode" = "consult" ]; then
  [ -n "$brief" ] || fail "consult needs a brief — it is the whole question"
elif [ -n "$resume" ]; then
  [ -n "$brief" ] || fail "a verification pass needs a brief: what changed, where, and any rebuttals"
fi

dirty_tree=0
if [ "$mode" = "audit" ]; then
  case "$scope" in
    commit:?*|base:?*)
      if [ -n "$(git -C "$repo" status --porcelain --untracked-files=all -- . ':(exclude).local' 2>/dev/null)" ]; then
        dirty_tree=1
      fi ;;
  esac
fi

# The requested scope keeps driving diff generation after the verify relabel.
gate_scope="$scope"
ref=""
if [ "$mode" != "consult" ]; then
  case "$gate_scope" in
    commit:?*) ref="$(git -C "$repo" rev-parse --short "${gate_scope#commit:}" 2>/dev/null || echo n/a)" ;;
    *) ref="$(git -C "$repo" rev-parse --short HEAD 2>/dev/null || echo n/a)" ;;
  esac
fi

if [ "$mode" = "consult" ]; then
  out_dir="${repo}/.local/consults"
  report_prefix="consult"
else
  out_dir="${repo}/.local/audits"
  report_prefix="audit"
fi
mkdir -p "$out_dir"

run_started="$(date +%s)"
tmpdir="$(mktemp -d)"
# Composed next to its final home so the closing hard link stays on one filesystem.
composed="${out_dir}/.compose-$$-${slug}"
reviewer_pid=""
mcp_pid=""
claim=""
status=0
session=""
detected_model=""
snapshot_dir=""
snapshot_marker=""
snapshot_mode=""
review_root="$repo"
body="${tmpdir}/body" events="${tmpdir}/events" errlog="${tmpdir}/err"

# Every run leaves one provenance line in the report directory's ledger, whether
# it lands a report, fails, or is interrupted.
append_ledger() {
  jq -cn \
    --arg ts "$(TZ=Europe/Berlin date '+%Y-%m-%dT%H:%M:%S%z')" \
    --arg slug "$slug" --arg mode "$mode" --arg scope "$gate_scope" \
    --arg pass "$([ -n "$resume" ] && echo resumed || echo fresh)" \
    --arg ref "$ref" --arg reviewer "$reviewer" \
    --arg model "${detected_model:-${model}}" \
    --argjson exec "$([ "$exec_mode" = 1 ] && echo true || echo false)" \
    --arg snapshot "${snapshot_mode:-none}" \
    --arg session "$session" --arg resumed "$resume" \
    --arg outcome "$1" --arg path "$2" --argjson exit "$3" \
    '{ts: $ts, slug: $slug, mode: $mode,
      scope: (if $mode == "consult" then null else $scope end),
      pass: $pass, ref: (if $ref == "" then null else $ref end),
      reviewer: $reviewer, model: (if $model == "" then null else $model end),
      exec: $exec, snapshot: $snapshot,
      session: (if $session == "" then null else $session end),
      resumed: (if $resumed == "" then null else $resumed end),
      outcome: $outcome, path: $path, exit: $exit}' \
    >> "${out_dir}/index.jsonl" 2>/dev/null || true
}

# A signalled run drops a breadcrumb so a torn-down launch is never a silent
# empty log.
breadcrumb() {
  local crumb
  # The pid keeps concurrent same-slug runs from clobbering each other's crumbs.
  crumb="${out_dir}/${report_prefix}-$(TZ=Europe/Berlin date '+%Y-%m-%d-%H%M%S')-${slug}-$$-interrupted.log"
  {
    printf 'interrupted by %s at %s after %ss — no report written\n' "$1" \
      "$(TZ=Europe/Berlin date '+%Y-%m-%d %H:%M:%S %Z')" "$(( $(date +%s) - run_started ))"
    printf 'slug: %s\nmode: %s\nreviewer: %s\n' "$slug" "$mode" "$reviewer"
    [ "$mode" = "consult" ] || printf 'scope: %s\n' "$gate_scope"
  } > "$crumb" 2>/dev/null || true
  append_ledger interrupted "$crumb" "$2"
}

# Shared-clone snapshots: every object the snapshot creates lands inside the
# clone, so in steady state the audited repo's .git is never written and
# teardown is a plain delete. The path is stable per repo and reviewer because
# claude session resume is cwd-scoped (see references/reviewers.md).
snapshot_base="${TMPDIR:-/tmp}"
snapshot_base="${snapshot_base%/}/kindly-snapshots"

cleanup_snapshot() {
  [ -n "$snapshot_dir" ] || return 0
  rm -rf "$snapshot_dir" 2>/dev/null || true
  if [ -n "$snapshot_marker" ]; then
    rm -f "$snapshot_marker" 2>/dev/null || true
  fi
  # The shared base directory leaves with its last snapshot; a concurrent
  # run's dir or marker keeps the rmdir from firing.
  rmdir "$snapshot_base" 2>/dev/null || true
  snapshot_dir=""
}

create_snapshot() {
  local commit tree temp_index owner repo_exclude repo_excludes_file legacy_admin
  local stray stray_dir stray_owner needs_prune=0
  local add_pathspec=(-- . ':(exclude).local')
  # Naming an ignored path in a pathspec makes `git add` fail outright, so the
  # exclude is passed only where .local is not already ignored.
  if git -C "$repo" check-ignore -q .local 2>/dev/null; then
    add_pathspec=(-- .)
  fi
  mkdir -p "$snapshot_base" 2>/dev/null || return 1
  snapshot_dir="${snapshot_base}/$(basename "$repo")-$(printf '%s' "$repo" | shasum | cut -c1-12)-${reviewer}"
  snapshot_marker="${snapshot_dir}.pid"
  legacy_admin="$(git -C "$repo" rev-parse --path-format=absolute --git-path worktrees 2>/dev/null || true)"
  # Sweep crashed step-aside leftovers — no later run claims their suffixed
  # paths. A live owner keeps its checkout.
  for stray in "${snapshot_dir}"-*; do
    [ -e "$stray" ] || continue
    stray_dir="${stray%.pid}"
    stray_owner="$(cat "${stray_dir}.pid" 2>/dev/null || true)"
    if [ -n "$stray_owner" ] && kill -0 "$stray_owner" 2>/dev/null; then
      continue
    fi
    rm -rf "$stray_dir" 2>/dev/null || true
    rm -f "${stray_dir}.pid" 2>/dev/null || true
    if [ -n "$legacy_admin" ] && [ -d "${legacy_admin}/$(basename "$stray_dir")" ]; then
      needs_prune=1
    fi
  done
  # Heal a crashed run's leftovers; step aside onto a unique path (losing resume
  # continuity for this run only) when a live concurrent run owns the stable one.
  if [ -e "$snapshot_dir" ] || [ -e "$snapshot_marker" ]; then
    owner="$(cat "$snapshot_marker" 2>/dev/null || true)"
    if [ -n "$owner" ] && kill -0 "$owner" 2>/dev/null; then
      snapshot_dir="${snapshot_dir}-$$"
      snapshot_marker="${snapshot_dir}.pid"
      echo "agent.sh: another review holds the stable snapshot path; using ${snapshot_dir} (its session will not be resumable)" >&2 || true
    else
      rm -rf "$snapshot_dir" 2>/dev/null || true
      rm -f "$snapshot_marker" 2>/dev/null || true
    fi
  fi
  # An earlier git-worktree implementation registered snapshots in the repo's
  # .git/worktrees; prune clears any lingering admin files. A still-present
  # worktree (a live concurrent old-style run) is left alone.
  if [ -n "$legacy_admin" ] && [ -d "${legacy_admin}/$(basename "$snapshot_dir")" ]; then
    needs_prune=1
  fi
  if [ "$needs_prune" = 1 ]; then
    git -C "$repo" worktree prune 2>>"$errlog" || true
  fi
  printf '%s' "$$" > "$snapshot_marker" 2>/dev/null || return 1
  git clone --quiet --shared --no-checkout -- "$repo" "$snapshot_dir" 2>>"$errlog" || return 1
  # The clone's fresh git-dir drops the repo's info/exclude and local
  # excludes-file config; carry both over so local-only ignore rules hold.
  repo_exclude="$(git -C "$repo" rev-parse --path-format=absolute --git-path info/exclude 2>/dev/null || true)"
  if [ -n "$repo_exclude" ] && [ -f "$repo_exclude" ]; then
    { mkdir -p "${snapshot_dir}/.git/info" \
        && cp "$repo_exclude" "${snapshot_dir}/.git/info/exclude"; } 2>>"$errlog" || return 1
  fi
  repo_excludes_file="$(git -C "$repo" config --local --get core.excludesFile 2>/dev/null || true)"
  if [ -n "$repo_excludes_file" ]; then
    git -C "$snapshot_dir" config core.excludesFile "$repo_excludes_file" 2>>"$errlog" || return 1
  fi
  case "$gate_scope" in
    uncommitted)
      # A temp-index commit captures the live working tree into the clone's
      # object store, parented on HEAD so the snapshot's own HEAD-diff is the
      # scope; the live index is never touched.
      temp_index="${tmpdir}/snapshot-index"
      { GIT_INDEX_FILE="$temp_index" git -C "$repo" --git-dir="${snapshot_dir}/.git" --work-tree=. read-tree HEAD \
          && GIT_INDEX_FILE="$temp_index" git -C "$repo" --git-dir="${snapshot_dir}/.git" --work-tree=. add -A "${add_pathspec[@]}" \
          && tree="$(GIT_INDEX_FILE="$temp_index" git -C "$repo" --git-dir="${snapshot_dir}/.git" --work-tree=. write-tree)" \
          && commit="$(git --git-dir="${snapshot_dir}/.git" -c user.name=kindly -c user.email=kindly@localhost \
              commit-tree "$tree" -p HEAD -m 'kindly audit snapshot')"
      } 2>>"$errlog" || return 1
      ;;
    *)
      commit="$(git -C "$repo" rev-parse HEAD 2>>"$errlog")" || return 1
      ;;
  esac
  git -C "$snapshot_dir" checkout --quiet --detach "$commit" 2>>"$errlog" || return 1
}

trap 'cleanup_snapshot; rm -rf "$tmpdir" "$composed"; if [ -n "$claim" ]; then rmdir "$claim" 2>/dev/null || true; fi' EXIT
trap 'if [ -n "$reviewer_pid" ]; then kill_tree "$reviewer_pid"; fi; if [ -n "$mcp_pid" ]; then kill_tree "$mcp_pid"; fi; breadcrumb SIGINT 130; exit 130' INT
trap 'if [ -n "$reviewer_pid" ]; then kill_tree "$reviewer_pid"; fi; if [ -n "$mcp_pid" ]; then kill_tree "$mcp_pid"; fi; breadcrumb SIGTERM 143; exit 143' TERM

if [ "$mode" = "audit" ] && [ "$exec_mode" = 0 ]; then
  # A snapshot commit holds one version per path; when staged content diverges
  # from the on-disk tree, only the live review (with its staged-and-unstaged
  # prepared diff) keeps both versions in scope.
  staged_divergence=""
  if [ "$gate_scope" = "uncommitted" ]; then
    staged_divergence="$(comm -12 \
      <(git -C "$repo" diff --cached --name-only -- . ':(exclude).local' 2>/dev/null | sort) \
      <(git -C "$repo" diff-files --name-only -- . ':(exclude).local' 2>/dev/null | sort) \
      2>/dev/null || true)"
  fi
  if [ -n "$staged_divergence" ]; then
    snapshot_mode="live-fallback"
    echo "agent.sh: staged content diverges from the working tree ($(printf '%s' "$staged_divergence" | tr '\n' ' ')); reviewing the live tree so both versions stay in scope — keep it frozen until the report lands" >&2 || true
  elif create_snapshot; then
    snapshot_mode="checkout"
    review_root="$snapshot_dir"
    if [ -e "${snapshot_dir}/.gitmodules" ]; then
      echo "agent.sh: warning: the reviewed checkout defines submodules; the snapshot leaves them unpopulated, so findings inside submodule content need a live-tree pass" >&2 || true
    fi
    if [ "$dirty_tree" = 1 ]; then
      echo "agent.sh: the working tree is dirty; the snapshot pins context reads to a clean HEAD checkout" >&2 || true
    fi
  else
    cleanup_snapshot
    snapshot_mode="live-fallback"
    echo "agent.sh: warning: could not create a snapshot checkout; reviewing the live tree — keep it frozen until the report lands" >&2 || true
    if [ "$dirty_tree" = 1 ]; then
      echo "agent.sh: warning: the working tree is dirty; context reads can leak uncommitted edits into this committed-scope review" >&2 || true
    fi
  fi
elif [ "$mode" = "audit" ] && [ "$dirty_tree" = 1 ]; then
  echo "agent.sh: warning: the working tree is dirty; context reads can leak uncommitted edits into this committed-scope review" >&2 || true
fi

scope_phrase=""
if [ "$mode" = "audit" ]; then
  case "$gate_scope" in
    uncommitted)
      if [ "$snapshot_mode" = "checkout" ]; then
        scope_phrase="the changes this checkout's HEAD commit introduces over its parent — a snapshot of the caller's uncommitted work (staged, unstaged and untracked, excluding .local/)"
      else
        scope_phrase="the uncommitted changes (staged, unstaged and untracked, excluding .local/ audit and consult artifacts)"
      fi ;;
    commit:?*)
      scope_phrase="commit ${gate_scope#commit:}"
      if [ "$snapshot_mode" = "checkout" ]; then
        scope_phrase="${scope_phrase} (the checkout is a clean snapshot at HEAD)"
      else
        scope_phrase="${scope_phrase} (ignore uncommitted edits)"
      fi ;;
    base:?*)
      scope_phrase="all committed changes since the merge-base with ${gate_scope#base:}"
      if [ "$snapshot_mode" = "checkout" ]; then
        scope_phrase="${scope_phrase} (the checkout is a clean snapshot at HEAD)"
      else
        scope_phrase="${scope_phrase} (ignore uncommitted edits)"
      fi ;;
  esac
fi

if [ "$mode" = "consult" ]; then
  prompt="${consult_prompt}

${brief}"
elif [ -n "$resume" ]; then
  scope="verify"
  snapshot_note=""
  [ "$snapshot_mode" = "checkout" ] \
    && snapshot_note=" — this pass runs in a fresh checkout of the repository, so resolve every file by its repository-relative path"
  prompt="The numbered findings from your review have been worked on. Re-read every file involved (your earlier context is stale${snapshot_note}), then verify each finding individually: read the fix code itself and judge whether it truly closes the finding without introducing new issues. Report the outcome per finding number, plus every remaining or new finding of any severity and kind, under the same classification, citation and verdict rules as your first pass. This pass converges the gate; closing it takes a fresh full-scope review reporting zero findings. ${conduct}

${brief}"
else
  prompt="Review ${scope_phrase} and give me a full report of any issues, inconsistencies, unhandled edge-cases, bloated or stale documentation, or suboptimal, half-baked or sloppy work. ${classification}${brief:+

${brief}}"
fi

# Status alone keeps the same line when an already-dirty file is edited again,
# and composing the report itself must not flag the tree as mutated.
tree_state() {
  { git -C "$review_root" status --porcelain --untracked-files=all -- . ':(exclude).local'
    git -C "$review_root" diff -- . ':(exclude).local'
    git -C "$review_root" diff --cached -- . ':(exclude).local'
    (cd "$review_root" && git ls-files -z --others --exclude-standard -- . ':(exclude).local' \
      | sort -z | xargs -0 shasum)
  } 2>/dev/null | shasum | cut -d' ' -f1
}

start_reviewer_attempt() {
  local attempt="$1"
  printf -- '--- %s attempt %s stderr ---\n' "$reviewer" "$attempt" >> "$errlog"
  case "$reviewer" in
    codex)
      (cd "$review_root" && trap - PIPE && exec_isolated codex "${args[@]}" "$prompt" </dev/null) >"$events" 2>>"$errlog" &
      ;;
    claude)
      (cd "$review_root" && trap - PIPE && exec_isolated env -u CLAUDECODE -u CLAUDE_CODE_SESSION_ID -u CLAUDE_CODE_ENTRYPOINT \
        -u CLAUDE_CODE_ENABLE_TASKS -u CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS -u AI_AGENT \
        claude "${args[@]}" </dev/null) >"$events" 2>>"$errlog" &
      ;;
  esac
  reviewer_pid=$!
}

run_with_first_event_watchdog() {
  local attempt waited final_no_event no_event_msg reviewer_started first_event_secs heartbeat_secs
  first_event_secs="$(positive_int_env KINDLY_FIRST_EVENT_SECS 120)"
  heartbeat_secs="$(positive_int_env KINDLY_HEARTBEAT_SECS 30)"
  attempt=0
  while :; do
    attempt=$((attempt + 1))
    : > "$events"
    : > "$body"
    status=0
    reviewer_started="$(date +%s)"
    start_reviewer_attempt "$attempt"
    waited=0
    while [ "$waited" -lt "$first_event_secs" ] && ! [ -s "$events" ] && kill -0 "$reviewer_pid" 2>/dev/null; do
      sleep 1
      waited=$((waited + 1))
      if [ $((waited % heartbeat_secs)) -eq 0 ] && [ "$waited" -lt "$first_event_secs" ] \
        && ! [ -s "$events" ] && kill -0 "$reviewer_pid" 2>/dev/null; then
        echo "agent.sh: waiting for first ${reviewer} reviewer event (${waited}s elapsed)" >&2 || true
      fi
    done
    if [ -s "$events" ]; then
      wait_with_progress "$reviewer_pid" "$events" "$reviewer" "$reviewer_started" || status=$?
      break
    fi
    if ! kill -0 "$reviewer_pid" 2>/dev/null; then
      if wait "$reviewer_pid"; then
        status=0
      else
        status=$?
      fi
      no_event_msg="no ${reviewer} reviewer events before process exit ${status}"
      echo "agent.sh: ${no_event_msg}" >&2 || true
      echo "$no_event_msg" >> "$errlog"
      break
    fi
    final_no_event=0
    no_event_msg="no ${reviewer} reviewer events after ${waited}s; ending attempt ${attempt} and retrying"
    if [ "$attempt" -ge 2 ]; then
      final_no_event=1
      no_event_msg="no ${reviewer} reviewer events after ${waited}s on attempt ${attempt}; giving up"
    fi
    echo "agent.sh: ${no_event_msg}" >&2 || true
    echo "$no_event_msg" >> "$errlog"
    kill_tree "$reviewer_pid"
    wait "$reviewer_pid" 2>/dev/null || true
    if [ "$final_no_event" = 1 ]; then
      status=1
      break
    fi
  done
}

sandbox="read-only"
[ "$exec_mode" = 1 ] && sandbox="workspace-write"
# A tree-hash hiccup must never abort the run — least of all after a finished review.
tree_before="$(tree_state || true)"
# Backgrounded reviewers stay under the INT/TERM traps, which kill the full
# process tree through reviewer_pid.
if [ "$reviewer" = "codex" ]; then
  args=(exec)
  [ "$mode" = "consult" ] && args+=(--skip-git-repo-check)
  if [ -n "$resume" ]; then
    args+=(resume "$resume" -c "sandbox_mode=\"${sandbox}\"")
  else
    args+=(--sandbox "$sandbox")
  fi
  args+=(-c 'web_search="live"')
  args+=(-c 'approval_policy="never"')
  if [ "$exec_mode" = 1 ]; then
    args+=(-c 'sandbox_workspace_write.network_access=true')
  fi
  # The operator-curated MCP servers stay enabled. They run outside the shell
  # sandbox, so the operator's curation is the boundary for what they can reach.
  args+=(--json -o "$body")
  [ -n "$model" ] && args+=(-m "$model")
  args+=(-c 'model_reasoning_effort="xhigh"')
  run_with_first_event_watchdog
  session="$({ jq -Rr 'fromjson? | select(type == "object" and .type == "thread.started") | .thread_id // empty' "$events" 2>/dev/null || true; } | head -1)"
else
  if [ "$exec_mode" = 1 ]; then
    tool_note="Your tools are Read, Glob, Grep and Bash, plus WebSearch and WebFetch for upstream documentation."
  else
    tool_note="Your tools are Read, Glob and Grep — no shell — plus WebSearch and WebFetch for upstream documentation."
  fi
  if [ "$mode" = "consult" ]; then
    prompt="${tool_note}

${prompt}"
  else
    diff_file="${tmpdir}/scope.patch"
    case "$gate_scope" in
      uncommitted)
        if [ "$snapshot_mode" = "checkout" ]; then
          # The snapshot commit carries the whole scope, untracked content included.
          git -C "$review_root" diff "HEAD^" HEAD > "$diff_file" 2>/dev/null || : > "$diff_file"
        # On a clean verify-pass tree, emit no diff file rather than bare section headers.
        elif [ -n "$(git -C "$repo" status --porcelain --untracked-files=all -- . ':(exclude).local' 2>/dev/null)" ]; then
          { git -C "$repo" status --short --untracked-files=all -- . ':(exclude).local'
            echo
            echo "# Staged (index vs HEAD)"
            git -C "$repo" diff --cached -- . ':(exclude).local'
            echo
            echo "# Unstaged (worktree vs index)"
            git -C "$repo" diff -- . ':(exclude).local'
            echo
            echo "# Untracked files (contents live on disk — read each)"
            git -C "$repo" ls-files --others --exclude-standard -- . ':(exclude).local'
          } > "$diff_file"
        fi ;;
      commit:?*)
        git -C "$review_root" show "${gate_scope#commit:}" > "$diff_file" ;;
      base:?*)
        git -C "$review_root" diff "$merge_base" HEAD > "$diff_file" ;;
    esac
    if [ -s "$diff_file" ]; then
      untracked_note=""
      [ "$gate_scope" = "uncommitted" ] && [ "$snapshot_mode" != "checkout" ] \
        && untracked_note=" Untracked files appear as listed paths only — read each one directly."
      prompt="${tool_note} The exact diff for this scope is already prepared at ${diff_file}; read it first, then read repository files for context.${untracked_note}

${prompt}"
    fi
  fi
  # Under dontAsk a visible-yet-unapproved MCP tool is silently denied, so
  # each curated server gets an allow rule, named with the installed CLI's
  # sanitizer. The listing health-checks servers over the network, so it runs
  # bounded — a wedged listing degrades to no rules, never a hang.
  echo "agent.sh: enumerating MCP servers for the reviewer allowlist" >&2 || true
  mcp_rules=""
  mcp_listing_file="${tmpdir}/mcp-list"
  : > "$mcp_listing_file"
  (cd "$review_root" && claude mcp list </dev/null >"$mcp_listing_file" 2>/dev/null) &
  mcp_pid=$!
  mcp_waited=0
  mcp_timed_out=0
  mcp_list_secs="$(positive_int_env KINDLY_MCP_LIST_SECS 120)"
  while kill -0 "$mcp_pid" 2>/dev/null; do
    if [ "$mcp_waited" -ge "$mcp_list_secs" ]; then
      mcp_timed_out=1
      kill_tree "$mcp_pid"
      # A partial listing must not pass for a complete one: no rules beat
      # silently missing some servers.
      : > "$mcp_listing_file"
      echo "agent.sh: warning: the MCP server listing timed out after ${mcp_waited}s; proceeding without allow rules" >&2 || true
      break
    fi
    sleep 1
    mcp_waited=$((mcp_waited + 1))
  done
  mcp_status=0
  wait "$mcp_pid" 2>/dev/null || mcp_status=$?
  mcp_pid=""
  # The listing exits zero even when servers are unhealthy (verified live), so a
  # nonzero exit means the listing itself broke mid-write — discard it whole.
  if [ "$mcp_status" -ne 0 ] && [ "$mcp_timed_out" = 0 ]; then
    : > "$mcp_listing_file"
    echo "agent.sh: warning: the MCP server listing failed (exit ${mcp_status}); proceeding without allow rules" >&2 || true
  fi
  while IFS= read -r mcp_line; do
    case "$mcp_line" in
      (*": "*" - "*) mcp_name="${mcp_line%%: *}" ;;
      (*) continue ;;
    esac
    [ -n "$mcp_name" ] || continue
    mcp_tool="$(printf '%s' "$mcp_name" | sed 's/[^A-Za-z0-9_-]/_/g')"
    case "$mcp_name" in
      ("claude.ai "*) mcp_tool="$(printf '%s' "$mcp_tool" | sed -e 's/__*/_/g' -e 's/^_*//' -e 's/_*$//')" ;;
    esac
    mcp_rules="${mcp_rules:+${mcp_rules},}mcp__${mcp_tool}__*"
  done < "$mcp_listing_file"
  if [ "$mcp_status" -eq 0 ] && [ "$mcp_timed_out" = 0 ] && [ -z "$mcp_rules" ] && [ ! -s "$mcp_listing_file" ]; then
    echo "agent.sh: warning: could not enumerate MCP servers; the reviewer may see MCP tools it cannot use" >&2 || true
  fi
  reviewer_tools='Read,Glob,Grep,WebSearch,WebFetch'
  [ "$exec_mode" = 1 ] && reviewer_tools="${reviewer_tools},Bash"
  args=(-p "$prompt" --permission-mode dontAsk --effort xhigh)
  args+=(--output-format stream-json --verbose)
  args+=(--tools "$reviewer_tools")
  args+=(--allowedTools "${reviewer_tools}${mcp_rules:+,${mcp_rules}}" --add-dir "$tmpdir")
  if [ "$exec_mode" = 1 ]; then
    args+=(--settings '{"sandbox":{"enabled":true,"autoAllowBashIfSandboxed":false,"allowUnsandboxedCommands":false,"failIfUnavailable":true,"filesystem":{"allowWrite":["/tmp","/private/tmp","/var/folders"]},"network":{"allowLocalBinding":true}}}')
  fi
  [ -n "$model" ] && args+=(--model "$model")
  [ -n "$resume" ] && args+=(--resume "$resume")
  run_with_first_event_watchdog
  jq -Rr 'fromjson? | select(type == "object" and .type == "result") | .result // empty' "$events" >"$body" 2>/dev/null || true
  session="$({ jq -Rr 'fromjson? | select(type == "object" and .type == "system" and .subtype == "init") | .session_id // empty' \
    "$events" 2>/dev/null || true; } | head -1)"
  detected_model="$({ jq -Rr 'fromjson? | select(type == "object" and .type == "system" and .subtype == "init") | .model // empty' \
    "$events" 2>/dev/null || true; } | head -1)"
fi

reviewer_exit="$status"
usable_report=0
final_result_subtype="$(event_final_result_subtype "$events")"
final_result_failed="$(event_final_result_failed "$events")"
if [ -s "$body" ]; then
  usable_report=1
  if [ "$reviewer" = "claude" ]; then
    case "$final_result_subtype" in
      ""|success) ;;
      *) usable_report=0 ;;
    esac
    [ "$final_result_failed" = "yes" ] && usable_report=0
  fi
fi

if [ "$usable_report" != 1 ]; then
  # tmpdir (with the reviewer log) is wiped on EXIT; keep the diagnostics first.
  # The pid keeps concurrent same-slug runs from clobbering each other's logs.
  fail_log="${out_dir}/${report_prefix}-$(TZ=Europe/Berlin date '+%Y-%m-%d-%H%M%S')-${slug}-$$-failed.log"
  {
    printf 'slug: %s\nmode: %s\n' "$slug" "$mode"
    [ "$mode" = "consult" ] || printf 'scope: %s\n' "$gate_scope"
    printf 'reviewer: %s\nexit: %s\nelapsed: %ss\n\n--- stderr ---\n' \
      "$reviewer" "$status" "$(( $(date +%s) - run_started ))"
    cat "$errlog" 2>/dev/null
    printf '\n--- failure details ---\n'
    event_failure_details "$events"
    stderr_failure_details "$errlog"
    printf '\n--- event metadata tail ---\n'
    event_metadata_tail "$events"
  } > "$fail_log" 2>/dev/null || true
  {
    echo "agent.sh: ${reviewer} review failed (exit ${status}$([ -s "$body" ] || echo ', empty report')$([ "$stalled" = 1 ] && echo ', stalled' || true))"
    tail -5 "$errlog" 2>/dev/null
    event_error_summary "$events" || true
    event_recovery_hint "$events" || true
    stderr_recovery_hint "$errlog" || true
    [ -s "$fail_log" ] && echo "agent.sh: diagnostics saved to ${fail_log}"
  } >&2 || true
  failed_exit="$((status == 0 ? 1 : status))"
  append_ledger failed "$fail_log" "$failed_exit"
  exit "$failed_exit"
fi

if [ "$reviewer_exit" -ne 0 ]; then
  echo "agent.sh: warning: ${reviewer} exited ${reviewer_exit} after writing a usable report; preserving report and marking reviewer-exit" >&2 || true
fi

{
  printf -- '---\n'
  printf 'slug: %s\n' "$slug"
  if [ "$mode" = "consult" ]; then
    printf 'mode: consult\n'
  elif [ "$scope" = "verify" ]; then
    printf 'scope: verify (%s)\n' "$gate_scope"
  else
    printf 'scope: %s\n' "$scope"
  fi
  printf 'repo: %s\n' "$repo"
  if [ "$mode" != "consult" ]; then
    printf 'ref: %s\n' "$ref"
  fi
  printf 'reviewer: %s%s (xhigh)\n' "$reviewer" "${model:+ ${model}}"
  if [ -n "$detected_model" ]; then
    printf 'model: %s\n' "$detected_model"
  fi
  if [ "$same_reviewer" = 1 ]; then
    printf 'same-reviewer: yes\n'
  fi
  if [ -n "$snapshot_mode" ]; then
    printf 'snapshot: %s\n' "$snapshot_mode"
  fi
  if [ "$reviewer_exit" -ne 0 ]; then
    printf 'reviewer-exit: %s\n' "$reviewer_exit"
  fi
  if [ "$exec_mode" = 1 ]; then
    printf 'exec: enabled\n'
  fi
  if [ "$dirty_tree" = 1 ] && [ "$snapshot_mode" != "checkout" ]; then
    printf 'dirty-tree: yes\n'
  fi
  tree_after="$(tree_state || true)"
  if [ "$tree_after" != "$tree_before" ]; then
    printf 'tree-mutated: yes\n'
    echo "agent.sh: warning: the working tree changed during the review — inspect git status before trusting it" >&2 || true
  fi
  printf 'session: %s\n' "${session:-unknown}"
  if [ -n "$resume" ]; then
    printf 'resumed: %s\n' "$resume"
  fi
  printf 'date: %s\n' "$(TZ=Europe/Berlin date '+%Y-%m-%d %H:%M %Z')"
  printf -- '---\n\n'
  cat "$body"
  printf '\n'
} > "$composed"

# A Berlin-local timestamp plus the slug names the report; the mkdir claim is the
# mutex so parallel or panel runs finishing in the same second never clobber one
# another, and the hard link lands the report in full or not at all.
stamp="$(TZ=Europe/Berlin date '+%Y-%m-%d-%H%M%S')"
report=""
suffix=""
n=1
while [ -z "$report" ]; do
  name="${report_prefix}-${stamp}-${slug}${suffix}"
  if [ ! -e "${out_dir}/${name}.md" ] \
    && mkdir "${out_dir}/.claim-${name}" 2>/dev/null; then
    claim="${out_dir}/.claim-${name}"
    report="${out_dir}/${name}.md"
    ln "$composed" "$report"
    rm -f "$composed"
    rmdir "$claim"
    claim=""
  else
    n=$((n + 1))
    suffix="-${n}"
  fi
done

append_ledger report "$report" 0
echo "report: $report" || true
echo "session: ${session:-unknown}" || true
