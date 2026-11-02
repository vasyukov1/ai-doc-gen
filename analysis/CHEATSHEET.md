# Niche and product decision cheat sheet

## The rule

Use [evidence_ledger.csv](evidence_ledger.csv) for every quantitative market, pricing, or quality claim. Public product pages establish offers, not demand. Simulated users establish neither demand nor willingness to pay.

## Thirteen steps

| Step | Action | Done when |
|---|---|---|
| 1 | Choose the product decision and candidate niches | [niche_brief.md](niche_brief.md) names both |
| 2 | Map each niche's latest documentation job | Interview notes include an actual artifact and deadline |
| 3 | Inventory rules and standards | Exact profile basis, edition, locale, and rights are recorded |
| 4 | Scan alternatives, including manual work | [competitors.csv](competitors.csv) has primary sources |
| 5 | Collect public conversations | Deduplicated export and source limitations are recorded |
| 6 | Extract recurring pains and desired outputs | Manual review checks extraction accuracy |
| 7 | Map the user journey | The first value moment and likely stop points are explicit |
| 8 | Test a clickable prototype | Users submit real code and request a real artifact |
| 9 | Build a consented evaluation set | Source repos and reviewed target documents are available |
| 10 | Measure the generation baseline | Support, completeness, compile success, edit time, latency, and cost are reported |
| 11 | Improve extraction and section generation | Same held-out set shows a measured gain |
| 12 | Compare local and hosted models | Quality, privacy, throughput, and total cost are compared |
| 13 | Choose launch niche and price | Real behavior supports the choice; uncertainties remain visible |

## Immediate decisions

The owner selected students, the existing academic document set and generic technical documentation, public GitHub and ZIP inputs, Russian or English output, hosted LLM generation with checks and caching, and asynchronous Go/PostgreSQL jobs. The technical specification is uploaded or entered through a form; the English academic title page preserves the existing structure with translated text. The current stage is analysis only. Use [interview.md](interview.md) to resolve the remaining editing and deployment constraints before implementation. The future job contract is in [product_analysis.md](product_analysis.md).

## Stop conditions

Pause a profile's release if its required source rules are unavailable, users cannot correct output, unsupported claims remain in the document, or the editable Typst bundle cannot reproduce the PDF. Pause a commercial launch decision if no real user has completed a task with a real repository.
