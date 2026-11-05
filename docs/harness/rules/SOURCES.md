# Sources: which internal rules each student rule condenses

For the teacher. File names only; the sections are named as they appear in each file.
Analytics rules come from the analytics rules folder of the knowledge-base repository (and its
`.windsurf/rules/`); development rules come from this repository (`.windsurf/rules/`).
`11d_evidence_discipline.md` is identical in both repositories.

## Analytics rules

| Id | Rule (short) | Condenses |
|---|---|---|
| A1 | Start from the reader's decision | `general_news.audience.md` (whole file: who we write for, what earns a slot, lead with the "so what"); `general_news.rules.md` § The audience comes first, § Section necessity (the four-question test); `sales_scans.rules.md` § The offer comes first |
| A2 | One workflow, with a human gate | `general_news.workflow.md` (Phase 1, Human gate, Phase 2, Invariants); `sales_scans.workflow.md` (Phase 1, Human gate, Phase 2); `sales_scans.rules.md` § Human gate; `general_news.rules.md` § Discussion state (written before any reader-facing content) |
| A3 | Every number: source, date, kind | `general_news.rules.md` § Source links (every number traces to a link or is labelled as our own arithmetic; vendor figures marked); `quarterly_reports.rules.md` § Required sections ("The numbers": value, period, source, confidence); `sales_scans.offer.md` § Source discipline (person-level facts marked unverified) |
| A4 | Link each claim to the original, and open it | `general_news.rules.md` § Source links (what counts as a key statement, the form, density), § When a primary source blocks you ("never publish analysis of a document you have not opened"); `general_news.workflow.md` step 8 (find and verify the initial source); `quarterly_reports.rules.md` § Report conventions (two-tier citation; a homepage alone is not evidence) |
| A5 | Grade evidence; state flatly only A or B | `general_news.rules.md` § Discussion state & evidence analysis (§2 grades A–E; vendor self-reporting never clears the bar), § Source links (tier table); `general_news.workflow.md` step 9a (`intel/sources.md`: per-claim confidence) |
| A6 | Keep fact, analysis and speculation apart | `quarterly_reports.rules.md` § Binding constraints (separate the three registers; never fabricate a speaker position); `general_news.rules.md` § Source links (analysis stays unlinked), § Podcast short takes (wild takes labelled as speculation, never our opinion laundered as someone else's) |
| A7 | Search twice; find independent sources and critics | `general_news.workflow.md` step 4 (second pass excluding the first), step 9a (`intel/contradictions.md`: hunt disagreement, who contests whom); `sales_scans.rules.md` § Research rules for step 4 (second pass for the awkward facts; primary over aggregator); `general_news.rules.md` § Discussion state §2 (two independent sources that do not derive from each other) |
| A8 | Say what each source cannot see | `sales_scans.rules.md` § Phase 1 step 2 (what the group's composition biases); `general_news.rules.md` § Discussion state §4 (what the evidence cannot settle; confounders: definitions, horizons, populations); `quarterly_reports.rules.md` § Report conventions (name the evidence gaps); `sales_scans.workflow.md` step 5 ("no trigger found is a valid result"); `sales_scans.offer.md` § Source discipline ("a guess is not") |
| A9 | Argue the other side once, honestly | `general_news.rules.md` § Contrarian analysis (mechanism, not mood; steelman once); `sales_cases.rules.md` § Working a case step 5 (test the client's hypothesis and ours); `quarterly_reports.rules.md` § Required sections (contrarian view + steelman) |
| A10 | Confidence, and what would change your mind | `general_news.rules.md` § Phase 2 (investigation: what would change our mind), § Discussion state §4 (provisional read with a confidence level; what would resolve it); `sales_scans.offer.md` § Budget bands (confidence on every estimate); `sales_cases.rules.md` § Working a case step 6 (what would kill it); `quarterly_reports.rules.md` § Required sections (speculation with the observable that confirms or kills it) |
| A11 | Date the work; never rewrite the past | `quarterly_reports.rules.md` § Report conventions (no content dated after the report date; week marker on drafts; flag a source's vintage), § Binding constraints (`knowledge_cut`); `sales_scans.rules.md` § Cadence and re-scanning (new directory, new key, never rewrite an old scan); `general_news.rules.md` § Back-references (when we got it wrong); `regional_monitoring.rules.md` § Directory structure (dated archive copies) |
| A12 | Make every number reproducible | `sales_cases.rules.md` § Working a case step 1 (client sources copied verbatim and numbered); `general_news.rules.md` § Short episode transcripts, § Vertical slides (generated, never hand-authored; never edit the book directly); `quarterly_reports.rules.md` § Report conventions (a count is accepted only if its items are named, with a count check) |
| A13 | Keep personal data out | `sales_scans.rules.md` § What this pipeline does not do; `sales_cases.rules.md` § What a case does not do; `sales_scans.offer.md` § Source discipline (not researched by design); `quarterly_reports.rules.md` § Binding constraints (public-facing filter; no client names); `general_news.rules.md` § Key-actor portraits (never generate a photo of a real person); also `26_simulated_user_journeys.md` (personas are fictional) and `27_push_check.md` (scan for secrets before pushing) |
| A14 | Write the guardrails before you publish | `general_news.rules.md` § Discussion state §5 (publishing guardrails: safe, attribute, never; copied verbatim downstream), § Source links (check before publishing; every link resolves); `general_news.workflow.md` § Invariants |
| Grades table | A–E | `general_news.rules.md` § Discussion state §2 and the tier table in § Source links |

## Development rules

| Id | Rule (short) | Condenses |
|---|---|---|
| D1 | Requirements before code | `03-requirements-rules.md` (product and functional requirements; every product requirement linked to a user situation; the generation protocol: propose requirements when information is insufficient); `02-pack-documentation-structure.md` § User situations (each action unambiguous); `12_architecture.md` § New component documentation (before any code) |
| D2 | One id per thing, one index | `02-pack-documentation-structure.md` § Indexing scheme, § Index files, § General rules; `01_organization.md` § General rules (the ID in the file name makes an object findable) |
| D3 | Plan in small tasks first | `13_planning_workflows.md` § Feature planning (plan file before any work; task contract: inputs, process, outputs, tests, verify, handoff); `04-testing-planning-rules.md` § Rules (tasks of 20–30 minutes); `11_core_operations.md` § Planning triggers (tick immediately, never batch) |
| D4 | Tests planned from user situations | `04-testing-planning-rules.md` § Test maps, § Rules (every test covers a product requirement); `14e_usat_testing.md` (a test for every user action in every user situation); `13_planning_workflows.md` § Test writing planning (failing tests first); `11a_test_execution_mandate.md` § Tests define correctness; `14_testing_strategy.md` § New feature tests (smoke, contract, end-to-end) |
| D5 | Mockup data first | `05-development-style.md` § Component functioning modes (mockup mode: real API calls, mocked data, no hard-coded static data); `09-mockup-data-rules.md` (typical, empty and error files; realistic values); `21-project-conventions.md` § Architecture (mockup-first) |
| D6 | Components with explicit contracts | `12_architecture.md` § Service boundaries, § Dependencies; `02-pack-documentation-structure.md` § Integration requirements; `14c_contract_testing.md`; `18_change_coordination.md` § Boundary verification; `15_code_quality.md` § System design (no hard-coded keys or environment values) |
| D7 | Run tests yourself; fix code, not tests | `11a_test_execution_mandate.md` (green definition; no avoidance; no test simplification); `14_testing_strategy.md` § Test change policy, § Test execution discipline; `15_code_quality.md` § Anti-stub policy; `11_core_operations.md` § Quality mandate; `13_planning_workflows.md` § Capability integrity; `14b_ui_and_complexity.md` (the agent runs interface tests, never the user) |
| D8 | No claim without evidence | `11d_evidence_discipline.md` (whole file); `11_core_operations.md` § Stop-ship checklist |
| D9 | Simulated users before done | `26_simulated_user_journeys.md` (whole file); `14e_usat_testing.md` (one action at a time, with a screenshot) |
| D10 | Stated defaults; stop only for the irreversible | `06-autonomous-uncertainty.md` (self-resolve, research, switch and document; never block on one uncertainty); `25_coordination.md` rule 8 (questions carry their default; wait only on what is irreversible or outward-facing) |
| D11 | Commit each step, then push | `15_code_quality.md` § Committing changes; `18_change_coordination.md` § Coordination workflow (commit frequently); `27_push_check.md` (push, scan for secrets, say why anything stays local) |
| D12 | Docs in the same commit as the code | `17_documentation.md` (update docs whenever code changes); `02-pack-documentation-structure.md` § Integration test mandate (contract change and test in the same commit); `09-mockup-data-rules.md` § Shared entity fixtures (same commit); `18_change_coordination.md` § Coordination workflow (before marking done) |
| D13 | One task, one branch; stage only your files | `07-workspace-isolation.md` (isolated workspaces; never modify another agent's workspace); `25_coordination.md` (a worktree of your own; work cut from another team's active branch adds new files only); `27_push_check.md` § Never (never stage files kept out of git) |
| D14 | Coordinate several agents | `25_coordination.md` rules 1–7 (read the board, claim before you write, never change a path another owns, reserve names, lock shared resources); `18_change_coordination.md`; `19_schema_coordination.md`; `20_merge_coordination.md` |
| D15 | End with the next exact command | `13a_session_continuity.md` (resume point with the exact next command); `13_planning_workflows.md` § Session continuity; `25_coordination.md` rule 9 |

## The pack itself

| Part | Condenses |
|---|---|
| `AGENTS.md`: product directive | `21-project-conventions.md` § What this repo is, § The rule that matters most; `03-requirements-rules.md` (the human approves requirements); `17_documentation.md` § Breaking changes (explicit approval) |
| `AGENTS.md`: "Read when" table | `00_aNAVIGATION_TREE.md` § Task classification; `00_NAVIGATION.md` (the analytics entry point) |
| `AGENTS.md`: project facts | `21-project-conventions.md` § Commands; `11_core_operations.md` § Environment |
| README: always on vs on demand | the `trigger:` front matter of every file in `.windsurf/rules/` (`always_on`, `model_decision`, plus our own `context` and `semantic_auto`); `10_rule_loading_triggers.md` (always-on, semantic and context triggers; 30–40% saving); `01-pack-triggers.md` |
| README: the navigation file | `00_aNAVIGATION_TREE.md` ("load this file first"; loading strategy); `10_rule_loading_triggers.md` § Design principles (intent over syntax; prefer false positives) |
| README: how to tell they work, how to change | `11_core_operations.md` § File size limits (rule files ≤ 293 lines); `10_rule_loading_triggers.md` § Maintenance (monthly review); `26_simulated_user_journeys.md`, `27_push_check.md` (rules that quote the request that created them); `25_coordination.md` (one canonical copy, copied unchanged) |

## Read but not condensed

Specific to our stack, not to a student project: `08-pack-asset-standards.md` (document templates),
`14a_test_infrastructure.md` and `23_docker_mcp_usage.md` (containers), `14d_ui_performance_optimization.md`
(browser-test speed), `16_navigation_reference.md` (our search tools), `22_article_processor_integration.md`
(one service's integration guide), `24_tier0_checklists.md` (a shorter view of rules condensed above);
`general_news.rules.md` channel formats (posts, transcripts, slides, carousels, media plan timing),
`robot_scenarios.rules.md`, and `sales_scans.offer.md` § Products, § Proof, § Public co-funding.
