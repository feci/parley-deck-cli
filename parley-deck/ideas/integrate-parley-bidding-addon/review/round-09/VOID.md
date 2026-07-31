# Round 9 is void — facilitator note, not a participant artifact

Round 9 ran twice. Neither run is usable, and both reasons are the facilitator's.

1. **First launch — network outage.** `codex-1` failed DNS lookup on `wss://chatgpt.com` after
   five reconnects, `hermes-1` returned `Connection error`, `kimi-1` failed
   `getaddrinfo ENOTFOUND auth.kimi.com`. No artifact from anyone. Recorded as an outage, never
   as an accept.
2. **Second launch — two facilitator errors.** `codex-1` was given a sandbox writable root
   covering only `parley-deck-skill`, so it completed a full review and could not write it here
   (verdict BLOCK with one MAJOR, stdout only). And the facilitator edited `lib/installer.js`
   and `test/bidding-addon.test.js` while `hermes-1` and `kimi-1` were still reading the tree.

`hermes-1.md` (ACCEPT) and `kimi-1.md` (BLOCK) are kept as evidence. They count as **neither
signoff nor accept**: their front matter says `reviewed-commit: dcd200e`, and the tree did not
stay at `dcd200e` for the duration.

Both MAJORs raised in this round were verified independently by the implementer and fixed in
cycles 12 and 13. Round 10 is the fresh full-scope round the strict gate requires, at `9ed2081`,
with the tree frozen for its duration.
