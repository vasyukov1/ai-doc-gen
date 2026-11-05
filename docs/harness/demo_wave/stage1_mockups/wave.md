# Wave — stage1_mockups

65 findings from 4 personas × 2 runs on: stage1_stream.jpg, stage1_course.jpg, stage1_your_map.jpg.

Critique: `claude-sonnet-5-5` · triage: `claude-opus-5-5` · 2026-10-29 23:30

## Accepted into this wave (5)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W1 | Remove other students' activity cards (such as nour's 'worked / finished / askedFor' record) from every student's Stream, render the student's own activity cards as plain sentences ('You finished 3 of 5 on the sharding exercise sheet'), and print on each personal card, such as twin-ada's note, whom it is for and who can see it. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night), Petra Halden (the teacher) | F02, F11, F23, F28, F34, F42, F52, F53, F59 | I4: 'their questions and notes stay their own', so another student's 'askedFor' must not reach classmates. The card is removed, not reworded with a name and face as F02/F52 suggest. Raised by 4 personas; the filter is cheap. |
| W2 | On the Stream, course and Your map screens, use the same ask box: placeholder 'Ask the library about this course…'; one fixed line under it, 'Seen by you, the library and your teaching team, never by other students'; the button grey while the box is empty and solid green once there is text; the unlabelled ▾ replaced by its options written as words, none of which posts to the class; and bottom padding so the last card scrolls clear of it. Also label the top 'Everyone' dropdown with what it filters ('Show: Everyone') and remove it from Your map. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night), Petra Halden (the teacher) | F03, F12, F22, F36, F37, F43, F57, F60 | I4: questions stay the student's own and 'the teaching team can read what is said to the library'. So the line names the teaching team (not 'only you and the library', as F37/F41 ask), and no option posts into the course stream. Raised by 4 personas at severity 3; cheap. |
| W3 | On Your map, replace each part's 'Reply' link with two labelled buttons: 'Test myself', which opens that part's check, and 'Unclear? Say why', which opens the ask box headed 'About: <part title>' with the audience line from W2. Add one button at the top, '5-minute check · N questions', that runs the checks of today's parts. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night) | F08, F15, F32, F41, F49 | Directive: 'the student can test their own understanding of it and say in their own words what is clear and what is not'. I2 keeps every part's check one tap away. Severity 4, raised by 3 personas; it is a relabel plus two links. |
| W4 | Pin a one-line strip at the top of the Stream ('Today: N new parts · 5-minute check · Due: <homework>, <day time>') whose items open those parts, the check and the homework. Put day headings (Today, Yesterday, then the date) between the cards and mark parts the student has not opened 'New'. Keep the cards in taught order and their text unchanged. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night) | F01, F10, F18, F27, F47 | Directive: parts arrive 'top to bottom in the order they are taught, like a channel'. The strip and day headings sit over that order without changing it, so F01's newest-first re-sort (I1) and F10's rewording of the parts (I5) are not taken. Raised by 3 personas at severity 3. |
| W5 | Put a one-word label under each of the six bottom-bar icons, highlight the current icon on every screen including Your map, and make each number say what it counts ('8 new', '3 due'). | Ayla (design, year 1, phone), Rune (engineer, shortest path), Thao (data science, laptop at night) | F04, F26, F39, F45 | Navigation only: it touches no part, badge or mark, so every invariant holds, and it lets the student find what is due and keep 'their own pace'. Raised by 3 personas at severity 3; the cheapest change in the set. |

## Deferred (9)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W6 | Rename 'I understand everything here' to 'Mark all 3 as got it', show those parts afterwards as 'Got it (your mark)', and show an 'Undo' bar for 10 seconds after the tap. | Ayla (design, year 1, phone), Rune (engineer, shortest path) | F07, F16, F24, F33 | I3: the mark is the student's own word and must stay visibly correctable, so the chip stays. Removing it or gating it behind a test, as F07 offers, would put a score above the student's word. Sound, but 2 personas and over the cap of 5. |
| W7 | Replace the two different pill rows on the course and Your map screens with one tab row in the same place on both (Parts · My map · Homework · Skills · More). It fits 390 px without sideways scrolling and highlights the current tab. Also remove the empty sticky box that repeats the section title on Your map. | Ayla (design, year 1, phone), Rune (engineer, shortest path), Petra Halden (the teacher) | F09, F14, F20, F31, F65 | Navigation only, so every invariant holds, and 3 personas raised it. But it rebuilds the structure of both screens, more than this wave can re-walk; W5 handles finding your way for now. |
| W8 | Cut the course header to the title and one line ('12 parts in 4 chapters'), and move the description, keepers and revision behind an 'About this course' link. Replace the grey paragraph under the tabs with one line ('Each part is checked on its own; tap its badge to see by whom'), so the first part shows without scrolling. | Ayla (design, year 1, phone), Rune (engineer, shortest path) | F05, F13, F19, F30 | I2 holds: badges stay on every part and the per-part rule is still stated in one line. Sound and cheap, but only 2 personas and over the cap; first in the next wave. |
| W9 | Rewrite the badge text 'N experts · N twin' in plain words, give each badge icon one meaning, show the checkers' initials without overlap, and open a one-line legend with their full names when the badge is tapped. | Thao (data science, laptop at night), Petra Halden (the teacher) | F35, F44, F54, F62 | Serves I2 ('a badge reads the same to every student') and keeps every badge. Unclear until a plain-word name for 'twin' is agreed, so it waits for the next wave. |
| W10 | Collapse the pinned 'Ask the library' panel into a one-line bar that expands on tap, drop its duplicate heading, and give it the normal side margins. | Ayla (design, year 1, phone), Rune (engineer, shortest path) | F06, F21, F29 | Gives the parts the screen back without touching them (I1). Sound, but 2 personas at severity 2–3, and W2 already settles what the box does; re-walk that before shrinking the box. |
| W11 | On Your map, replace 'not yet' and 'worked on' with the student's own marks ('Got it', 'Unclear', 'Not started'). Replace the intro paragraph with one line ('Your own marks, not a grade: tap a part to read it or change your mark'), and show '0 of 12 got it' as a small bar. | Ayla (design, year 1, phone), Thao (data science, laptop at night) | F17, F40, F48 | I3: the states stay the student's own marks, not test results (F40's 'tested and understood' would turn them into a score). I2 keeps the word 'checked' for the badge. Sound, but over the cap. |
| W12 | Use one set of words on every screen: 'parts' for elements and pieces, 'chapters' for sections, 'Course content' for 'Content kernel', and 'Maintained by' for 'Kept by' and 'Keepers'. | Thao (data science, laptop at night) | F38, F46 | I5 holds: this renames the app's own labels, never what a part says. Only one persona raised it, so fold it into W7/W8 next wave. |
| W13 | Show each part's first full sentence on its Your map card and let the card expand in place to the whole part. | Rune (engineer, shortest path) | F25 | Keeps I1: each part stays its own card. Raised once, at severity 2. |
| W14 | Add a teacher 'See as a student' mode. A fixed banner names the class, the session day and the revision students see now; scheduled parts are marked 'Scheduled'; every marking and asking control is disabled, with the label 'Preview: nothing is saved'. | Petra Halden (the teacher) | F50, F51, F56, F58, F61, F64 | Serves 'the teacher stays in control of what is taught and when' (I4). But it belongs to the teacher's side, not the student app this wave builds, and only one persona raised it. |

## Rejected (1)

| # | Change | Who asked | Findings | Why |
|---|---|---|---|---|
| W15 | Show the teacher a class overview on Your map: for each part, how many students have learned it, and a list of the students who are behind, by name. | Petra Halden (the teacher) | F55, F63 | Breaks I3 ('nothing grades them silently'): students' own marks would become a named 'behind' list they never see. Still solvable: label this screen plainly as 'My map', and let the teaching team read, part by part, what students asked the library (I4). |

## What must not change

- I1: The Stream shows the course as separate part cards, top to bottom in the order taught. No screen merges parts into one long page, and no control re-sorts them (newest first or any other order).
- I2: Every part card on the Stream, course and Your map screens shows its verification badge. Every part that has a check still opens that check, and a part's badge reads the same for every student.
- I3: The student's marks (got it, unclear, not now) are shown to them as their own and can be changed by them. No hidden score, rank or test result replaces or overrides them.
- I4: A student has no control that adds, reorders or removes a course part (the ask box only asks the library). A student's questions and notes are never shown to other students, and the teaching team can still read what is said to the library.
- I5: Personalisation changes only the form, the order of extras and the length. A part's text and what the course requires read the same for every student.

## Walk again after the wave

1. Stream, first screen at 390 px: the Today strip names today's new parts, the 5-minute check and the next homework due, and each item opens its target (ayla, rune).
2. Stream, catching up: read yesterday's parts under the 'Yesterday' heading, in taught order, with unopened ones marked 'New' (thao).
3. Stream as a student: no card shows another student's activity or request, your own activity reads as a sentence, and each personal card (twin-ada's note) says whom it is for and who sees it (thao, petra).
4. Ask box on the Stream, course and Your map screens: same placeholder and the same audience line naming the teaching team; ▾ options written as words with none posting to the class; button grey when empty and solid once text is typed; last card fully visible above the box; 'Show: Everyone' label on the Stream and absent on Your map (thao, petra).
5. Your map: 'Test myself' opens that part's check, 'Unclear? Say why' opens the ask box headed with the part's title, and '5-minute check' starts today's questions (ayla, rune, thao).
6. Bottom bar on all three screens: every icon has a word, the current one is highlighted, every number says what it counts, and Homework is one tap away (rune).
7. Invariant check (I1, I2, I3): scroll the Stream top to bottom and confirm the part cards are in taught order with their badges unchanged. Then, on Your map, change one part's mark and change it back.

## Spread between runs

| Persona | Findings per run | Agreement (screen + heuristic) |
|---|---|---|
| Ayla (design, year 1, phone) | [9, 8] | 0.5 |
| Rune (engineer, shortest path) | [8, 8] | 0.56 |
| Thao (data science, laptop at night) | [8, 8] | 0.15 |
| Petra Halden (the teacher) | [8, 8] | 0.3 |
