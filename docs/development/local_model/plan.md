# Bounded free document drafts

Status: milestone 9 implementation. Supersedes the earlier unverified local-model and free-tier estimates. Local inference is deferred. This work does not start the final automated testing milestone.

## Product contract

A small project should produce a minimal source-linked draft and PDF within five minutes on the target M1 / 8 GB machine. The draft remains useful when every cloud quota is exhausted. This is a measured target, not an availability guarantee for external services or arbitrarily large source downloads and queues.

Every profile heading is retained. Baseline sections quote relevant source declarations and user requirements, with references and specific gaps. They do not infer behavior from source declarations. Optional models improve the prose. Section provenance is `baseline`, `model` or `user`; the PDF reports their counts. A successful render is not a claim of factual correctness or academic completeness.

## Pipeline

1. Intake and evidence extraction retain their existing safe, separate-worker flow.
2. Append attributed title, specification and additional user context to the evidence view used by generation, review and export.
3. Rank facts by heading terms and evidence kind. Keep at most four excerpts per section. Build and save all missing baseline sections before network access.
4. Process ordered groups of at most four leaf sections. The largest TZ has nine primary groups. Compact each group's evidence JSON toward a 6,000-byte bound, removing lower-ranked excerpts first; a specification contributes only a matching excerpt. Output is capped at 2,400 tokens, with a 100–180-word target per section.
5. Check per-section caches before provider discovery. Validate cache data again; invalid entries are cache misses. Cache identity includes evidence, profile, language, prompt, compact algorithm, provider account scope and model. Baseline records never enter the model cache.
6. Validate each response member independently. Save good sections immediately. Failed/missing/duplicate members get at most one repair attempt, preferring another connection. Network/quota failures consume the global attempt budget but not the invalid-output repair allowance.
7. End external work at the persisted deadline, then compile all saved baseline and model sections. No background improvement starts after completion without an explicit user action.

## Providers and zero-cost boundary

Connections are optional and ordered Gemini → Groq → OpenRouter. The form accepts the exact model ID and account limits rather than shipping unverified model names. Models are checked against the provider catalog at runtime. OpenRouter must advertise JSON output and zero published pricing; generation also uses a zero `max_price` routing constraint and required parameter support. Unsupported JSON formats fail safely on other providers.

Gemini must use a project without billing and a free-tier model. Groq must use the Free Plan. The form requires these declarations when a key is supplied; the API cannot verify these account settings. The application never enables billing or silently chooses another model. Use a separate free project/account as appropriate. Keys are encrypted per job and never returned through the API. Only connected providers receive selected evidence, and only unfinished groups are sent to a fallback provider.

Official references (checked 2026-11-04):

- [Gemini compatibility](https://ai.google.dev/gemini-api/docs/openai), [billing](https://ai.google.dev/gemini-api/docs/billing), [limits](https://ai.google.dev/gemini-api/docs/rate-limits).
- [Groq compatibility](https://console.groq.com/docs/openai), [limits](https://console.groq.com/docs/rate-limits).
- [OpenRouter model catalog](https://openrouter.ai/docs/api/api-reference/models/list-all-models-and-their-properties), [free models](https://openrouter.ai/collections/free-models).

No free-tier token-rate table is treated as a promise. Set RPM/TPM no higher than the account's current limits. Defaults (5 RPM, 15,000 TPM) are conservative application settings, not provider entitlements. Keys from the same provider account may share limits; use one key per provider account with this local deployment.

## Durable budgets and scheduling

Migration 9 creates `job_draft_runs` and `provider_limits`, and adds call duration to `provider_calls`. Jet types are generated from the migrated schema. Old jobs without a draft-run record retain the legacy pipeline and credentials; new submissions use the free-draft path.

The initial deadline is database acceptance time + 240 seconds. At most 12 API calls may be admitted, including catalog checks and failures. Reservation is persisted before I/O. A crash may spend an attempt without an answer, but cannot reset its budget. One section regeneration explicitly starts a new budget and excludes all other saved sections, including other baselines.

PostgreSQL admission uses a short advisory transaction lock and row locks. It permits at most two active calls globally and one per provider across workers. RPM/TPM windows are keyed by provider/key fingerprint; input UTF-8 bytes plus output cap conservatively reserve tokens. Reservations are not refunded. API header exhaustion and Retry-After impose cooldowns. Catalog checks use the same admission mechanism.

A single job submits groups sequentially, using fallback providers while another has a cooldown; multiple workers may use different providers concurrently. This avoids redundant work and fits the upper concurrency bounds. Request timeouts are 45 seconds (catalog: 8 seconds); a 60-second reservation recovers after a worker crash. Cancellation closes requests within the existing heartbeat cycle. An unresponsive provider may continue remote processing after cancellation; local admission limits cannot control that.

## Persistence and safety

Only the active lease holder can publish sections. A model result replaces only a baseline section; existing model text and user edits survive. Provider JSON contains only section content fields and cannot set cache keys, provenance or edit flags. Unsupported citations and unsafe text are rejected. Completed artifacts remain private and use literal JSON text in the trusted Typst template.

`generation_reason`, attempts, deadline and model-improved count are visible in job responses and pages. Waiting for rate limits keeps the job running with a progress event because the worker continues to own the bounded drafting stage. It does not schedule an unbounded retry after the deadline.

## Verification

See [../verification.md](../verification.md). Do not label fixture latency as real model speed. Actual quality, account availability and token accounting remain unverified until a real free key is supplied through the website. No benchmark request may print or embed that key in documentation.
