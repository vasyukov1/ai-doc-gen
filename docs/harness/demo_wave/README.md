# Demo: three waves of simulated-user critique on ReLibrary

One run of the loop in `../simulate_wave.py` on our own product, at the three stages it went through: **mockups**
(the student app on fixtures, before the service existed), **the implementation structure** (brief, maps, user
situations, task specs) and **the built product** (screens from the scenario reports of 16–25 October). Four fictional
personas (`personas.yaml`), each run twice; one fixed directive (`directive.md`); a triage capped at five changes per
wave. Run on **29 October 2026, 23:15–23:30 CET**. Rules: `../rules/USER_SIMULATION.md`.

## The three waves

| Stage | What the personas saw | Findings | Accepted / deferred / rejected | Top 3 accepted | One rejected, and why |
|---|---|---|---|---|---|
| 1 · Mockups | 3 phone screens (390×844) of `apps/recipient` in mockup mode: [stream](../../../slides/img/s10/stage1_stream.jpg), [course](../../../slides/img/s10/stage1_course.jpg), [your map](../../../slides/img/s10/stage1_your_map.jpg) | 65 | 5 / 9 / 1 | **W1** take other students' raw activity cards ("nour · worked · finished: 3 of 5 · askedFor: an easier sheet") out of everyone's stream (all 4 personas) · **W2** one ask box with one honest audience line, no unlabelled ▾ (all 4) · **W3** on *Your map*, "Reply" becomes "Test myself" and "Unclear? Say why" (3 personas, severity 4) | **W15** (Petra, the teacher): a class overview on *Your map* naming the students who are behind — breaks **I3** ("nothing grades them silently"); her real need, seeing what students struggle with, stays as per-part questions to the library (I4) |
| 2 · Structure | [structure.md](stage2_structure/structure.md): `PLAN.md` §1, the measured counts ([diagram](../../../slides/img/s10/stage2_structure.png)), all 103 situation titles, 6 situations in full, 1 task spec (≈ 3,400 words) | 63 | 5 / 15 / 2 | **W2** put the audience on every send button and give the three library entries one name (3 personas, 11 findings) · **W1** one top-down order for the course screen, opening with an unpersonalised "Due" strip (2 personas) · **W4** a teacher situation "Preview a day as students see it" (Petra, severity 4 in both runs) | **W21** (Rune): a "Skip basics" switch that removes basic parts from his stream — breaks **I5** (personalisation "never changes … what the course requires") and **I4** (only the teacher removes parts); kept: per-area choices that shorten basics while every part stays, with its check (W12, deferred) |
| 3 · Product | 4 screens of the built app: [day one](../../../slides/img/s10/stage3_day_one.jpg), [section test](../../../slides/img/s10/stage3_section_test.jpg), [what you missed](../../../slides/img/s10/stage3_what_you_missed.jpg) (phone), [your week](../../../slides/img/s10/stage3_your_week_laptop.jpg) (laptop) | 64 | 5 / 11 / 2 | **W5** "What you missed" as one numbered column in taught order (all 4 personas) · **W4** the "About: …" chip on the ask box only when the student chose it (all 4) · **W3** one true sentence under "Ask the library" instead of "Just you and the library" (3 personas) | **W18** (Rune): sort "What you missed" by what Homework 1 needs, hide the other 16, let "I know this" remove parts — breaks **I1** (no "re-sorted feed") and **I5**; kept: a "for Homework 1" label on rows that stay in taught order (W11, deferred) |

Each stage folder has `findings.json` (every finding with its persona and run, and the spread between runs),
`wave.md` (the decisions, what must not change, the steps to walk again), and the exact `prompts/` and `responses/`.

**Spread between the two runs** (share of screen + heuristic pairs both runs found): 0.15–0.56 on the mockups,
0.36–0.67 on the structure, 0.36–0.45 on the product. Half or more of what one run finds, the other run does not: run
each persona at least twice (U7), and treat a finding as a hypothesis to check on the next walk (U8).

## What the three waves show together

- **One question survived all three stages: who reads what I ask the library?** The structure promises "Just you
  and the library" (`US.C.70.11`); the built app also says "your teachers can read it", and so does the directive
  (I4, taken from the app). The personas raised it in the written situations (stage 2, W2), in the mockup's ask box
  (stage 1, W2) and in the built app, where two screens print opposite sentences (stage 3, W3: Thao in both runs, Ayla,
  Petra). A wave on the structure would have settled it before it reached the mockup and the product.
- **The directive held at every stage.** All five rejections are requests that would help one persona today: a ranked
  list of who is behind (Petra, three stages), a stream without the basics or sorted by the homework (Rune). Each was
  rejected with the clause it breaks and says what can still be solved inside the directive; three point at a deferred
  change in the same wave, two only describe it, which U10's check would flag.
- **The cap decides what waits.** Bare bottom-bar icons were found on the mockup and on the built app: stage 1
  accepted labels (W5), stage 3 deferred the same change behind five bigger ones (W14). A deferred finding comes back
  until it is built.
- **Done once before, by hand.** `.pack_documentation/suggestions/` read the same structure by 25 audience segments and
  ranked 11 gaps in three tiers while leaving the plan untouched; its tier-3 "no way to ask the library anything" is now
  a box on every student screen. Stage 2 is the same move in minutes, with a directive deciding instead of tiers.

## How it was run

Mockup mode started: the student app on fixtures with its own dev server (`VITE_DATA` unset, port 5291, stopped after),
screenshotted with Playwright at 390×844 on `/`, then the *Shelf* tab (`/kernel/k_sharding`) and the *Your map* tab.
Stage 3's screens are copies of the scenario reports' images (`24_student_week_wk02/06/16`, `24_student_week_laptop_wk29`
in `.pack_documentation/documents/scenario_reports_2026-10-16/img/`), resized to 900 px. Stage 2's `structure.md` was
assembled from `PLAN.md` §1, `.pack_documentation/{architecture,requirements,user_situations,plans}` and a count of the
test files; the counts are in its section 2.

```bash
cd services/objects/simulation/courses/ds402/plan/examples/session-10
I=../../slides/img/s10
W="--directive demo_wave/directive.md --personas demo_wave/personas.yaml --runs 2"
python3 simulate_wave.py $W --screens $I/stage1_stream.jpg $I/stage1_course.jpg $I/stage1_your_map.jpg \
    --out demo_wave/stage1_mockups
python3 simulate_wave.py $W --text demo_wave/stage2_structure/structure.md --out demo_wave/stage2_structure
python3 simulate_wave.py $W --screens $I/stage3_day_one.jpg $I/stage3_section_test.jpg \
    $I/stage3_what_you_missed.jpg $I/stage3_your_week_laptop.jpg --out demo_wave/stage3_product
```

**Models:** critique `claude-sonnet-5-5` (24 calls: 4 personas × 2 runs × 3 stages), triage `claude-opus-5-5`
(3 calls). **27 model calls in all.** All answers are saved, so re-running these commands now costs nothing.

**How the calls were made.** `ANTHROPIC_API_KEY` was not set in the session that ran this, so each command ran in
three passes: with `--dry-run` it wrote every prompt to `prompts/`; each prompt was then answered by a fresh agent
pinned to the model above, which read only that prompt file and its screenshots and saved its JSON answer to
`responses/<call>.json`; the same command run again read the saved answers, built `findings.json`, wrote the triage
prompt and, once the triage answer was saved, `wave.md`. This is the script's by-hand route (`../ux_critic.md`), with
the model doing the pasting. Two differences from the API route: the answers' JSON shape was checked after the fact
(all 27 valid; every finding lands in exactly one triage change) instead of being enforced by structured outputs, and
each agent read the prompt's system part as instructions inside its own harness rather than as the API's `system`.
With a key set, the same commands do all of it in one pass.
