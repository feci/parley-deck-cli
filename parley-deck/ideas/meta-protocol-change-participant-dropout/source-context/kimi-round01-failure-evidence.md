# Kimi round-01 failures

Two separate configured launches of the same Phase-1 step; no model or provider was changed. Private raw logs remain under .parley-runtime/invocations and are not copied or committed. This 400 is neither 429/503 nor a timeout, so the 15-minute retry delay does not apply.

## Attempt 1

Invocation: `809e8add-4b99-45f8-8ceb-b2f46d2259f8`. Started 2026-10-08T21:13:43.150209Z; completed 2026-10-08T21:13:44.939891Z; exit 1; status failed. No canonical artifact.

Verbatim decisive stderr (no secrets):

```text
error: failed to run prompt: provider.api_error: 400 Ambiguous model 'k3'. Use provider/model prefix (ex: kmca/k3 or kmc/k3).
```

## Attempt 2

Invocation: `7358f2b2-3ea6-4ba0-80d9-399d2cfa323a`. Started 2026-10-08T21:14:27.215495Z; completed 2026-10-08T21:14:28.588701Z; exit 1; status failed. No canonical artifact.

Verbatim decisive stderr (no secrets):

```text
error: failed to run prompt: provider.api_error: 400 Ambiguous model 'k3'. Use provider/model prefix (ex: kmca/k3 or kmc/k3).
```
