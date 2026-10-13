# Release evidence

These files record CLI 1.55.0 and skill 2.19.0 delivery.
The signed source plan and closed IMPLEMENTATION remain in the parent idea.
The separate channel report and final release inbox note provide the final delivery result.

Preserve raw logs and tool output, including whitespace. The three WinGet manifests retain their required CRLF bytes.
Their recorded hashes cover those exact bytes. A whitespace cleanup would invalidate that evidence.
Whitespace checks apply to authored text. Git's cr-at-eol setting permits the retained manifest format during the audit check.

The scripts preserve the actual delivery mechanism. They are audit material, not a batch entry point.
Some scripts mutate live releases, repositories, or local runtimes. Do not replay them without the required state and authority.
The first Windows comparison script records one baseline sample. The final comparison also includes the later main baseline and independent Zcode adjudication.
See windows-final-comparison.json in the parent source-context directory for both witnesses and their limits.

The private runtime backup remains outside git. No credential or private prompt file belongs in this directory.
