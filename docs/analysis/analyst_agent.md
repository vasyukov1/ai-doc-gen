# Instructions for an analyst agent

Read [niche_brief.md](niche_brief.md) and [evidence_ledger.csv](evidence_ledger.csv) first.

1. State the single decision being supported and the output schema.
2. Write an analysis plan before processing data: sources, filters, joins, expected row counts, metrics, uncertainty method, and checks.
3. Validate data counts, duplicate keys, missing values, dates, units, and source accessibility. Report every exclusion.
4. Do not invent numbers, quotes, users, prices, or standard requirements. Label assumptions and simulations explicitly.
5. Keep vendor marketing claims separate from independently observed user outcomes.
6. For a standards claim, record the standard identifier, edition, exact accessible text or approved rubric, and what the template actually implements.
7. For a product-quality claim, use held-out repositories, a reviewer rubric, and an appropriate interval; report failures and correction effort.
8. Return a concise decision, supporting evidence, uncertainty, next test, and ledger updates. If evidence is insufficient, say so.
