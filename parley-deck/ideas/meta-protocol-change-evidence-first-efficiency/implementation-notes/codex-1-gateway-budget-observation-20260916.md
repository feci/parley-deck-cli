# Codex gateway budget observation — 2026-09-16

Read-only ego-browser inspection of the actual configured gateway adds a concrete
limit to the experiment spending analysis. No model request, key creation, budget
edit, subscription change or account configuration change was made.

The selected Kimi model `kimi-code/k3` uses the local config provider
`omniroute-kimi`, type `kimi`, base URL `https://omniroute.marao.sk/v1`
(`/Users/tomasfecko/.kimi-code/config.toml`, only model/provider nonsecret fields
read). Thus public direct Kimi API prices alone would not establish this route's
billing terms. The model and route were not changed.

At gateway UI version 3.8.51, the Budget page
<https://omniroute.marao.sk/dashboard/costs/budget> shows daily/weekly/monthly
per-key limit fields but explicitly says **“Hard-cap policy and email alerts are
coming soon”**. Existing task client rows show no configured daily/monthly limit.
No Save Limits action was used. The UI's monetary totals are not invoices and
are not scoped Parley experiment costs; they are excluded from audit accounting.
Retained observation: runtime `billing-research/omniroute-hard-cap-ui.json`.
This is direct UI evidence, not source-level proof of every backend path. It
cannot serve as the required pre-execution USD 15 campaign ceiling.

The public Kimi membership FAQ
<https://www.kimi.com/membership/pricing?from=kfc_overview_topbar> says extra
usage can be purchased after subscription credits run out, with an Enable step,
payment-method binding and top-up. That is not evidence that this account has
extra usage disabled or that the gateway cannot select another connection.
The retained public text and URL are under `billing-research/kimi-membership*`.
No zero marginal cost, token tariff or subscription-only routing is inferred.

Provider spending remains unproven; the USD 15 total experiment budget and
900-second shared task-arm limit remain unchanged. No treatment or freeze has
started. Implementation/source-review calls remain separately accounted.

## Bounded public source follow-up

The public repository linked by the deployed docs is
`diegosouzapw/OmniRoute`. Tag v3.8.51 was absent (404 retained); the public
release/v3.8.51 tree was captured at
`cde49c93724f74604ba1e8ffcd6d8827ed9b2dce`. Eleven exact public source blobs and
SHA-256 digests are retained under runtime `billing-research/source-provenance.json`.
Deployed-byte identity is unverified.

The UI notice must not be read as proof that no backend checks exist. In that
source, `src/domain/costRules.ts:437` has `checkBudget(apiKeyId, additionalCost=0)`,
which compares recorded period spend plus the caller's additional cost with the
configured limit. `src/lib/usage/apiKeyUsageLimits.ts:478` fails closed on
previous unpriced usage for a limited window. Neither inspected function by
itself establishes a conservative reservation for all in-flight experiment
requests or provider-invoiced cost. End-to-end enforcement was not tested.

`src/lib/usage/flatRateProviders.ts` explicitly calls subscription classification
a display-only signal; budget/quota/routing keep token estimates. In
`costCalculator.ts:33`, flatRateAsZero is opt-in for display. Therefore zero on a
dashboard is not proof of zero billed experiment cost. This follow-up narrows
the observation to what was actually inspected; no gateway fix or policy change
is proposed as part of the Parley source patch.
