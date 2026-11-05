# User simulation rules (U1–U13)

Read these before you let simulated users critique a screen, or triage their findings into a wave of changes.
Each rule has an id, the rule in one line, why it exists, and a check that you or your agent can run.
When the agent applies a rule, it names the id, e.g. "(U9)".

The loop: freeze the directive → write personas → show them the stage you are at → each critiques, twice →
triage against the directive → build the wave → walk again. Then D9 before you call it done.

| Step | directive | personas | what they see | critique | triage | the wave | walk again |
|---|---|---|---|---|---|---|---|
| Rules | U1, U2 | U3, U4 | U5 | U6, U7, U8 | U9, U10 | U11 | U12, U13 |

Tool: `simulate_wave.py` in the course's `examples/session-10/`, with its prompts `ux_critic.md` and `wave_triage.md` and a worked run in `demo_wave/`.

---

### U1. Write the directive before the first wave, and freeze it
One paragraph (who it serves, what they get, in what form, who stays in control) and 3–6 invariants, in `directive.md`.
- **Why:** Simulated users ask for whatever helps them most today; without a fixed directive the product drifts one reasonable request at a time.
- **Check:** `git log -- directive.md` shows no commit between the first and the last wave; each invariant is a sentence you could test on a screen.

### U2. The directive belongs to the human
The agent may propose a new directive as a separate question with its default (D10); it never edits it inside a wave.
- **Why:** A wave that changes its own yardstick cannot be compared with the wave before it.
- **Check:** Every change to `directive.md` is its own commit, made by the human, outside any wave.

### U3. Choose personas that disagree, and invent them
3–5 personas that differ in device, time, knowledge, language and role; at least one is not your main user. Fictional names only (A13).
- **Why:** Four similar personas are one opinion counted four times, and a caricature of a real person is not something you publish.
- **Check:** No two personas share device and goal; no name belongs to a real student, client or colleague; the teacher or admin is in the cast.

### U4. Give each persona mistakes, not opinions
Every persona has a goal, a device, what they know, and two or more typical mistakes ("taps the biggest button first").
- **Why:** Mistakes take a simulated user off the happy path, which is where the problems are; "likes clean design" finds nothing.
- **Check:** `personas.yaml` has a `mistakes` line per persona; some findings start from one of those mistakes.

### U5. Show the stage you are at, and only what exists
Mockups: screenshots of the app on fixtures. Structure: brief, user situations and task specs as text. Product: screenshots of the built app, phone and laptop.
- **Why:** The earlier the stage, the cheaper the change; a persona shown a feature the product lacks critiques another product.
- **Check:** Every input comes from the app or the documents at this commit; no screenshot shows a real person's name, an e-mail or a key (A13).

### U6. One persona per conversation, and every run starts fresh
Never put two personas in one chat or one prompt; a second run never sees the first.
- **Why:** A persona that has read another's findings repeats them, and the agreement you then measure is your own artefact.
- **Check:** One prompt file and one response file per persona per run; no response mentions another persona.

### U7. Run each persona at least twice, and report the spread
Report findings per run and the share of (screen, heuristic) pairs both runs found.
- **Why:** A model's answer varies between runs; a finding seen in both runs is more stable, not more true.
- **Check:** `findings.json` has a `spread` per persona, and `wave.md` quotes it.

### U8. A finding is a hypothesis, not a measurement
Each finding names the screen, what was tried, what happened, the problem, a severity 1–4, one Nielsen heuristic by name and one suggestion.
- **Why:** Simulated users find what a heuristic review finds; they cannot say how many real people would hit it (A6: label it simulated).
- **Check:** Every finding has all seven fields; nothing in the report says "N% of users" about simulated users.

### U9. Triage against the directive: ACCEPT, DEFER or REJECT, with the clause
Deduplicate first; then decide every change, with a one-line reason that quotes the directive or names an invariant.
- **Why:** The directive is what stops a good-sounding request from turning the product into a different one.
- **Check:** Every finding id lands in exactly one change; every decision's reason names a clause or an invariant.

### U10. A REJECT keeps the user's problem
For each REJECT, say which clause it breaks and which part of the user's problem can still be solved inside the directive.
- **Why:** The request breaks the directive; the pain behind it is still real, and losing it loses a user.
- **Check:** Each rejected change names an alternative, and that alternative is an ACCEPT or a DEFER in the same wave.

### U11. Cap the wave at five
At most five ACCEPTs per wave: several personas, high severity, cheap to build. Everything else is DEFER, not deleted.
- **Why:** A wave must be small enough to build and walk again before the next one, or you cannot tell which change fixed what.
- **Check:** `wave.md` has ≤ 5 accepted changes (the script demotes the rest); the deferred list is carried into the next wave.

### U12. Build the wave, then walk again: same personas, same steps, one invariant
After building, re-run the same personas on new screenshots of the steps in "Walk again", one of which checks an invariant still holds.
- **Why:** A change that fixes one persona's problem can create another's, or quietly break the directive; only walking again shows it.
- **Check:** The next wave's screenshots post-date the wave's commits; fixed findings do not recur; the invariant step passes (then D9).

### U13. Keep it reproducible and cheap
Save every prompt and every model answer; re-runs reuse saved answers; log the command, models, date and number of calls.
- **Why:** A wave nobody can re-run is an anecdote, and a loop that costs a lot gets run once.
- **Check:** `prompts/` and `responses/` exist for every stage; the README gives the exact command, the model ids, the date and the call count.
