#!/bin/zsh
# claude-1 round-05 focused PRIMARY test runs; no provider (PATH shadowed by no-provider-bin)
cd "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude"
export PATH="/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/quota-implementation/fixup-2/host/no-provider-bin:$PATH"
PKGS=(./internal/telemetry ./internal/membership ./internal/app ./internal/pidlease ./internal/fsutil ./internal/consensus ./internal/quota ./internal/protocol ./internal/driver ./internal/runner ./internal/config ./internal/runcontrol ./internal/runstate)
RUN='Quota|Cycle2|Lease|SyncFile|EmbeddedDefaultMatchesLiveDeck'
echo "start $(date -u +%FT%TZ) git=$(git rev-parse HEAD) go=$(go version)"
echo "== shared TMPDIR"
TMPDIR="/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/claude1-r5/tests/shared-tmp" go test $PKGS -run "$RUN" -count=1 > "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/claude1-r5/tests/shared-focused.log" 2>&1; echo "shared exit=$?"
echo "== local TMPDIR"
TMPDIR=/tmp/claude1-r5-local go test $PKGS -run "$RUN" -count=1 > "/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/worktrees/quota-auto-exclude/.parley-runtime/claude1-r5/tests/local-focused.log" 2>&1; echo "local exit=$?"
echo "end $(date -u +%FT%TZ)"
