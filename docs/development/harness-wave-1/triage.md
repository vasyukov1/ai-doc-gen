# Wave 1 triage

Simulation input: frozen implementation, not executed browser journeys. Six fresh-context Codex subagents, two independent runs per fictional persona. No project provider keys or external simulation API were used. Primary-agent browser checks separately establish observed behavior. Reviewers wrote JSON responses; their findings remain hypotheses (A6, U8).

## Accepted changes (U9, U11)

1. R1 / I4: correct source/specification mistakes before upload; visible file removal; keep entered data. Covers beginner source-choice findings. Validation limits come from existing server contracts.
2. R2 / I2: preserve tables on prose save, inline save, unsaved-edit guard and explicit regenerate confirmation. Covers maintainer table-loss and dirty-edit findings. Confirmation mitigates accidental regeneration, but does not add version history.
3. R3 / I3: label PDF compilation separately from content quality; show model/baseline/gap counts and reachable review/continuation actions. Covers beginner baseline and reviewer completed/retry findings.
4. R4 / I4: visible stale status, bounded polling and manual refresh. Added from direct code inspection: catch previously discarded every polling error.
5. R5 / I1, I4: explain document types and key-free output, bilingual profile guidance, accurate server/provider processing disclosure. Covers privacy and partial beginner terminology findings.

## Deferred, not discarded

- Full Russian interface localization (beginner): R5 provides focused bilingual guidance; translating every page and error consistently is a separate wave. It serves I4 but is larger than this correction.
- Regeneration history / compare-and-accept (maintainer): requires a durable version model and explicit replacement policy. For this wave, retain explicit confirmation and do not present regeneration as reversible. This remains a data-loss risk after a deliberate confirmed regeneration.
- Save-all, cross-session draft recovery, and table editing/removal: preserving existing table data and saving individual sections without navigation addresses the observed loss first. No browser storage of code, keys or drafts is introduced.

Confidence: high for observed code paths, medium for predicted user impact. Falsifier: manual phone/laptop checks cannot reproduce the loss/confusion or the new flow prevents completion. No claim about the percentage of real users is made.
