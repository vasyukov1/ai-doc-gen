# Employee rewalk — fictional Lena

Reader: wave reviewer deciding whether the accepted changes preserve an honest architecture-document workflow. Cut-off: 2026-11-05. Simulated review, not observed user research (U8, A6). One isolated review invocation; no provider APIs or keys used. The runtime model identifier was not exposed to this reviewer.

Evidence: [laptop](../screens/after-laptop.png), [phone](../screens/after-phone.png), [employee status](../screens/employee-status.png), [walk log](../checks/walk.log), and current `web/templates/{index,job,sections,document}.html`, `web/src/usage.ts`. Screenshots establish visible copy/layout, templates establish conditional UI intent, and the supplied log reports a fixture run; none demonstrates actual user comprehension or provider billing. This review did not perform browser interactions.

| Step: laptop then phone | Observed evidence | Simulated assessment |
| --- | --- | --- |
| Lena prepares to upload confidential source | Both upload screenshots place server processing, selected cloud context and sanitized-project guidance above the source controls. The submit note repeats the boundary. | PASS for disclosure placement. Hypothesis: Lena can correct the mistake of interpreting browser-private access as local-only processing before uploading. |
| Choose architecture | Both upload screenshots expose Architecture; the template maps it to arc42. Employee status explicitly names architecture documentation and arc42; the log reports an outline result. | PASS for discoverability and observed result. Selecting the shortcut itself was not exercised here. |
| Mistake: interpret full progress as quality certification | Employee status displays full section availability beside an explicit non-quality explanation, outline-only warning and zero model-improved sections. The PDF card warns that compilation does not verify requirements, code behavior or test results. | PASS for the no-quality-certification invariant in the supplied phone view and shared template. Laptop status rendering was not supplied. |
| Mistake: infer free generation from absent usage | Status says no provider calls recorded, and explicitly separates missing usage from zero cost and billing. `usage.ts` supports refresh, failures and capped-report caveats. | PASS for honest wording. Actual provider counts, billing and refresh interactions remain unverified. |
| Read intermediate text, then edit and rebuild | Status exposes Read available sections and explains early access. The sections template links Review and rebuild, and warns saving does not update the PDF. The document template requires a review checkbox before compilation when rendering is available. | PASS for visible entry point and source-level instructions; live intermediate reading, saving and rebuilding were not exercised or shown in these screenshots. |

The explicit privacy boundary, purpose shortcut, quality disclaimer, usage caveat and intermediate-text entry point address the intended confusion at the relevant screens. This is a hypothesis about improvement, not a measured before/after effect.

Remaining hypotheses (defer; no new accepted changes):

- **Employee-after-H1 — upload / Error prevention / severity 2 (assumed).** Lena tries to skim directly to upload; the nearby disclosure is visible but small. She might miss it on a phone. Suggest a later comprehension check before increasing visual emphasis. Confidence: medium; falsifier: representative employees consistently explain server/cloud processing before selecting files.
- **Employee-after-H2 — status / Visibility of system status / severity 2 (assumed).** Lena skims the full progress bar and compiled badge; despite nearby warnings, these may dominate her reading. Suggest a later task asking what is safe to share, preserving the existing disclaimer. Confidence: low; falsifier: readers correctly distinguish availability, compilation and verified quality without prompting. Strong contrary reading: repeated explicit outline warnings already make the limitation sufficiently clear.

No certification claim found in the reviewed screen/template paths. Actual comprehension, laptop status layout, model-connected generation, costs and edit-to-rebuild behavior remain outside this review's evidence.
