---
from: codex-1
to: user
idea: meta-protocol-change-lean-organizer
phase: release
blocking: yes
date: 2026-09-24
---

## Current blocker

npm latest remains 2.12.1. The exact audited 2.13.0 tarball is NOT published. Initial publication returned E404 and whoami E401; a later whoami succeeded, but publishing required EOTP web verification. Both unattended verification links expired. No active authorization challenge remains, and the expired URLs have been removed from this current handoff. No credentials were read or copied; agent browser launching was disabled.

## Owner action

Complete publication and npm's interactive verification directly in your own terminal:

```sh
npm publish '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-24-lean-organizer/parley-deck-skill-2.13.0.tgz' --access public
```

Alternatively, notify this session when ready to approve a fresh short-lived request. Do not send credentials or OTPs in chat. The preserved tarball SHA256 is dcf9c75ca84c3d8f0eb475cf745e4e4bd0973269740e334b0cd8c01d7608d802. After publication, the independent auditor must compare registry latest/integrity and the downloaded registry tarball with this exact artifact.
