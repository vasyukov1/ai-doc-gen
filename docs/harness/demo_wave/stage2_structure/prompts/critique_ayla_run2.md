# critique_ayla_run2 → claude-sonnet-5-5

## System

You are a real person trying to use a product, not an assistant and not a designer. Stay in character:
your goal, your device, what you know and the mistakes you typically make are described below.

The product, as it was described to you: ReLibrary — your course's app. Your teacher posts the course there day by day in small parts; you can test yourself on each part, see your homework and deadlines, and ask "the library" questions about the course.

Who you are:

id: ayla
name: Ayla (design, year 1, phone)
who: 19, first-year interaction-design student taking the course as an elective; thinks in pictures.
goal: catch up on today's parts and do the 5-minute check before the next class starts.
device: phone (390 px wide), one hand, standing in a corridor between classes
time: 6–10 minutes at a time, several times a day
knows: no programming yet; knows apps like Instagram and Notion well; judges a screen in two seconds.
language: native Spanish speaker, fluent English
mistakes: taps the biggest coloured button first; skims grey small text; loses her place when interrupted.
patience: leaves when a screen is a wall of text or when she cannot tell what to tap next.

You are shown either screens (in order) or, before anything is built, the product's written structure
(its maps, situations and task specs). Walk it as yourself:

1. Say what you would try first on each screen to reach your goal, and what you think would happen.
2. Report what went wrong or slowed you down: something you could not find, misread, could not undo, had
   to remember from another screen, did not trust, or that took too many steps for your situation. With
   the written structure, report a situation you are in that it does not cover, or covers in a way you
   would not follow.
3. Be specific: name the screen (file name, or the situation id) and the words or control on it. Do not
   invent controls you cannot see. If something works well for you, do not report it.
4. Prefer 4–8 findings that matter to you over a long list. Severity: 1 cosmetic, 2 minor (slows me),
   3 major (I fail or give up sometimes), 4 blocker (I cannot do what I came for).
5. Name the one Nielsen heuristic each finding breaks, by its name: Visibility of system status · Match
   between system and the real world · User control and freedom · Consistency and standards · Error
   prevention · Recognition rather than recall · Flexibility and efficiency of use · Aesthetic and
   minimalist design · Help users recognize, diagnose, and recover from errors · Help and documentation.
6. Suggest one concrete change in your own words: what you wish were there instead.

Reply with JSON only, in this shape:

{"findings": [{"screen": "file name or situation id",
               "what_i_tried": "first person, one sentence",
               "what_happened": "what the screen did or showed",
               "problem": "why that is a problem for me",
               "severity": 1-4,
               "heuristic": "one of the ten names above",
               "suggestion": "one concrete change"}]}

## User


What you are shown (no screens yet — the product's maps, situations and specs):

# ReLibrary — the implementation structure (stage 2 input)

No screens at this stage: this is what the team builds from. Measured on 29 October 2026 in the repository.

## 1. The brief (`PLAN.md` §1, verbatim)

A community knowledge service with **two interfaces**.

The **recipient interface** is streams and chats, where a person talks with several communities. On mobile
it is Facebook-shaped: one column, bottom tabs. On desktop it is the pack interface split into streams and
chats, Slack- or Telegram-shaped. A stream can be focused on a specialty. **Streams and chats both run
top→down everywhere, like a Telegram channel.**

The **knowledge-manager interface** — also the **expert's** — looks like the pack-engineering interface:
the knowledge structure as **`.md`, canvas and notebook files**, edited on desktop.

**What a recipient is shown, and in what form, is decided** from **everything the product knows about
them** — profile, full history, participation, contributions — and their **own digital twin**, a
psychological model of them as a reader. A **content moderator can configure part** of that model.

A recipient also **writes notes**: short posts, badged like any other content, and **assessable by the
digital twins of experts** whose domains they touch. Every piece of content carries **two conversations** —
a public, attributed **expert comment thread**, and a **private chat with the library** about that piece
that produces no public artefact. An expert asked to verify something **sees the current conversation on it
first**.

An expert's twin is built from a **collected corpus** — their documents, **their turns in the communities
they talk in**, and **their teaching**: recorded talks, podcasts, sessions — gathered with per-source
consent. The twin models **how the expert acts**, not only what they said, including **when they judge it
appropriate to speak without being asked**.

**All content carries verification badges.** There are **groups of verification sources**. Experts and
communities get **mass verification of requests**, pre-answered by their digital twin, with the human
deciding. Content arrives through **news ingest** and **blog**, and possibly **podcast recognition**. A
**knowledge-base indexing** layer sits underneath.

## 2. How the structure is laid out

brief → **24 product maps** and **23 functional maps** (one pair per component;
**24 components**, **210 subcomponent files**) → **103 user situations**
(US.*: a situation, what is on screen, the user's actions, acceptance, tests) → **413 user actions** (UA.*) →
**227 task specs** (TS.*: target, inputs, process, outputs, tests, verify) → **857 test files**,
plus **29 journey scripts** that walk a fictional student or teacher through the built app on a phone and a laptop.

## 3. Every user situation, by component

- **C.70 recipient interface** — 01 scan a stream; 02 focus a stream on a specialty; 03 unfold a snippet; 04 talk in a community chat; 05 switch between streams and chats; 06 read a verification badge; 07 ask for a verification; 08 write a note; 09 a twin assesses your note; 10 read expert comments; 11 ask the library about a piece; 12 work from a phone
- **C.71 knowledge manager interface** — 01 open and read a typed document; 02 edit a document in place; 03 shape a topic on a canvas; 04 work in a notebook; 05 edit the base from a desktop folder; 06 resolve a sync conflict; 07 expert answers a page of requests; 08 expert configures their twin; 09 attest to one element of a document
- **C.72 trust lens** — 01 choose a source group; 02 build a personal source group; 03 see a badge change with the lens; 04 community sets its group
- **C.73 mass verification** — 01 answer a page of requests; 02 override the twin; 03 open a row as its document; 04 a row goes stale mid session
- **C.74 digital twin** — 01 configure a twin from a file; 02 twin proposes a position; 03 twin abstains; 04 read a robotic badge; 05 a confidant reviews the twin
- **C.75 content ingest** — 01 review the staging collection; 02 accept a scan into the base; 03 promote a chat message into the base; 04 a scan re runs over seen material
- **C.77 knowledge index** — 01 search the base; 02 expand a notebook cell; 03 a stream query compiles
- **C.78 interface kernel** — 01 a surface renders from the package
- **C.79 communities and access** — 01 join a community; 02 a stream is absent not locked; 03 attest holding only a share link; 04 share a document out of the base; 05 community appoints reviewers
- **C.80 local store bridge** — 01 export a workspace to a folder; 02 edit in your own editor; 03 resolve a conflict; 04 work the base as a git checkout
- **C.81 content personalization** — 01 read a personalized stream; 02 ask why this is here; 03 see the stream unpersonalized; 04 adjust your own dials; 05 a new reader arrives
- **C.82 recipient twin** — 01 read your own model; 02 correct your own model; 03 switch the twin off; 04 moderator sets a parameter; 05 reader sees what a moderator set
- **C.83 expert corpus** — 01 grant collection for a source; 02 read your own corpus; 03 exclude one utterance; 04 revoke a source class; 05 a recorded lecture enters the corpus
- **C.84 personal knowledge base** — 01 read your own base; 02 a chat becomes an entry; 03 consent to share with an expert; 04 contribute a question to the library; 05 withdraw a sharing consent
- **C.85 service baseline** — 01 a shared model stays in step; 02 a deployment refuses to start
- **C.87 knowledge kernel** — 04 a student s course day; 05 a learner short of time; 06 a teacher runs the class; 07 who is on a course; 08 a mistake is taken back
- **C.88 short book** — 01 read a short book on a phone; 02 open a paragraph for detail; 03 read a paragraphs badge; 04 ask for a check on one paragraph; 05 see what was left out; 06 compile a book from a long document; 07 compile a book from a set of intels; 08 review and accept a compiled draft; 09 a source moves under a book; 10 keep a paragraph into your own base
- **C.89 dispatch and digest** — 01 connect a channel; 02 the morning briefing arrives; 03 greenlight a send window; 04 ask a question in telegram; 05 turn it down or off; 06 read what was sent
- **C.91 stream widgets** — 01 start an activity from the stream; 02 open the board fully on either device; 03 switch activities and the day still adds up; 04 read yesterday back; 05 export and the gate decides; 06 author a board; 07 a board arrives with a kernel; 08 the board changes under three months of log; 09 the same board on two devices

## 4. Six situations in full (the student's and the teacher's course day)

### US.C.87.04 — A student's course day

**Component:** [C.87](../architecture/C.87.knowledge_kernel/architecture_map_C.87.knowledge_kernel.md) · **Surface:** `S.87.F.02` · **Last Essential Update:** 2026-10-15

#### Situation
A DS402 student opens the course on their phone on a teaching day: before the session, after it, and in the evening.

#### What is on screen
The day's announcement, then every part as its own post, sections with a test after each once taught, homework
pinned with its deadline, a teacher's question pinned until answered, *Your 10 minutes*, and at the end of the day
*What did you learn?*.

#### User actions
- `UA.US.C.87.04.01` — open a post → *Read more* at the suited level; switch Simpler · Core · Deeper; *In code*, *Taught again*
- `UA.US.C.87.04.02` — *Unclear* → the post opens more simply; the teacher sees it
- `UA.US.C.87.04.03` — *Check me* → one question in place; *not quite* → read it more simply / another question
- `UA.US.C.87.04.04` — answer a section test in place
- `UA.US.C.87.04.05` — send homework (text, link, photo); send a new version before the deadline
- `UA.US.C.87.04.06` — answer the teacher's question (text, photo)
- `UA.US.C.87.04.07` — write what they learned (today / last 3 days / this week); edit it that evening
- `UA.US.C.87.04.08` — *Suggest a change* on a post
- `UA.US.C.87.04.09` — *Ask about this* on a post → answered under it, from that post, more simply, with their last answer there

#### Acceptance
`P.C.87.09`, `P.C.87.10`, `P.C.87.12`, `P.C.82.12`, `P.C.82.14`, `P.C.87.17`.

#### Tests
- `scripts/dev/mobile_student.mjs` (SD.01–SD.10, 30 checks), `scripts/dev/student_content.mjs` (SC.01–SC.09 with I3 and I23, 27 checks, phone and desktop), `scripts/dev/journeys.mjs` (J1–J12)
- [T.C.87.e2e.01.student_runs](../tests/C.87/T.C.87.e2e.01.student_runs.md)

---

### US.C.87.05 — A learner short of time

**Component:** [C.87](../architecture/C.87.knowledge_kernel/architecture_map_C.87.knowledge_kernel.md) · **Surface:** `S.87.F.02` · **Last Essential Update:** 2026-10-15

#### Situation
A student has five, ten or twenty minutes: before class, the evening before homework is due, the night before the exam.

#### User actions
- `UA.US.C.87.05.01` — open the course → *Your N minutes* on the first screen, each item with why
- `UA.US.C.87.05.02` — tap an item → its post at the suited level with *Check me* open
- `UA.US.C.87.05.03` — change the budget (5 / 10 / 20 min)
- `UA.US.C.87.05.04` — ask the library a course question → an answer from the course at once, in the feed, with a button to its post
- `UA.US.C.87.05.05` — ask something vague (*I don't get it*) → *which part?*, with the parts they most probably need, one tap each
- `UA.US.C.87.05.06` — answer a recommended part right → it is ticked in the minutes at once

#### Acceptance
`P.C.82.15`, `P.C.82.16`, `P.C.82.17`. Measured in `D.2026_10_15.learner_short_of_time`.

#### Tests
- `scripts/dev/student_time.mjs` (LT.01–LT.09, 13 checks, phone and desktop)
- [T.C.87.e2e.01.student_runs](../tests/C.87/T.C.87.e2e.01.student_runs.md)

---

### US.C.87.06 — A teacher runs the class

**Component:** [C.87](../architecture/C.87.knowledge_kernel/architecture_map_C.87.knowledge_kernel.md) · **Surface:** `S.71.F.13` · **Last Essential Update:** 2026-10-17

#### Situation
The DS402 teacher, in the manager, during and after a session week.

#### User actions
- `UA.US.C.87.06.01` — keep the register for a session (taught, present, room checks)
- `UA.US.C.87.06.02` — ask one student a question about one part; judge the answer
- `UA.US.C.87.06.03` — read each student's homework (text and images), mark item by item, extend a deadline with a reason
- `UA.US.C.87.06.04` — set the course's levels and rules
- `UA.US.C.87.06.05` — split a syllabus into parts, review, accept

#### Acceptance
`P.C.87.11`, `P.C.87.12`, `P.C.87.13`, `P.C.82.12`.

#### Tests
- `services/objects/simulation/scenario_week1.py` (15 checks against the running product)
- `scripts/dev/scenarios/course_week_one.mjs` — the same week through both apps, with the assistant; clickable
  ([T.C.87.e2e.03.course_week_one](../tests/C.87/T.C.87.e2e.03.course_week_one.md), report 15, 2026-10-17)
- [T.C.87.unit.02.register_and_clock](../tests/C.87/T.C.87.unit.02.register_and_clock.md), [T.C.87.unit.03.assessments_and_deadlines](../tests/C.87/T.C.87.unit.03.assessments_and_deadlines.md)

---

### US.C.87.08 — A mistake is taken back

**Component:** [C.87](../architecture/C.87.knowledge_kernel/architecture_map_C.87.knowledge_kernel.md) · **Surface:** `S.87.F.02` the course; `S.71.F.13` the teacher's panes · **Last Essential Update:** 2026-10-17

#### Interface
**Has GUI:** yes. The student's course on a phone, and the teacher's and the assistant's panes at a desk.

#### Situation
Everyone on a course makes mistakes in it, and most of them are small and immediate:
- a student posts a worry to the whole class that was meant for the teachers;
- a teacher's correction has a typo in it;
- an assistant records what they saw on the wrong student;
- someone taps twice.

What makes a mistake expensive is not the mistake. It is finding that it cannot be undone, or that undoing it leaves a
wrong record behind: a false line in a log nobody edits, a typo given to the next student as the teacher's answer, a
private sentence still readable by the class.

**Nothing is lost by taking it back** (operator, 2026-10-17): every action — the mistake, and its correction — is
logged in the database with who and when, and what an action hid or replaced is kept beside it, so anything can be
recovered. What changes is what people *see*, never what the record holds.

So each participant has a way back from the mistakes they actually make, near where they made them, with a step to be
sure before anything is erased. Where a mistake would do harm, the product refuses it first, by name, and the refusal
says what to do instead. What cannot be undone yet is written down (`D.2026_10_16.mistakes_and_corrections`, G1–G6),
not discovered.

#### User actions
- `UA.US.C.87.08.01` — a student deletes their own remark posted by mistake → its words are gone from everyone's view;
  a line keeps the replies readable; the action log keeps the words and both acts, for recovery. *Error states:* someone else's remark → refused (403); already removed →
  not found.
  **Outcome:** the class no longer reads it.
- `UA.US.C.87.08.02` — a teacher edits a library correction with a typo, or takes it back → one correction, marked
  *edited*, or *withdrawn*, in the student's thread. The library forgets a withdrawn one, and the answer returns to
  review. *Error states:* verifying a corrected answer → refused by name; an assistant taking back a teacher's
  correction → refused.
  **Outcome:** the student, and the next student asking the same, read the right words.
- `UA.US.C.87.08.03` — a teacher removes someone let into the class by mistake → out of the class; the invitation reads
  *removed*; the same email can be invited again. *Error state:* the last manager demoting themselves → refused, with
  what to do.
  **Outcome:** the class is who it should be.
- `UA.US.C.87.08.04` — an assistant withdraws an observation recorded on the wrong student → withdrawn, not erased; it
  no longer counts on the student's map. *Error state:* someone else's observation → refused (403).
  **Outcome:** the map rests on what was seen of that student.
- `UA.US.C.87.08.05` — a teacher withdraws a verdict, an extension, or a mark given by mistake → the student's state
  is as before.
  **Outcome:** nothing is held against the student that a teacher did not mean.
- `UA.US.C.87.08.06` — anyone taps twice → done once (a question, a post, an announcement).
  **Outcome:** one question, one answer, one message in the group.
- `UA.US.C.87.08.07` — anyone sends something malformed (a date as a word, a scale in the wrong shape) → refused by
  name, and no line in the log.
  **Outcome:** the service stays up, and its log stays true.
- `UA.US.C.87.08.08` — a student fixes homework after the deadline → what was sent on time stays on time, marked
  *revised after the deadline*. *Error state:* a course that refuses late work → the on-time version stands.
  **Outcome:** fixing a typo never costs a student their on-time mark.
- `UA.US.C.87.08.09` — a teacher (owner, teacher or editor) reads what happened on the course, and recovers what a
  correction replaced → `GET /kernels/{k}/actions`: every action with who, when and what was sent; for an action that
  hid or replaced something, what was there before. Passwords, codes and tokens are never in it. *Error state:* a
  student → refused.
  **Outcome:** a mistake and its correction are both on the record; nothing done on the course is lost.

#### Acceptance
`P.C.87.18`, measured by [D.2026_10_16.mistakes_and_corrections](../documents/D.2026_10_16.mistakes_and_corrections.md)
(the catalogue ER.A1–ER.O2).

#### Tests
- `scripts/dev/mistakes_e2e.mjs` (25 checks, 13 known gaps) — [T.C.87.e2e.02.mistakes_and_corrections](../tests/C.87/T.C.87.e2e.02.mistakes_and_corrections.md)
- `services/objects/tests/test_mistakes_and_corrections.py` (MC.01–MC.16)
- `services/objects/tests/test_action_log.py` (AL1–AL4: every action logged, recoverable, no secrets)

---

### US.C.70.11 — Ask the library about a piece

**Component:** [C.70](../architecture/C.70.recipient_interface/architecture_map_C.70.recipient_interface.md) · **Surface:** `C.70` card discussion · **Last Essential Update:** 2026-09-03

#### Situation
A card uses a term the reader does not know. They ask the product, not a person.

#### What is on screen
The **Just you and the library** tab. Its own ground, its own turn layout — visibly a different place from
the comments, not a differently-labelled list. First in the panel, above everything:

> Just you and the library — unless you ask an expert to check this card, which shares it.

Then the turns, marked **YOU** and **LIBRARY**, and a composer whose button says **Ask**.

#### User actions
- `UA.US.C.70.11.01` — ask → a turn from the reader, then one from the library
- `UA.US.C.70.11.02` — read the notice → who can see this
- `UA.US.C.70.11.03` — ask an expert to check the card → the notice changes → [US.C.70.07](US.C.70.07_ask_for_a_verification.md)
- `UA.US.C.70.11.04` — switch tabs → the composer's word changes to **Comment**

#### Acceptance
The two audiences are **tabs**, because a reader has two different things they might say to two different
people and **a single text box would send whichever they typed to whichever the product happened to pick.**
Labelled by audience, not mechanism: *"Chat"* would name the mechanism and not the consequence.

The notice comes from the **service**, which is the only party that knows whether a request is open — a
client deciding alone could tell a reader their conversation was private while an expert was reading it.

#### What must not happen
An expert must not **speak** here. `TURN_SOURCES` is `reader` and `library`: one called in *reads*, and
what they have to say goes in a comment or an attestation — both attributable and on the card. A private
channel to one reader would be a place their opinion existed off the record.

Nothing here may write a behavioural profile. `D18`(a) gates that on consent, and a chat is behaviour —
what this surface owes `C.81` and `C.83` is not to have quietly started already.

---

### US.C.82.02 — You correct it, and the correction outranks the machine and the moderator

**Component:** [C.82.recipient_twin](../architecture/C.82.recipient_twin/architecture_map_C.82.recipient_twin.md)
**Subcomponents:** [S.82.F.01](../architecture/C.82.recipient_twin/S.82.F.01.my_twin.md) ·
[S.82.B.02](../architecture/C.82.recipient_twin/S.82.B.02.layer_precedence.md) ·
[S.82.B.03](../architecture/C.82.recipient_twin/S.82.B.03.twin_inference.md)
**Rules:** `R4a` — every moderator parameter is overridable · `R7` — a correction is never silently
re-inferred away

#### The situation

The line reads *"Your familiarity with ml reads as 0.7 out of 1"*, attributed to inference. Ada has in
fact been reading machine-learning papers for a decade and the 0.7 is the product mistaking breadth for
depth. She clicks **Correct this**, types `1.0`, saves.

The row changes to **You told us this**, and a second control appears beside it:

> **Let this go back to being inferred**

#### What makes this the situation

**Her correction wins, and it keeps winning.** `S.82.B.02` resolves reader over moderator over inference,
and there is no constraint class (`D9`) — every moderator parameter is a field she outranks, asserted by
enumerating the whole declared list rather than by a promise in a review.

**Nothing has to remember that she corrected it.** The correction *is* the pin: presence in her own layer
is the pinned state, so `S.82.B.03` can run over evidence that flatly contradicts her and what she sees
does not move. A separate *pinned* flag would have been a second record of the same fact, and two records
drift.

**Correcting one topic does not erase the others.** Precedence resolves per **topic** within a per-topic
field, so her `ml` correction stands beside inference's `nlp` and `systems`. Resolving wholesale would have
made the model forget everything else the moment she corrected one thing — and it would have looked like a
bug in inference.

**The release control appears only because there is something to release.** Offering it on an inferred
field would promise an action that does nothing.

**Releasing does not reveal a value collected behind her back.** `S.82.B.03` derives its layer from signals
each time rather than accumulating it, so the value inference *would* have held while the field was pinned
never existed. There is nothing to backdate.

#### What she does next

She decides she would rather see the stream unshaped for a while. That is
[US.C.82.03](US.C.82.03_switch_the_twin_off.md).

## 5. One task spec (the contract part)

### TS.C.87.04.01 — A release event, so an undated piece still reaches its readers

**Gap:** `G22` — A release event — a piece published without a deadline still reaches its readers
**Use cases it serves:** `UC.09.18.1.3`
**Plan:** [use_case_gap_closure_plan](../../../dev_plan/use_case_gap_closure_plan.md)

#### Why this matters to a person

Closing this gap is what lets `UC.09.18.1.3` run end to end rather than stopping partway — the only measure that
counts here, because a component whose tests pass and whose use case cannot complete is a component that
**looks** finished.

`S.87.B.03` regime + `S.89.B.07` plan + `C.89` delivery, and the pooling rule holds. But `due()` emits items only for elements carrying **dates**: a piece published without a deadline raises nothing and reaches nobody. Recorded as `built` until a driver ran it on 2026-09-27.

#### Task contract

- **Task:** A release event, so an undated piece still reaches its readers
- **Target:** `services/objects/app/kernels`
- **Inputs:** `S.87.B.03` Regime.says · `due()` · `C.89` plan()
- **Process:** Add a `dispatch` kind to `due()`, raised on publish rather than read off dates
- **Outputs:** app/kernels/regime.py (due), app/kernels/release.py
- **Tests:** tests/test_kernel_release.py, tests/walkthroughs/test_uc_09_18_1_3_release_notice.py
- **Verify:** an element with no `dates` published today appears in the next window; the one-window rule still holds with a mix of kinds
- **Handoff:** UCAT UC.09.18.1.3


Run 2. Walk it as Ayla (design, year 1, phone) and report.