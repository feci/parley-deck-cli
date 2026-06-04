---
agent: hermes
idea: unified-tui-home
review-cycle: 3
reviewed-commit: 245b9eb
---

AF1 re-review: launchReaper + waitForActive verified. N-launched runs now complete and record sessions before parent cancel fires on TUI detach; ctrl+c during wait still aborts as expected. TestLaunchReaperWaitsForInFlightRuns passes and asserts blocking behavior. AF2–AF6 remain unchanged and correct. Build and all tests green.