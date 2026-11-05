# Analytics rules (A1–A14)

Read these before any work that collects data, computes a number, or writes a finding or a recommendation.
Each rule has an id, the rule in one line, why it exists, and a check that you or your agent can run.
When the agent applies a rule, it names the id, e.g. "(A3)".

The rules follow the order of the work:

| Stage | question | sources | evidence ledger | analysis | findings | recommendation, publishing |
|---|---|---|---|---|---|---|
| Rules | A1, A2 | A4, A7, A8 | A3, A5, A12 | A6, A9 | A10, A11 | A14; A13 throughout |

A3 and A5 are the session 9 ledger, written as rules. Its `kind` column maps like this: `open` and `partial` are
observed, `research` is reported, a number you compute is estimated, `assumption` is assumed. A `simulated` answer
counts as assumed: it may rank options, never set a level.

---

### A1. Start from the reader's decision
Name who will read the analysis and the decision it must inform, before you collect any data.
- **Why:** The reader decides what earns a place; without one, nothing tells you what to cut.
- **Check:** The first lines name the reader and the decision; every section passes "would they miss it if it were gone?"

### A2. Follow one workflow, with a human gate
Go question → sources → evidence ledger → analysis → findings → recommendation, and stop for a human before the expensive part.
- **Why:** A fixed order stops conclusions being written before the evidence exists; the gate stops effort on the wrong question.
- **Check:** The work shows each stage in this order, and the plan records where a human chose what to pursue.

### A3. Give every number its source, date and kind
Every number carries where it came from, the date it is true as of, and its kind: observed, reported, estimated or assumed.
- **Why:** A number without them cannot be checked, refreshed or weighed against another.
- **Check:** Pick any number in the report: its ledger row shows source, as-of date and kind; your own arithmetic says it is yours.

### A4. Link each claim to the original source, and open it
Put the link at the claim itself, point it at the original (the dataset, paper or filing, not news about it), and read it.
- **Why:** Coverage reshapes what the original says, and a list of sources at the end cannot show which one backs which claim.
- **Check:** Every number, quote and dated event has its link beside it; every link opens; nothing rests on "reportedly".

### A5. Grade your evidence, and state flatly only what clears the bar
Grade every claim from A to E (table below), and state flatly only grades A and B.
- **Why:** A claim made by the party that profits from it is not settled, however plausible it sounds.
- **Check:** Each ledger row has a grade; no self-reported figure (a company's own, a vendor's own) appears as plain fact.

### A6. Keep fact, analysis and speculation apart
Label each statement as a sourced fact, your own analysis, or speculation, and never blend two of them in one sentence.
- **Why:** Most errors in reports are a guess presented as a finding, not a wrong number.
- **Check:** Analysis paragraphs carry no source links; each speculation says it is one and names who made it.

### A7. Search twice: once for the obvious, once for what it missed
Do a second search that excludes what the first found, looking for independent sources and for the people who disagree.
- **Why:** The first search returns the obvious and the flattering, and sources that copy each other repeat one error.
- **Check:** The ledger lists sources that do not derive from each other, and one section names who disagrees with whom.

### A8. Say what each source cannot see
For every source, write whom it covers, whom it misses and how it is skewed.
- **Why:** A sample speaks only for the people in it, and a silent gap reads as "nothing happened".
- **Check:** Each source has a "covers / misses" note; where nothing was found, the report says "not found" instead of guessing.

### A9. Argue the other side once, honestly
Write the strongest opposite reading (the steelman: the best opposing case) and test your hypothesis against it.
- **Why:** An analysis that only confirms the first idea, or the client's idea, has not been tested.
- **Check:** A contrary-view section names a mechanism (an assumption, a skipped step), not just doubt, and what would decide.

### A10. State your confidence and what would change your mind
Give each conclusion a confidence (high, medium or low) and the observation that would overturn it.
- **Why:** A reader cannot weigh a conclusion that sounds equally sure about everything.
- **Check:** Every finding and recommendation has a confidence and a falsifier (an observation that proves it wrong).

### A11. Date the work, and never rewrite the past
Put a cut-off date on every analysis, use nothing from after it, and make a new dated version instead of editing an old one.
- **Why:** The old version is the record of what was known then, and saying when a call failed makes the others credible.
- **Check:** Each report states its cut-off; no source postdates it; old versions stay in git; a failed call is named as failed.

### A12. Make every number reproducible
Keep raw inputs untouched, produce every derived number and file with code, and never hand-edit a generated file.
- **Why:** A number you cannot regenerate cannot be checked or refreshed, and a hand edit is lost at the next run.
- **Check:** A fresh run from the raw inputs reproduces every number in the report; every count lists the items behind it.

### A13. Keep personal data out
Use only data that is public or volunteered: nothing behind a login, no contact details, no real person as a persona or picture.
- **Why:** You answer for every person in your data, and a leaked name or key cannot be unpublished.
- **Check:** Before each push, search the repo and your prompts for names, emails and keys; every persona is fictional.

### A14. Write the guardrails before you publish
Before any slide, post or summary, list what may be stated flatly, what only with attribution, and what must never be stated.
- **Why:** Short versions travel without their context; the list stops a forbidden claim reaching a reader.
- **Check:** Every output obeys the list: no "never" claim appears, every "attribute" claim names its source, every link opens.

---

## Evidence grades (for A5)

| Grade | What it is | How you may state it |
|---|---|---|
| **A** | the original: the dataset, official statistic, paper, filing, your own measured data | flatly |
| **B** | two independent sources that do not copy each other | flatly |
| **C** | one credible secondary source | with attribution: "according to …" |
| **D** | self-reported by the party that benefits | with attribution, marked as theirs |
| **E** | circulating, no traceable source | never |

## The ledger, one row per number (A3, A5, A8)

`name · value or range · unit · kind · source (link) · as of · grade · covers / misses`
