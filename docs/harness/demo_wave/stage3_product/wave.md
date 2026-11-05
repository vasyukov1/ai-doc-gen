# Wave — stage3_product

64 findings from 4 personas × 2 runs on: stage3_day_one.jpg, stage3_section_test.jpg, stage3_what_you_missed.jpg, stage3_your_week_laptop.jpg.

Critique: `claude-sonnet-5-5` · triage: `claude-opus-5-5` · 2026-10-29 23:30

## Accepted into this wave (5)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W1 | On every check result, show 'Correct' or 'Not correct yet' with a tick or cross icon and colour, then 'Your answer: <exactly what the student entered>', and after a wrong answer add a 'Show answer' button beside 'Try again' that shows the expected answer, the sentence of the part it comes from and a link to that part. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night) | F05, F11, F21, F28, F34, F44 | Serves 'the student can test their own understanding' and keeps I2 ('a part with a check keeps its check'): the check stays as the teacher wrote it (no switch to multiple choice), only its feedback becomes readable and recoverable. |
| W2 | Replace the underlined 'Test this section' link with a pill button styled like 'Try again', placed at the end of the section it tests and labelled with that section and its size, for example 'Test §1 · How this course works · 4 questions'. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Petra Halden (the teacher) | F04, F12, F20, F30, F60 | The directive's 'the student can test their own understanding' needs a findable, unambiguous entry to the test; the section's parts, their order and their checks are unchanged (I1, I2). |
| W3 | Show one fixed sentence under 'Ask the library' on every screen, 'The library answers from the course; your teachers can read your questions, and they are not marked.', and remove 'Just you and the library'. | Ayla (design, year 1, phone), Thao (data science, laptop at night), Petra Halden (the teacher) | F14, F33, F43, F54, F58 | I4 says 'the teaching team can read what is said to the library', so 'Just you and the library' is false; the honest sentence goes everywhere, and no private library mode is added because I4 rules it out. |
| W4 | Attach an 'Asking about: <part title>' chip to the Ask the library composer only after the student taps a part's own ask control, keep it fixed while the student scrolls, and show 'Asking about: the whole course' with a clear × when nothing is chosen. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night), Petra Halden (the teacher) | F06, F29, F35, F42, F53 | The student controls 'what they ask' and 'the library answers from the course'; a topic the student never chose, changing on scroll, sends the question to the wrong part. |
| W5 | Show 'What you missed' as one numbered column in the order the parts were taught, one part per row with the bullet on the title's first line and its minutes at the right, a tick on parts the student has opened, a top line 'N of 21 opened · about M min left', and keep the floating ↓ button out of the list. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night), Petra Halden (the teacher) | F07, F13, F23, F31, F36, F41, F56, F62 | I1 ('separate parts ... in the order they were taught'): the two-column wrap breaks parts apart and scrambles their order; the fix restores it without hiding or re-sorting any part. |

## Deferred (11)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W6 | On phones, collapse 'Ask the library' to a one-line bar with a full-colour 'Ask' button that expands on tap, and give the ▾ on the button a visible word. | Ayla (design, year 1, phone), Rune (engineer, shortest path) | F01, F19 | Compatible: the student still controls 'what they ask' and more of the I1 stream is visible; outside this wave's cap, and W3 and W4 first fix what the panel says. |
| W7 | Show the day title once at the top of the day, drop the duplicate green day card and its '1 min ago', and merge the orange deadline card and 'Your deadlines' into one line in the stream's own day unit, for example 'Homework 1 · covers Days 1–4 · due Fri 23 Oct, 09:00 · 4 days left'. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night) | F02, F09, F18, F39, F48 | Form only (I5) and cheap, but lower severity than the five accepted; the unit wording must match the homework the teacher set (I4), which needs confirming first. |
| W8 | Default 'Quick knowledge sync' to the length the student picked last time (5 min the first time), rename it 'Quick review: pick how much time you have', add a visible 'Start' button and show the full 'Simpler version · 2 min' label. | Ayla (design, year 1, phone), Thao (data science, laptop at night) | F03, F40 | I5 allows it ('the form fits the student ... and the length'); two personas at moderate severity, outside this wave's cap. |
| W9 | Label the number on the floating ↓ button ('29 new parts below') and place the button in the margin so it never covers card text. | Ayla (design, year 1, phone) | F10 | Compatible with I1 (it points down the stream in taught order); raised once at severity 2, and W5 already keeps it off the missed list. |
| W10 | Rename the five 'Your week' counts to plain phrases that state each rule and who set it (for example '172 not opened yet', '3 you marked as got it'), show the meaning in one line on tap or hover, and keep the counts visible. | Ayla (design, year 1, phone), Thao (data science, laptop at night) | F08, F15, F37, F45 | I3 ('nothing grades them silently') backs saying who set each state, and the counts stay visible rather than dropped since the record is the student's to see; the rule behind 'understood' and 'fluent' must be defined first, so unclear for this wave. |
| W11 | Make the Homework 1 card one tap target that opens the homework with 'Needs: N parts, M not checked yet', each linking to that part in its place in the stream, and show the next deadline as the first line of 'Your week'. | Rune (engineer, shortest path) | F17, F24, F25, F27 | Compatible with I4 and I1 if the needed parts are the ones the teacher's homework names and they open in taught order; one persona and it needs homework-to-part links, so outside this wave's cap. |
| W12 | On each 'Three to look at next' card, name the part it builds on, add an 'Open' button, add one line saying why it was picked, and replace 'you have not shown' with 'not checked yet'. | Thao (data science, laptop at night) | F46 | Compatible (form, I5; plain words support I3); raised once at severity 2. |
| W13 | After Save on 'What did you learn?', show which period the entry was saved into with Edit and Delete, relabel 'today 12:00' as 'Posted today at 12:00', and add a line 'To keep something to yourself, write it in Notes'. | Thao (data science, laptop at night) | F47 | I3 ('correctable by them') backs edit and delete, and I4 already keeps notes the student's own, so private writing has a place without an 'Only me' switch on the reflection; raised once, outside the cap. |
| W14 | Put a one-word label under each icon of the phone's bottom bar (Home, Library, Alerts, Me, Ask, More). | Ayla (design, year 1, phone) | F16 | Compatible with the directive; raised once at severity 2. |
| W15 | Rename the overlapping 'Your map', 'Your skills', 'Your knowledge' and 'Your page' places to names that say what is inside, with a one-line grey description under each menu item. | Rune (engineer, shortest path), Thao (data science, laptop at night) | F32, F38 | Compatible with the directive; laptop only at severity 2, outside this wave's cap. |
| W16 | When a teacher opens the student app, show a fixed banner 'Teacher preview of the student view · DS403 · Day N · students see up to Day M', make check answers, 'Try again', 'Test this section', the reflection box and Save read-only with 'Preview only: nothing is saved', and never show 'You answered' text the viewer did not enter. | Petra Halden (the teacher) | F50, F51, F52, F57, F59, F64 | Supports 'the teacher stays in control of what is taught and when' (I4) and keeps a teacher's clicks out of students' records (I3), but it is a new mode raised by one persona and not cheap; first in line for the next wave given F57's severity 4. |

## Rejected (2)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W17 | Add a teacher switch in the student app that shows class counts and the names of students who are behind, on 'Your week' and on 'What you missed'. | Petra Halden (the teacher) | F49, F55, F61, F63 | Outside 'ReLibrary (the student's side of a course)' and against I3 ('what the student marks ... is theirs'): naming who is behind turns students' own marks into a class score; the 'whose numbers are these' part is solved by W16's preview banner, and class progress stays in the teacher's own tools. |
| W18 | Sort 'What you missed' so the parts Homework 1 needs come first, tag the others 'extra', hide them behind 'Show the other 16' or a 'Show only what Homework 1 needs' filter, and let 'I know this' remove parts from the list. | Rune (engineer, shortest path) | F22, F26 | Breaks I1 ('no wave replaces the stream with ... a re-sorted feed') and I5 ('never changes ... what the course requires'); the fast path can still come from a 'for Homework 1' label on rows kept in taught order (W11) and the student's own 'got it' mark (I3). |

## What must not change

- I1: Every screen shows the course as separate parts in one top-to-bottom stream in the order they were taught; no screen sorts, filters or merges parts into one long page.
- I2: Every part shows its verification badge with the same text for every student, and every part that had a check still has that check, unchanged in wording and form.
- I3: Every got it / unclear / not now mark the student made is shown to them and can be changed by them, and no new score, grade or ranking appears anywhere in the student app.
- I4: Students have no control that adds, reorders or removes the course's parts; their notes stay their own; and the teaching team can still read every question sent to the library.
- I5: Every personalised form (shorter, simpler, extras in another order) shows the same part text and the same required parts and checks as the teacher set.

## Walk again after the wave

1. stage3_section_test.jpg (phone): answer 'Why data' wrongly and confirm 'Not correct yet', 'Your answer: ...' as typed, and that 'Show answer' shows the expected answer with a link to its part (W1).
2. stage3_section_test.jpg (phone): find and start the test for §1 and for Section 2, each from a labelled button at the end of its own section (W2).
3. stage3_day_one.jpg, stage3_section_test.jpg and stage3_what_you_missed.jpg (phone): read the line under 'Ask the library' and confirm it is the same sentence on all three (W3).
4. Ask the library on stage3_section_test.jpg (phone): with nothing tapped the chip reads 'the whole course'; tap a part's ask control, scroll, and confirm the chip still names that part (W4).
5. stage3_what_you_missed.jpg (phone): read all 21 parts in one numbered column in taught order, with minutes and opened ticks, and nothing covered by the ↓ button (W5).
6. Invariant check (I1, I2, I3): scroll Day 1 from top to bottom and confirm every part is still separate, in taught order, with its verification badge and its check, and that pressing 'Show answer' leaves the student's own marks unchanged and adds no score.

## Spread between runs

| Persona | Findings per run | Agreement (screen + heuristic) |
|---|---|---|
| Ayla (design, year 1, phone) | [8, 8] | 0.45 |
| Rune (engineer, shortest path) | [8, 8] | 0.36 |
| Thao (data science, laptop at night) | [8, 8] | 0.36 |
| Petra Halden (the teacher) | [8, 8] | 0.45 |
