# Simulated user journey

Simulate one user from a specified candidate niche completing one real documentation job. Ask for a plausible path through these steps:

1. Discovers the tool and decides whether to trust it with code.
2. Supplies a public GitHub link or ZIP and an optional description.
3. Selects a profile and enters required metadata.
4. Reviews included files and corrects the selection.
5. Sees extracted facts and missing information.
6. Reviews generated sections and fixes errors.
7. Downloads editable Typst and PDF, then checks acceptance against their actual rubric.

Return JSON with `stopped_at`, `steps`, `first_value_moment`, `unresolved_risk`, and `alternative_used`. Do not invent a payment probability. Repeat across different tasks and constraints only to locate likely failure points. Simulated journeys cannot estimate conversion or demand.
