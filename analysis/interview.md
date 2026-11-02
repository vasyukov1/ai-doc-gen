# Product-owner interview and decision log

Decisions below reflect the product owner's answers on 2026-11-02. Open items remain questions, not implicit approval of a default.

| Decision | Selected direction or recommendation | Alternatives and trade-off | Status |
|---|---|---|---|
| First audience | Students preparing coursework or thesis projects | Teams remain a later expansion audience | Selected by owner, 2026-11-02 |
| First document families | Existing academic document set adapted to the website, plus generic technical project documentation | The existing set has five document variants; English content needs its own review | Selected by owner, 2026-11-02 |
| Initial repository access | Public GitHub URL and ZIP | Private GitHub access adds OAuth, scopes, and retention obligations | Selected by owner, 2026-11-02 |
| Output language | Russian or English, selected for each job | The current academic layout and fixed text are Russian; English needs translated and checked template content | Selected by owner, 2026-11-02 |
| Technical specification input | Upload a file or complete a web form | Generating it from code alone would invent intended requirements | Selected by owner, 2026-11-02 |
| English academic title page | Translate fixed text while preserving the existing title-page structure | A neutral title page would change the selected academic format | Selected by owner, 2026-11-02 |
| Current delivery stage | Analysis only | Website implementation and template migration begin after a later instruction | Selected by owner, 2026-11-02 |
| Generation | Hosted LLM API with checks and caching | Local model and fine-tuning remain later comparisons | Selected by owner, 2026-11-02 |
| Job processing | Go workers with PostgreSQL job leases and status polling | Kafka only if later workload or integration requirements justify it | Selected simplification, 2026-11-02 |
| Deployment | Single-region EU-hosted beta | Other regions may fit the audience; data residency affects vendor selection | Open |
| Editing | Section editor plus Typst bundle download | Browser Typst editor offers more control but increases implementation scope | Open |
| Pricing | Defer until observed use | Per-document, credits, subscription, or team seat pricing can be tested | Open |

## Questions that need user answers

1. Can you provide repositories and accepted outputs that can be used for evaluation? Are those materials permitted to be shared with an LLM provider?
2. Is there an EU-hosting or no-third-party-model requirement for uploaded code?
3. Should users edit sections in the browser, edit Typst locally, or both?
4. What is the target demonstration date and who will judge whether the result is acceptable?

## Decision rule

Record each answer here with date, owner, and reason. Update scope and quality gates after the choices are made.
