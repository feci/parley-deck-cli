# Kimi readiness failures — 2026-10-09

Two fresh measured CLI readiness attempts exited 1; neither authored an artifact. The retry started more than five seconds after the first terminal outcome. No credential/model/gateway change was made. These are actual process failures, not quota-recognizer-positive native evidence.

- Invocation `d8a37618-2af7-4072-be4c-d339e662ec85`; started 2026-10-09T06:45:51.846135Z; completed 2026-10-09T06:45:53.393432Z; duration 1583ms; exit 1; class process_failure.
```text
error: failed to run prompt: provider.auth_error: 403 [kimi-coding] 1 connection(s) exist but are excluded by this API key's connection allowlist / quota scope — add them to the key in the dashboard, or use a key without that scope
```

- Invocation `7939933e-5ba9-4e0a-8f85-e25c53542035`; started 2026-10-09T06:47:02.321007Z; completed 2026-10-09T06:47:03.781862Z; duration 1492ms; exit 1; class process_failure.
```text
error: failed to run prompt: provider.auth_error: 403 [kimi-coding] 1 connection(s) exist but are excluded by this API key's connection allowlist / quota scope — add them to the key in the dashboard, or use a key without that scope
```
