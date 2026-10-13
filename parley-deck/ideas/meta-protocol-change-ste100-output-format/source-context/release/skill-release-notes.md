Parley Deck now uses ASD-STE100 Issue 9 as the default for new English protocol prose.

The new `references/STE100.md` has the same rule body as protocol Appendix B. The short core pointer keeps SKILL.md at 19,482 bytes under its unchanged 20,000-byte limit.

Exact source quotes, including non-English quotes, remain intact. The working summary defines exemptions, adoption, per-idea exceptions, and advisory style review. Existing artifacts remain unchanged.

Local validation passed 402 Node tests, 54 Python tests, six payload manifests, package checks, the native portable build, and a fresh installer/doctor check. CI also passed on Python 3.10 and 3.13.

The companion CLI is 1.55.0. Core 2.19.0 is staged for attended publication. This GitHub release does not publish the package to npm.
Windows portable assets come from the release workflow's cross-build. They are not evidence of a native Windows installation test.
