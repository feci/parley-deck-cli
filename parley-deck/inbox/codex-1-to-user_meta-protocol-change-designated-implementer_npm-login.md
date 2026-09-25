---
from: codex-1
to: user
idea: meta-protocol-change-designated-implementer
phase: release
status: owner-action
blocking: no
date: 2026-09-25
---

# npm login needed for 2.14.0

The organizer attempted the owner-authorized exact-tarball command:

```sh
npm publish --access public '/Volumes/My Shared Files/AI_WORKSPACE/parley-deck/release-delivery/2026-09-25-designated-implementer/skill/parley-deck-skill-2.14.0.tgz'
```

It exited 1: E404 on PUT https://registry.npmjs.org/parley-deck-skill, reporting the resource not found or insufficient access. The existing public package was independently observed at 2.13.0; 2.14.0 publication is not claimed. Per the owner's release brief, please run `! npm login` in the Claude session, then report success so the exact audited tarball can be retried. No credential is requested or copied. Other channels continue.

Audited tarball SHA256: 681b9076a017e67935403ad6cd9779cafa09a117d168331daca300404336135b. Log: release-delivery/2026-09-25-designated-implementer/npm-publication.log.
