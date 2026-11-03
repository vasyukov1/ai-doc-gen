# Niche brief — AI Doc Gen

## The product

One sentence: a web application that reads a GitHub repository or uploaded source archive, extracts verifiable project facts, drafts a selected documentation profile, and delivers editable Typst plus PDF.

The current alternatives include manually reading code and writing documents, using a general AI chat with pasted files, university templates, and documentation platforms. The intended advantage is a complete document artifact with traceable code evidence, section-level review, and reproducible PDF output. That advantage is a hypothesis until tested with users.

## Decision this analysis serves

The product owner selected students preparing coursework or thesis projects as the first audience, with an existing academic document set adapted to the website and generic technical project documentation as the first two families (2026-11-02). Public GitHub repositories and ZIP archives are the initial inputs; users choose Russian or English for output. The initial generation approach is a hosted LLM API with checks and caching. Go and PostgreSQL will run asynchronous jobs without Kafka for the first release. English template review, pricing, and launch channel still need validation. No credible launch threshold can be set from the available evidence.

## Candidate niches

| Niche | User and job | Likely access route | Main uncertainty |
|---|---|---|---|
| University students | Turn a course or thesis codebase into an institution-accepted project report | Student communities and course instructors | Institution-specific rules and academic integrity policies |
| Freelance developers | Deliver a handover or technical project report to a client | Developer communities and agencies | Willingness to pay for a one-off artifact |
| Small software teams | Produce architecture, requirements, or operations documentation for existing code | Team leads and developer channels | Approval, privacy, and integration needs |
| Technical writers or consultants | Draft evidence-backed documentation for a client repository | Professional communities | Fit with their review workflow |

## Evidence we have

- A prototype codebase demonstrates local Typst templates, section generation through hosted LLM APIs, and PDF compilation. It is implementation evidence, not demand evidence.
- Official product pages confirm that several established platforms offer documentation authoring, AI functions, and publishing. See [competitors.csv](competitors.csv).
- Official standards pages confirm the scope of candidate ISO/IEC/IEEE standards. See [standards_and_templates.md](standards_and_templates.md).

## Evidence we need next

1. Interview students about their latest real documentation task, artifact, deadline, review process, and alternative used. The other candidate niches remain expansion options. The sample size is a research target, not an observed result.
2. Collect a small, permissioned set of repositories and corresponding human-approved documents for a blinded quality comparison.
3. Test a clickable workflow and measure completion, correction effort, and willingness to use the output. Do not infer willingness to pay from a simulated persona.
4. Identify specific university or company formatting rules before labeling any profile as compliant.

## Evidence rule

Every quantitative claim needs a source, date, definition, and kind (`open`, `partial`, `research`, `assumption`, or `simulated`). Report uncertainty. Simulations may reveal likely failure points but cannot establish demand or conversion.
