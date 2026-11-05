# triage → claude-opus-5-5

## System

You triage usability findings from simulated users into one wave of product changes. You are the product's
owner, not one of the users. The product's directive is fixed; it does not change in this wave or because a
user asked:

<directive>
# Directive — ReLibrary (the student's side of a course)

Derived from `PLAN.md` §1 (the brief) and the student app as built (`apps/recipient`). Fixed for all three waves.

## Directive

ReLibrary brings a course's knowledge to each student in small parts, arriving top to bottom in the order
they are taught, like a channel, and in the form that fits that student. Every part can be checked: it
carries a badge saying who verified it, and the student can test their own understanding of it and say in
their own words what is clear and what is not. The teacher (the course's knowledge manager) stays in control
of what is taught and when; the student controls their own pace, their own record and what they ask, and
the library answers from the course.

## Invariants (what no wave may break)

- **I1 — Small parts, in order.** The course reaches the student as separate parts in a top-to-bottom stream
  in the order they were taught; no wave replaces the stream with one long page or a re-sorted feed.
- **I2 — Every part can be checked.** Every part shows its verification badge, and a part with a check keeps
  its check; a badge reads the same to every student.
- **I3 — The student's own word, not a hidden score.** What the student marks (got it, unclear, not now) is
  theirs, visible to them and correctable by them; nothing grades them silently.
- **I4 — The teacher decides what is taught.** Students cannot add to, reorder or remove the course's parts;
  their questions and notes stay their own, and the teaching team can read what is said to the library.
- **I5 — The form fits the student, the content does not change.** Personalisation changes the form, the
  order of extras and the length; it never changes what a part says or what the course requires.
</directive>

Do this:

1. **Deduplicate.** Group findings that describe the same underlying problem (across personas and runs).
   One change per group; list every finding id it covers and the persona ids that raised it.
2. **Decide each change:**
   - **ACCEPT** — it fixes a real problem *and* keeps every clause of the directive. At most 5
     ACCEPTs: prefer what several personas hit, high severity, and cheap to build; one wave must be small
     enough to build and re-walk before the next one.
   - **DEFER** — sound and compatible with the directive, but outside this wave's cap, unclear, or raised
     once at low severity.
   - **REJECT** — it conflicts with the directive (for example: it would hide what was checked, take
     control of what is taught away from the teacher, or replace small parts with one long page), or it
     asks for something the product deliberately does not do.
3. **Give a one-line reason for every decision that cites the directive** (quote the clause, or name the
   invariant). For a REJECT, say which clause it breaks and what part of the user's problem could still
   be solved without breaking it.
4. **List what must not change** in this wave: the directive's invariants, each as one checkable sentence.
5. **List the journey steps to walk again** after the wave is built: the screens or situations the
   accepted changes touch, plus one step that checks an invariant still holds.

Write each change as something a developer can build ("Show … on …", "Move …", "Rename … to …"), not as a
wish. Reply with JSON only:

{"changes": [{"id": "W1", "change": "...", "findings": ["F03", "F11"], "personas": ["..."],
              "decision": "ACCEPT | DEFER | REJECT", "reason": "one line citing the directive"}],
 "must_not_change": ["..."],
 "rewalk": ["..."]}

## User


Findings:

[
 {
  "id": "F01",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.04 (and US.C.87.05)",
  "what_i_tried": "I opened the course in the corridor, looking for the 5-minute check, and started scrolling to find it.",
  "what_happened": "US.C.87.04 lists the day's announcement, then every part as its own post, section tests, homework pinned, the teacher's question pinned, and only then 'Your 10 minutes' and 'What did you learn?'; US.C.87.05 says 'Your N minutes' is on the first screen.",
  "problem": "Two descriptions disagree about where the one thing I came for sits, and in the first one it is under a pile of posts and two pinned items, so on a phone I would see a wall and leave. It also says '10 minutes' when I have 5.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Put one big 'Your 5 minutes' card at the very top of the course screen, above the announcement and the pins, with the 5 / 10 / 20 choice right on it, and fold the rest of the day's posts underneath."
 },
 {
  "id": "F02",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.04 / US.C.87.05.04",
  "what_i_tried": "I came back to the course after a teacher walked up to me and tried to carry on where I stopped.",
  "what_happened": "Every part is a post in a top-to-bottom feed, and an answer from the library lands 'in the feed' with a button to its post; nothing in the structure says which parts are new or which I already read or did.",
  "problem": "I lose my place when I'm interrupted, and I would have to scroll and re-read titles to remember what I had finished; I also would not know where my question's answer appeared in a long feed.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Mark each part as new, read or checked, add a 'Continue where you left off' button at the top, and when I ask the library, take me straight to the answer instead of dropping it somewhere in the feed."
 },
 {
  "id": "F03",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.04.01",
  "what_i_tried": "I tapped Read more on a part and then looked for the level that suits me.",
  "what_happened": "The post opens 'at the suited level' and offers 'Simpler · Core · Deeper', plus 'In code' and 'Taught again'.",
  "problem": "I know no programming, so 'In code' means nothing to me and I would not tap it, 'Core' and 'Taught again' are unclear, and I was never told why it opened at that level; the switch is also hidden inside the opened post.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use plain words I can picture (for example 'Easier · Normal · More detail', 'See the code', 'What the teacher said in class'), say once in a small line why this level was picked, and show the switch on the post itself without opening it."
 },
 {
  "id": "F04",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.04.09 vs US.C.70.11",
  "what_i_tried": "I tapped 'Ask about this' on a post because I did not understand a term, and typed my question.",
  "what_happened": "The answer comes 'under it' on the post. A different place, the 'Just you and the library' tab, has a long grey notice about who can see it, and when I switch tabs the composer's button changes from Ask to Comment.",
  "problem": "I cannot tell whether my question under the post is private or whether the whole class reads it, and I skim small grey text, so I would miss the notice and could post a 'dumb' question to everyone by tapping the big coloured button.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Next to the send button show a coloured label that says 'Only you and the library' or 'Everyone in the class' in the same colour as the button, and ask me to confirm before anything goes to the class."
 },
 {
  "id": "F05",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.08 (with US.C.87.04.02 and .04)",
  "what_i_tried": "I tapped 'Unclear' and then a wrong answer in a section test by mistake, one-handed, and wanted to take it back.",
  "what_happened": "The only way back written for a student is deleting my own remark; 'Unclear' tells the teacher, a test is answered 'in place', and a double tap is simply done once.",
  "problem": "A thumb slip on 'Unclear' is seen by the teacher and I cannot see how to withdraw it, and a mistaken test answer cannot be changed; I would get embarrassed and stop trying.",
  "severity": 2,
  "heuristic": "User control and freedom",
  "suggestion": "Show a short 'Undo' line for a few seconds after 'Unclear' and after a test answer, and let me change my test answer once before it counts."
 },
 {
  "id": "F06",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.87.04.03 / US.C.87.05.06",
  "what_i_tried": "I tapped 'Check me' on one part after another to do my 5-minute check.",
  "what_happened": "I get one question in place on each post; 'not quite' offers 'read it more simply' or 'another question'; a recommended part is ticked in the minutes only when I answer right.",
  "problem": "The check is scattered across posts, a wrong answer leaves the item unticked so nothing shows I did it, and two choices after 'not quite' make me think about which to tap; I never see when my 5 minutes are done.",
  "severity": 2,
  "heuristic": "Visibility of system status",
  "suggestion": "Give me one big 'Start' button that runs the questions one after another with a counter like '3 of 5', tick each item as done even if I got it wrong (marked 'to revisit'), and make 'another question' the big default after 'not quite'."
 },
 {
  "id": "F07",
  "persona": "ayla",
  "run": 1,
  "screen": "US.C.82.02",
  "what_i_tried": "I looked at the line about my own model and tried to understand what it says about me.",
  "what_happened": "It reads 'Your familiarity with ml reads as 0.7 out of 1', attributed to inference, with 'Correct this', then 'You told us this' and 'Let this go back to being inferred'.",
  "problem": "'ml', '0.7 out of 1' and 'inferred' are not words I use, I would take it as a grade on me, and I would not know what the release button does.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Say it in a picture and a plain sentence ('We think you know a medium amount about machine learning'), with a slider I can drag, and label the release button 'Let the app guess again'."
 },
 {
  "id": "F08",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.04 and US.C.87.05",
  "what_i_tried": "I opened the course on my phone in the corridor and looked for the 5-minute check to do before class.",
  "what_happened": "US.C.87.05 says Your N minutes is on the first screen, but US.C.87.04 lists it after the announcement, every part's post, the section tests, the homework and the teacher's question, with What did you learn? at the very end. One calls it Your 10 minutes and the other Your N minutes, and neither says where the 5 / 10 / 20 switch is.",
  "problem": "I have 6 to 10 minutes and I judge a screen in two seconds. If the check is not the first thing I see, I scroll through posts looking for it, lose my place and give up. I would never find a switch for 5 minutes if it is not in plain sight.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Put one big card at the very top of the course screen, above the announcement, called the same thing every day (Your 5 minutes), with the 5 / 10 / 20 choice inside it and the items listed right under it."
 },
 {
  "id": "F09",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.04 (What is on screen, User actions 01 to 09)",
  "what_i_tried": "I scrolled today's course screen to catch up on the parts I missed.",
  "what_happened": "Every part is its own post, each followed by a section test once taught, plus pinned homework, a pinned teacher's question and the end-of-day box. Each post carries up to eight or nine actions: Read more, Simpler / Core / Deeper, In code, Taught again, Unclear, Check me, Suggest a change, Ask about this.",
  "problem": "That is a wall on a 390 px screen. I cannot tell which button is the one I need, so I tap the biggest coloured one, which might be Unclear or Suggest a change instead of Check me. With one hand and no time, I leave.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Show today's parts as a short list of titles with a tick, grouped by section, and open one post at a time. Make Check me the single coloured button on a post and tuck the rest into a small More menu."
 },
 {
  "id": "F10",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.04.01 and US.C.87.04.02",
  "what_i_tried": "I tapped Read more, then looked for a way to make the post easier to follow.",
  "what_happened": "The spec offers Read more at the suited level, a Simpler · Core · Deeper switch, In code, Taught again, and an Unclear button.",
  "problem": "I do not know what level I am at or who decided it. Core means nothing to me. In code is for programmers, and I cannot tell if Taught again is the teacher's recap or just a repeat of the same text. I would guess wrong and read the wrong version.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use plain words I would say aloud, such as Easier / Normal / More detail, and label the other two Show the code and Watch the explanation again. Show which level I am on right now."
 },
 {
  "id": "F11",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.04.09 and US.C.87.05.04",
  "what_i_tried": "I tapped Ask about this under a post and typed a basic question I was embarrassed to ask in front of the class.",
  "what_happened": "The answer comes under the post, or in the feed, with a button to the post. Nothing in these situations says whether my question and the answer are visible to classmates or the teacher.",
  "problem": "Under a post and in the feed look like public places, like comments on Instagram. I would not trust that it is private, so I would not ask. I skim small grey text, so a footnote would not reach me anyway.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Next to the Ask box, in readable type, show a lock and Only you see this. If a question can be seen by others, say so before I send it and make me choose."
 },
 {
  "id": "F12",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.08 (UA.US.C.87.08.01 and 08.06)",
  "what_i_tried": "I typed a worry about the homework into the first box I found and posted it, then realised it was for the teachers and not for the class.",
  "what_happened": "The structure only gives me a way back after the fact: I delete my remark and a line stays so the replies still read. Nothing stops the mistake or says who will read it before I post.",
  "problem": "I am the kind of person who posts while walking and gets interrupted. The class can already see it, and the leftover line tells everyone something was removed. I would be embarrassed and would stop writing in the course.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Show the audience on the box itself with a two-way choice, To the class / Only to the teachers, and after posting give me a short Undo bar (about 10 seconds) that removes the post without leaving a line."
 },
 {
  "id": "F13",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.87.04.03 and US.C.87.05.02",
  "what_i_tried": "I started a Check me question, got interrupted by a friend, and came back a few minutes later.",
  "what_happened": "The spec shows Check me opening one question in place, with not quite leading to read it more simply or another question. It says nothing about where I am in the check, which item is next, or how to resume.",
  "problem": "I lose my place when interrupted. I would not know whether I had finished the 5-minute check, which parts are still open, or whether not quite counted against me. I would have to start again, so I would quit.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Keep a visible progress strip on the Your 5 minutes card (for example 2 of 4 done, with the unfinished item highlighted), and when I reopen the app, take me straight back to the question I left."
 },
 {
  "id": "F14",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.70.11 (Just you and the library) and UA.US.C.70.11.03 to .04",
  "what_i_tried": "I asked the library what a term meant, then tapped to have an expert check the card and switched tabs to the comments.",
  "what_happened": "A long grey sentence sits at the top of the panel: unless you ask an expert to check this card, which shares it. The notice only changes after I ask the expert, and the composer's button quietly changes from Ask to Comment when I switch tabs. Nothing says what happens to text I have already typed.",
  "problem": "I skim small grey text, so I would not read that asking an expert makes my private chat shared. I could also type a private question, switch tabs and tap the same button, posting it publicly under my name.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Before an expert check is sent, show a plain confirm step with a big Cancel: The expert and others will see this chat. Keep a different colour and a lock icon on the private tab, and clear or warn about a typed draft when I switch tabs."
 },
 {
  "id": "F15",
  "persona": "ayla",
  "run": 2,
  "screen": "US.C.70.11, US.C.87.04.09 and US.C.87.05.04",
  "what_i_tried": "I wanted to ask the course something about a part I did not understand.",
  "what_happened": "There are three different places to do that: Ask about this under a post, a course question answered in the feed, and the Just you and the library tab on a card. They use different words (post, part, card, piece) and buttons (Ask about this, Ask, Comment).",
  "problem": "I do not know which one is the real one or whether they give the same kind of answer. I know Instagram and Notion, where there is one comment box in a fixed place, so here I would stall.",
  "severity": 2,
  "heuristic": "Consistency and standards",
  "suggestion": "Call it one thing in one place on every post, an Ask button with the same label everywhere, and use post (or part) consistently instead of card and piece."
 },
 {
  "id": "F16",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.87.04",
  "what_i_tried": "I open the course and look for the one thing I came for: what is due and when.",
  "what_happened": "The day opens with the announcement, then every part as its own post, a test after each section, and the homework 'pinned with its deadline' somewhere in that stream. The whole thing is one column that runs top to bottom, written for a phone.",
  "problem": "Homework is one item among a long run of posts. I have to scroll and skim a feed to find a deadline, and I do not know how far down it sits. I will not read announcements or section posts to get one date. I would give up and ask a classmate.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Put one 'Due' strip at the very top of the course: the homework title, the deadline and a Send button. The announcement and the part posts go below it, and I can ignore them."
 },
 {
  "id": "F17",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.87.04 / US.C.87.05",
  "what_i_tried": "I see the homework and ask myself which parts I must know to do it, and I look for a link from the homework to those parts.",
  "what_happened": "The structure has a pinned homework, a teacher's question, and a 'Your N minutes' list where each item comes 'with why'. Nowhere does it say that a homework item lists the parts it needs, or that 'Your N minutes' is built from the homework.",
  "problem": "My goal is 'what is due, and what I must know for it, nothing else'. The second half has no home. I would have to read the homework text, guess which of the many parts it means, and look them up one by one. That is exactly the digging I will not do.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Under each homework item show 'You need:' with the 2-4 parts as one-tap links, each opening with Check me ready. Anything outside that list stays collapsed."
 },
 {
  "id": "F18",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.87.05",
  "what_i_tried": "I have three minutes on the train, so I open the course and set the budget to the smallest value.",
  "what_happened": "The first screen shows 'Your N minutes', each item with its reason, and the budget switches between 5, 10 and 20 minutes. Nothing says the list starts with what is due, or that I can ask for only that.",
  "problem": "There is no choice for my real situation. The smallest budget is already more than I have. Each item carries a 'why' that I will skip, and I do not trust a list I did not see being filled by the deadline. A recommended list that ignores a due homework is worse than no list.",
  "severity": 2,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add an 'Only what is due' option next to 5 / 10 / 20, and always put due items first in the list. Show the 'why' as a short tag, not a sentence."
 },
 {
  "id": "F19",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.87.04 (situation: 'on a teaching day')",
  "what_i_tried": "On a Wednesday with no session, I open the course on my laptop to check what is still due from last week.",
  "what_happened": "US.C.87.04 is written for a teaching day: before the session, after it, in the evening. US.C.87.05 covers a learner short of time. No situation covers an ordinary day between sessions, or one place that lists all homework and deadlines across sessions.",
  "problem": "The brief says I can 'see your homework and deadlines', but the structure only puts them inside a day's feed. If last week's homework is on last week's day, I have to remember which day it was on and scroll back to it. A deadline I forget is the one I miss.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Add a situation and a screen: 'What is due', one list across the whole course, sorted by deadline, with sent / not sent next to each. Make it the first tab."
 },
 {
  "id": "F20",
  "persona": "rune",
  "run": 1,
  "screen": "UA.US.C.87.04.01",
  "what_i_tried": "I open a post and want the short, technical version with code, because I know Python and skip basics.",
  "what_happened": "I tap *Read more*, and it opens 'at the suited level'. Then I use Simpler · Core · Deeper and *In code* to change it. The structure does not say what 'suited' means for me or where I can see it.",
  "problem": "That is two taps for each post, and I may be starting from a level that someone else (the system) chose for me. With dozens of posts I will do this dozens of times. I do not see why it picked what it picked, so I cannot tell if I am missing something.",
  "severity": 2,
  "heuristic": "Visibility of system status",
  "suggestion": "Show the current level as a small label on every post. Let me pick 'Core + In code' once for the whole course and have it stay."
 },
 {
  "id": "F21",
  "persona": "rune",
  "run": 1,
  "screen": "UA.US.C.87.04.02 / US.C.87.08",
  "what_i_tried": "I tap *Unclear* on a post by mistake (I expect it to work like a reaction in Slack), and I give a wrong answer to *Check me* with a thumb slip.",
  "what_happened": "*Unclear* makes the post open more simply, and 'the teacher sees it'. US.C.87.08 lists ways back for remarks, corrections, invitations, observations, verdicts and double taps. It has nothing for taking back an *Unclear*, a test answer, or a 'What did you learn?' entry.",
  "problem": "In Slack or GitHub I can click a reaction off again. Here an accidental tap tells my teacher I did not understand something I did, and I cannot see how to remove it. A wrong tap on a test may also count against me.",
  "severity": 2,
  "heuristic": "User control and freedom",
  "suggestion": "Make *Unclear* a toggle that I can tap again to withdraw, with a short 'Undone' toast. Let me re-answer a *Check me* question, and add these cases to the list in US.C.87.08."
 },
 {
  "id": "F22",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.82.02",
  "what_i_tried": "I want to say once that I know Python well, so that the basics are skipped, and I look for where to do that.",
  "what_happened": "The path is to read a row such as 'Your familiarity with ml reads as 0.7 out of 1', click *Correct this*, type a number, and save. This is done topic by topic. The undo is a control called 'Let this go back to being inferred'.",
  "problem": "I do not think in 0.7 out of 1, and I will not walk through a list of topics to fix them one by one. The undo wording is long and unclear to me. I would not use this screen, so the system keeps its wrong guess about me and keeps giving me basics.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Offer plain choices per area ('I know this well / a bit / not yet') in one list, and a 'Skip basics' switch for the whole course. Name the undo 'Reset'."
 },
 {
  "id": "F23",
  "persona": "rune",
  "run": 1,
  "screen": "US.C.70.11",
  "what_i_tried": "I have a question about a term in a card, so I open the card's discussion and type it in the box, like in any chat.",
  "what_happened": "The panel starts with a notice sentence ('Just you and the library — unless you ask an expert to check this card, which shares it.'). Under it are the YOU and LIBRARY turns, and a composer whose button says **Ask**. When I switch tabs, the same composer's button becomes **Comment**, and a comment is public and attributed.",
  "problem": "I skip the notice, and the box looks the same in both tabs; only a single word on the button changes. In Slack, one box means one audience. I could post a half-formed question to the expert thread under my name, or think a comment is private. Also, the course's first screen (US.C.87.05.04) lets me ask the library a different way, with the answer in the feed, so there are two ways to ask with different results.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Make the two composers look different, with a placeholder that says 'Private: only the library sees this' or 'Public: everyone sees your name', and a different colour. Ask for a confirm on the first public comment, instead of relying on the notice."
 },
 {
  "id": "F24",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.87.04 and US.C.87.05",
  "what_i_tried": "I opened the course on a teaching day looking for one place that lists what is due and by when.",
  "what_happened": "The first screen is 'Your N minutes', which is built from a time budget. Homework is 'pinned with its deadline' inside the day's feed, between 'every part as its own post' and the 'What did you learn?' prompt. No screen lists the due items across days, and nothing says whether I have already sent mine.",
  "problem": "My whole goal is 'what is due'. I would have to scroll a day's feed, and then the earlier days, to assemble that list myself. A time budget is not a deadline. I would give up and ask a classmate.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Put a 'Due' strip at the very top of the course: one row per homework with the deadline and a state (not sent, sent v2 on time, revised after the deadline). One tap on a row opens it."
 },
 {
  "id": "F25",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.87.04 (homework pinned, parts as separate posts)",
  "what_i_tried": "I opened the pinned homework and wanted to see which parts I need to know to do it, so I could skip the rest.",
  "what_happened": "The structure shows the homework pinned with its deadline and every part as its own post with a test after each section. It does not say the homework points at the parts it needs. The only prerequisites mentioned are wired inside the kernel, not shown to me.",
  "problem": "I would have to remember which parts the homework is about and scan 40 or so posts to find them. That is exactly the reading I am trying to avoid, and I would miss one.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Under the homework, list 'Needs these parts' as links with a tick on each part I have already passed its Check me for. Tapping one opens that part at the right level."
 },
 {
  "id": "F26",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.87.04 (UA.01, UA.02, UA.03, UA.08, UA.09)",
  "what_i_tried": "I looked at a single part to get the code and the one fact I need, and move on.",
  "what_happened": "Each post carries Read more, Simpler · Core · Deeper, In code, Taught again, Unclear, Check me, Suggest a change and Ask about this. 'In code' sits behind Read more and then the level switch.",
  "problem": "That is eight or more controls per post and I have to open the post before I reach the one I want, which is the code. Two taps for what should be one. I will ignore most of them, and they may hide the one that matters.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Show the code toggle on the collapsed post. Keep Check me and Ask about this visible, and put Unclear, Suggest a change and Taught again under a single overflow menu."
 },
 {
  "id": "F27",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.87.05 (UA.01 to UA.03)",
  "what_i_tried": "I set 'Your N minutes' to my three minutes and took the first item.",
  "what_happened": "The budget has three values, 5, 10 and 20 minutes. Each item comes 'with why', and tapping it opens the post at the suited level with Check me open.",
  "problem": "I have three minutes, not five, so the list is already too long for me. The 'why' line is a sentence I will not read. There is also no way to say 'I know this, skip it', so Check me is the only shortcut and it is not labelled as one.",
  "severity": 2,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a 3 minute budget or a free number field. Add an 'I know this' button on each item that runs the one check and, if I pass, ticks the part off and moves to the next."
 },
 {
  "id": "F28",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.70.11 (Just you and the library tab) with UA.US.C.87.04.09 and UA.US.C.87.05.04",
  "what_i_tried": "I typed a question into the box under a card the way I would in Slack, and hit the button.",
  "what_happened": "There are two tabs. The privacy notice, 'Just you and the library — unless you ask an expert to check this card, which shares it', is a paragraph at the top of the panel. The only other difference is that the send button says Ask on one tab and Comment on the other. There are also two more entries for asking the library, 'Ask about this' on a post and 'ask the library a course question'.",
  "problem": "I skip paragraphs, so I will not read the notice, and the one-word button change is easy to miss. I could post a naive question publicly under my name, or write to the library thinking an expert reads it. Three differently named doors to the same library also make me unsure which one keeps my history.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Make the button itself carry the audience: 'Ask the library (private)' and 'Post comment (public, as Rune)'. Use one name, 'Ask the library', for all three entries."
 },
 {
  "id": "F29",
  "persona": "rune",
  "run": 2,
  "screen": "Section 1 (brief): 'what a recipient is shown is decided from everything the product knows about them' and US.C.81.03",
  "what_i_tried": "I assumed the due items and the parts I must know would be in the same place for everyone, and I would not have to check.",
  "what_happened": "The brief says the digital twin decides what and in what form a recipient is shown. The only way to see the stream unshaped is a separate situation, 'see the stream unpersonalized'.",
  "problem": "I cannot trust a list that a model may have reordered or trimmed for me. If a deadline is not in my view I fail without noticing. I also will not go looking for a personalization switch.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Keep deadlines and homework outside personalization, and show a fixed 'N due, all shown' counter on the course screen. Any personalized list gets a visible 'show all' toggle right next to it."
 },
 {
  "id": "F30",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.87.04, US.C.87.05 and US.C.70.05 (switch between streams and chats)",
  "what_i_tried": "I opened the course on my work laptop, between meetings.",
  "what_happened": "Both course situations are written as 'on their phone', with the course as a first screen with 'Your N minutes'. The brief's desktop is a Slack-like split of streams and chats, and the course day is not placed in it, so I cannot tell where the course lives or how a due item shows there.",
  "problem": "Most of my use is at a desk, and I expect the Slack or Telegram pattern: a left list with the course marked as due or unread. If the course is one stream among several communities I may open the wrong one and not see the minutes block at all.",
  "severity": 2,
  "heuristic": "Consistency and standards",
  "suggestion": "Specify the laptop course day too: the course pinned at the top of the streams list with a due count badge, opening to the same 'Due' strip and 'Your N minutes' block as the phone."
 },
 {
  "id": "F31",
  "persona": "rune",
  "run": 2,
  "screen": "US.C.89 (02 the morning briefing arrives, 05 turn it down or off) and TS.C.87.04.01 (release event)",
  "what_i_tried": "I connected Slack or Telegram so that I get only what is due and nothing else.",
  "what_happened": "A morning briefing arrives, and an undated piece now also raises a release notice, pooled in the next send window. The only control listed is 'turn it down or off', and the structure does not say what the briefing holds or whether I can filter it.",
  "problem": "Every undated reading now reaches me as a notice, so the useful message (what is due) is buried in noise. A piece with no date is told to me with no sign of whether it is required. If the only choice is all or off, I turn it off and miss the real deadlines.",
  "severity": 2,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a per-channel setting 'only things with a deadline'. Put the due items first in the briefing, each with a link to the homework, and mark undated pieces 'optional' or 'required'."
 },
 {
  "id": "F32",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 (UA.US.C.87.04.02, the button Unclear)",
  "what_i_tried": "I read a part about a method I did not follow and wanted to press Unclear, only to get an easier explanation, without anyone knowing that I was lost.",
  "what_happened": "The structure says: Unclear -> the post opens more simply; the teacher sees it. The same button does two things, it changes my reading and it reports me to the teacher.",
  "problem": "I take labels literally, so I read Unclear as a setting for me, but it is also a message to the teacher, and I only learn that from the small words after the arrow. I worry that a teacher who sees many Unclear taps from me thinks I am failing, so I would not press it at all and I would sit stuck instead. And once it is pressed there is no way to take it back.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Make two separate controls: a private one, 'Explain this more simply' (nobody sees it), and a second, optional one, 'Tell the teacher this part is unclear', which says before I press it what the teacher will see (for example 'Your teacher sees how many students pressed this, not who') and lets me undo it."
 },
 {
  "id": "F33",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 (Check me, UA.US.C.87.04.03) and US.C.70.11 (ask an expert to check this card)",
  "what_i_tried": "I saw a part with Check me and also, on another screen, the line 'unless you ask an expert to check this card', and I tried to work out which one is a quiz for me and which one sends my words to a person.",
  "what_happened": "In US.C.87.04 'Check me' gives me one question in place, with 'not quite' if I am wrong. In US.C.70.11 and C.88 'check' means an expert verifies the card or paragraph, and that shares it. The word check means two different things.",
  "problem": "My English is slow and I take the word as it is. 'Check me' could mean 'check my work' and send it to someone. I would hesitate to press it, or I would press the expert check thinking it is a self-test and share my private conversation. When a label has two meanings I stop trusting the app.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use the word 'check' for only one thing. Call the self-test 'Quiz me on this part' and the other 'Ask an expert to verify this', and show a short line under each one: 'only you see your answer' / 'an expert will read this'."
 },
 {
  "id": "F34",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 (a student's course day) and US.C.87.05",
  "what_i_tried": "At 22:30 I opened the course to read properly the teaching day I missed two days ago.",
  "what_happened": "The situation describes the student on the teaching day itself: before the session, after it, in the evening. The first screen shows the day's announcement, a pinned teacher's question, pinned homework, 'Your 10 minutes' and, at the end of the day, 'What did you learn?'. Nothing says how I reach an earlier day, which days I have not read, or whether the pinned question and 'What did you learn?' still wait for me for a day that is over.",
  "problem": "My goal is exactly this case, reading a missed day in order and fully. I would scroll a feed of 'every part as its own post' and not know where the missed day starts, which parts I did not read, and if the deadline pinned there already passed. I would worry that I am late and not see it explained.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Add a situation 'I missed a day' with a day list that shows for each day how many parts I have read and what is still open (homework and its deadline, the teacher's question), and a button 'Read this day from the start' that goes through the parts in order."
 },
 {
  "id": "F35",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.05 (UA.US.C.87.05.01 and .03, Your N minutes, budget 5 / 10 / 20 min)",
  "what_i_tried": "I looked for a way to tell the app I have 60 to 90 minutes tonight and want the whole day, not a short selection.",
  "what_happened": "The only budgets are 5, 10 or 20 minutes, and the first screen is built around 'Your N minutes'. The product decides what fits in the time, but does not say how it counts a minute.",
  "problem": "I read 'Your 10 minutes' literally as the time I should spend, and 20 is the largest, so my long session has no place. Also I read English slowly, so I think a count made for a fast reader would be too short for me, and I would not know if a 10-minute item takes me 25.",
  "severity": 2,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a budget 'All of it' (and 60 / 90 min), and show for each item how long it takes at my own reading speed, with a way to say 'I read slowly' once."
 },
 {
  "id": "F36",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 (UA.09 Ask about this), US.C.87.05 (UA.04 ask the library a course question) and US.C.70.11 (Just you and the library)",
  "what_i_tried": "I wanted to ask about one thing I did not understand in a part and I looked for the place where only I and the library are present.",
  "what_happened": "There are three names and two places: 'Ask about this' on a post, answered 'under it'; 'ask the library a course question', answered 'in the feed'; and a tab 'Just you and the library' whose notice reads 'unless you ask an expert to check this card, which shares it'. The words post, part, card and piece are used for what looks like the same thing.",
  "problem": "I cannot tell if these are one function or three, and I cannot tell if my question and the answer 'in the feed' are visible to the class. The notice is hard for me: 'which shares it' - shares what, with whom? Asking without feeling exposed is my main goal, and here I cannot trust that it is private.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one name and one place for asking, put the same plain sentence next to every box ('Only you can see this question and the answer'), and rewrite the notice as 'Only you and the library can see this. If you ask an expert to verify it, the expert will also see this conversation.' Use one word (part) for the piece everywhere."
 },
 {
  "id": "F37",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 and US.C.87.05 (the answer boxes: teacher's question pinned, section test, Check me, vague question 'which part?')",
  "what_i_tried": "I opened several parts at once on my laptop, answered a section test, started a Check me question, and asked the library 'I don't get it'.",
  "what_happened": "The structure puts several kinds of answer in place, in the feed, each under a different post. A vague question returns 'which part?' with the parts I most probably need, one tap each, and the library's answer to a course question arrives 'in the feed, with a button to its post'.",
  "problem": "I open many things and lose which one I was answering. After I tap one of the guessed parts I do not see my first question any more, and an answer that arrives in the feed can land far from the part it talks about. I can answer the teacher's pinned question when I think I am answering a test.",
  "severity": 2,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Always show my question quoted above its answer, with the name of the part it is about, and mark the one box that is waiting for me ('Waiting: your answer to the teacher's question'). Keep a small list 'My open questions' so I can return to them."
 },
 {
  "id": "F38",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.08 (UA.US.C.87.08.01 a student deletes their own remark, and UA.US.C.87.08.09 the teacher reads the action log)",
  "what_i_tried": "I wrote a worry in the class by mistake that I did not want the teacher to see, and I tried to delete it.",
  "what_happened": "The words are 'gone from everyone's view' and a line keeps the replies readable, but 'the action log keeps the words and both acts', and the teacher (owner, teacher or editor) can read every action with what was sent and recover what a correction replaced. There is no student action to take back an Unclear, a question or a Check me.",
  "problem": "'Delete' does not delete: the teacher can still read it, and nothing on the screen told me. I would feel safe and then be exposed, and when I learn it I would not trust any button in the app. The line that stays behind also shows to the class that I removed something.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Call it 'Hide from the class', and in the confirmation step say in one plain sentence who can still see it ('Your teacher can still see the original in the course record'); also give me 'Undo' for a few minutes after I press Unclear or send a question."
 },
 {
  "id": "F39",
  "persona": "thao",
  "run": 1,
  "screen": "US.C.87.04 (UA.US.C.87.04.01 Read more at the suited level; Simpler · Core · Deeper; In code; Taught again)",
  "what_i_tried": "I pressed Read more and wanted to read the statistics part in depth but the Python part in a simple way.",
  "what_happened": "The part opens 'at the suited level', chosen by the product, and I can switch Simpler · Core · Deeper. The structure does not say what level I am on, who decided that it suits me, or whether my choice is remembered. 'Taught again' and 'Core' have no explanation.",
  "problem": "I know statistics well and Python only medium, so one level for me is wrong for half of the parts, and I check claims, so I do not trust a level I did not choose. 'Core' and 'Taught again' are short labels that I take literally: core of what? taught again by whom, when?",
  "severity": 2,
  "heuristic": "Visibility of system status",
  "suggestion": "Show the current level on the post ('You are reading: Core'), say why it was picked, remember my choice per topic, and use plain labels: 'Easier', 'Standard', 'More detail', and 'What the teacher said again in class' instead of 'Taught again'."
 },
 {
  "id": "F40",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04 (UA.US.C.87.04.02, the Unclear button)",
  "what_i_tried": "I reached a paragraph about weighting that I did not follow and wanted to press Unclear, but first I looked for who would see it.",
  "what_happened": "The structure says that after Unclear the post opens more simply and 'the teacher sees it'. Nothing on screen is described that tells me, before I press, whether the teacher sees my name, a count, or only that the part was hard.",
  "problem": "I am afraid the teacher reads Unclear as me failing. The only place the text says the teacher is told is in the spec, not in front of me. So I would avoid the button and keep struggling quietly, which is the opposite of what I came to do.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Put one plain line next to Unclear that says exactly what the teacher gets, for example 'Your teacher sees how many students pressed this, not who', and offer a second way to get the simpler version that tells nobody."
 },
 {
  "id": "F41",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04.09 / US.C.87.05.04 / US.C.70.11 (Ask about this, ask the library a course question, Ask)",
  "what_i_tried": "I wanted to ask the library about a term in one post, and I saw three places that all say 'ask'.",
  "what_happened": "'Ask about this' on a post is answered under that post. Asking the library from Your N minutes is answered in the feed with a button to the post. In the card's Just you and the library tab the button is 'Ask', and when I switch to the other tab the same box says 'Comment'.",
  "problem": "Three things share one word and they end up in different places with different audiences. I take labels literally, so I cannot tell whether my question lands under the post for the whole class to read, in a private chat, or in the feed. After the second surprise I stop trusting the word.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one word for the private question to the library and say 'only you' on the button itself (for example 'Ask the library (private)'), and keep a different word for anything the class can see."
 },
 {
  "id": "F42",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04 (UA.US.C.87.04.01 and .03: Check me, Taught again, Core, In code)",
  "what_i_tried": "I looked for a way to test myself on the part I had just read and found the button 'Check me'.",
  "what_happened": "The labels are short phrases: 'Check me', 'Taught again', the level names 'Simpler · Core · Deeper', 'Your 10 minutes', 'Suggest a change'. None of them is explained on the screen as described.",
  "problem": "My English is good but slow, and idioms confuse me. 'Check me' could mean 'grade me and tell the teacher'. 'Taught again' does not say by whom, or whether it is a recording or a new text. 'Core' does not tell me whether it is the minimum or the main part. I would press and find out, and a graded-looking button is exactly what I hesitate to press.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use full plain wording, such as 'Test myself (one question, not graded)', 'Read the lecture again, shorter' and 'Main level', and put one short explanation under each when I hover or tap."
 },
 {
  "id": "F43",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04.01 and US.C.87.04.02 (Simpler · Core · Deeper, 'the post opens more simply')",
  "what_i_tried": "I pressed Unclear because the English was dense, expecting the same content in easier words.",
  "what_happened": "The structure says the post 'opens more simply' and offers Simpler. It does not say whether that means simpler words or less content. No situation covers language, or a reader whose first language is not English.",
  "problem": "My statistics is strong and my English reading is slow. If Simpler cuts the ideas I already know, I lose the details I want to check. If it only simplifies, I may still stumble on the idioms. I cannot ask for what I actually need, which is clear wording at the same depth.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a separate control, 'Same content, plainer English', next to Simpler, and let me keep the hard terms with a short definition so I can look them up."
 },
 {
  "id": "F44",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.70.11 (UA.US.C.70.11.03, ask an expert to check the card; the notice 'Just you and the library — unless you ask an expert to check this card, which shares it')",
  "what_i_tried": "After asking the library something I thought was a basic question, I wanted to ask an expert to check the card as well.",
  "what_happened": "The notice changes to say it is shared. The brief says an expert asked to verify 'sees the current conversation on it first'. The notice says the card is shared and does not say that my own earlier questions are shared.",
  "problem": "The word 'it' is unclear, because it could mean the card or my conversation. I would press the button to verify the card and expose every question I asked, including the naive ones, and I only find out afterwards. That is the exposure I am trying to avoid.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Before it sends, show a confirmation that lists exactly which of my messages the expert will read, with the choice 'Send the card only' or 'Send the card and my questions', and make the notice say 'my questions' instead of 'it'."
 },
 {
  "id": "F45",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.08 (A mistake is taken back)",
  "what_i_tried": "I pressed Unclear, or sent 'Ask about this', by accident, and looked for how to take it back.",
  "what_happened": "The situation gives a student only two ways back: delete their own remark, and a double tap counting once. Nothing says how to take back Unclear, Ask about this, a Check me answer or a What did you learn? entry. The log keeps everything.",
  "problem": "My typical mistake is pressing something and worrying what it shows the teacher. If the only undo is for remarks, I have to hope. I will be slow and careful before every tap, and a long evening of reading will not feel safe.",
  "severity": 2,
  "heuristic": "User control and freedom",
  "suggestion": "After Unclear and Ask about this, show 'Undo' for a few minutes, and next to my past actions list 'Take back' for each one, saying what the teacher will see once I do."
 },
 {
  "id": "F46",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04 and US.C.87.05 (the first screen: Your 10 minutes, budgets 5 / 10 / 20 min, no situation for a missed day)",
  "what_i_tried": "I opened the course at 22:00 to read the day I missed from the start, properly, in order.",
  "what_happened": "US.C.87.04 covers a teaching day as it happens. US.C.87.05 starts the screen with 'Your N minutes' and offers 5, 10 or 20 minutes. I found nothing about catching up on a past day, and nothing for 60 to 90 minutes.",
  "problem": "My whole goal is to read the day I missed in full, and the first screen pushes me toward a 20-minute summary. I would have to guess how to scroll back to yesterday's posts, and I could not tell what I had missed and what I had read.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a 'Catch up on a missed day' entry on the first screen that opens that day's parts in order, marks which I have read, and add a longer budget such as 'Whole evening (60+ min)'."
 },
 {
  "id": "F47",
  "persona": "thao",
  "run": 2,
  "screen": "US.C.87.04 (What is on screen: the pinned teacher's question, Check me, section tests, What did you learn?)",
  "what_i_tried": "I opened several posts at once to compare them and started typing an answer.",
  "what_happened": "The day shows a pinned teacher's question, a Check me question in place, a section test after each section, the homework and What did you learn? at the end. Each one takes an answer in place.",
  "problem": "Several different questions are open at once and they look alike. I lose track of which one I am answering, and I may type into the teacher's pinned question when I meant the test, and the other way round. I worry about sending the wrong thing to the teacher.",
  "severity": 2,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Keep the question, with its post title, visible above every answer box, mark only the teacher's question 'Goes to your teacher', and keep a drafted answer when I move to another post."
 },
 {
  "id": "F48",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.06 (A teacher runs the class) and US.C.87.04 (A student's course day)",
  "what_i_tried": "I looked for the situation where, the evening before, I open tomorrow's session and see exactly what my students will see on their phones.",
  "what_happened": "US.C.87.06 lists five things for me (keep the register, ask one student a question, read homework, set levels and rules, split a syllabus) and none of them is a look at tomorrow as a student; the student's day (US.C.87.04) is written only from the student's side.",
  "problem": "My whole evening is this check, and I would have to guess: I cannot tell what a student will see tomorrow or whether the manager and the student course show the same thing, so I would not trust anything I looked at.",
  "severity": 4,
  "heuristic": "Visibility of system status",
  "suggestion": "Add a teacher situation 'Preview tomorrow as my students will see it': the announcement, the parts, the tests, the pinned homework and my pinned question for a chosen day, in the same screen the student has, with a bar that says 'Preview, students see this on <date> at <time>'."
 },
 {
  "id": "F49",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.06 (UA.US.C.87.06.03) and the title 'C.87 07 who is on a course'",
  "what_i_tried": "I looked for the place that tells me who in my class is behind.",
  "what_happened": "I can read each student's homework one by one and mark item by item; the only other thing is a situation titled 'who is on a course', whose content is not shown and which sounds like a list of members, not of progress.",
  "problem": "Opening every student's homework to find who has not sent anything takes far more than my 30 minutes with a full class, and 'who is on a course' does not tell me it will answer 'who is behind'.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Give me one class overview for a chosen day: every student, with 'homework sent / not sent / late', tests done and the *Unclear* taps, sorted so those who are behind are at the top, and name that situation 'Who is behind' in my words."
 },
 {
  "id": "F50",
  "persona": "petra",
  "run": 1,
  "screen": "Section 1 (the brief: what a recipient is shown is decided from everything the product knows about them) and US.C.87.04",
  "what_i_tried": "I tried to take the screen I see on my phone as the student view and assume every student sees the same posts.",
  "what_happened": "The structure says each student gets a personalised stream, with Simpler · Core · Deeper levels, their own *Your 10 minutes*, and a moderator who can configure part of the model; nothing says what a teacher sees, or how to see one student's view or the class default.",
  "problem": "That is exactly my mistake: I will take my preview for what everybody sees. If two students see different things I cannot tell which one my screen is, and I lose trust.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "On any student-view screen I open, show a fixed label such as 'You see: the class default, Core level' with a switch to 'as <student name>' or 'as Simpler / Deeper', and say plainly what is the same for everybody and what is personal."
 },
 {
  "id": "F51",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.06 and US.C.87.08.05 (teacher's panes in the manager)",
  "what_i_tried": "I tried to act on a student's homework and extend a deadline, without being sure which class and which day was open.",
  "what_happened": "The text names the actions (mark item by item, extend a deadline with a reason, withdraw an extension or a mark) but never says that the pane shows which course, which class and which session day I am in, nor whom an extension applies to.",
  "problem": "I forget which class or day I am looking at, so I could extend a deadline or give a mark in the wrong place and only find out when a student writes to me.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Keep the course, the class and the day in a header that never scrolls away, and make every confirm line say it back: 'Extend the deadline of Session 9 homework for Anna K. in DS402, class B, to Friday?'."
 },
 {
  "id": "F52",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.04 (UA.04.02 *Unclear*, UA.04.03 *Check me*, UA.04.09 *Ask about this*) and US.C.70.11",
  "what_i_tried": "I tried the student view on my phone and tapped *Check me* and *Unclear* on a post to see what a student gets.",
  "what_happened": "These buttons write things: a test answer, an *Unclear* that 'the teacher sees', a question to the library with a chat that is behaviour; nothing says what happens when the person tapping is the teacher.",
  "problem": "I could add false *Unclear* taps or answers to my own class's data, or have my own questions on my teacher record, by only testing the screen, and I would not know.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "When I am looking as a teacher, make the student controls try-outs: they work, they say 'Preview: nothing here is saved or counted', and they leave no trace in the class's *Unclear*, tests or history."
 },
 {
  "id": "F53",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.08 (UA.US.C.87.08.02, 08.05 and 08.09)",
  "what_i_tried": "I tried to see how I take back something I changed by accident, and where I find what happened.",
  "what_happened": "A change to a library correction becomes marked *edited* or *withdrawn* in the student's thread; recovering what a correction replaced is described as `GET /kernels/{k}/actions`; nothing covers undoing an edit to a part, the day's announcement or something published for tomorrow.",
  "problem": "My slip is shown to students as 'edited' or 'withdrawn', and the way back is a web address meant for developers, so I would not use it. For the content I prepare most, no undo is described.",
  "severity": 2,
  "heuristic": "Help users recognize, diagnose, and recover from errors",
  "suggestion": "Add a button 'What changed on this course' in my pane, a list by date and person with 'Restore this version' next to every edit of a part or announcement, and a draft state where my changes are not seen by students until I press Publish."
 },
 {
  "id": "F54",
  "persona": "petra",
  "run": 1,
  "screen": "TS.C.87.04.01 (a release event, so an undated piece still reaches its readers), C.89 'greenlight a send window'",
  "what_i_tried": "I tried to find out if a part I wrote for tomorrow will actually reach my students, and when.",
  "what_happened": "The spec says a piece published without a deadline 'raises nothing and reaches nobody', until a release event is added, and speaks of 'the next window' and 'the one-window rule'; no situation shows me whether my part is a draft, published or already sent.",
  "problem": "A part can look finished to me and not reach any student, and I would find that out in the morning from the class; I cannot see published versus sent.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Show on each part and on the day a plain status line: 'Draft', 'Published, appears in the students' course', 'Sent to Telegram / Slack on <date>, <time>', and warn me in the evening if something for tomorrow has no release."
 },
 {
  "id": "F55",
  "persona": "petra",
  "run": 1,
  "screen": "US.C.87.06 (UA.06.01, UA.06.05) and Section 1 (the knowledge-manager interface: .md, canvas and notebook files)",
  "what_i_tried": "I tried to follow the words used for my side: 'the manager', 'keep the register (taught, present, room checks)', 'split a syllabus into parts', files in .md, canvas and notebook.",
  "what_happened": "The words come from the system (kernel, register, intels, typed document) and 'room checks' does not tell me what I record; I know the content, not these terms.",
  "problem": "I have to learn a vocabulary from the engineers' interface before I can prepare a session, and I will not know if 'register' is my attendance list or something else.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use teacher words on my screens: 'Attendance', 'Today's session', 'Students' homework', 'Course parts', and explain 'room checks' in one line where it appears."
 },
 {
  "id": "F56",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.06 — A teacher runs the class (and section 3, C.87 list)",
  "what_i_tried": "I looked for the one place where I can open tomorrow's course the way my students will see it, so I can check it tonight.",
  "what_happened": "US.C.87.06 gives me five things to do (register, ask one student, read homework, set levels and rules, split a syllabus); no situation in the whole list lets a teacher view the student course, and the only student day described, US.C.87.04, is written for a student on a phone, on a teaching day.",
  "problem": "My main goal tonight is to see what students will see tomorrow. If there is no preview, I will assume the manager pages I see are the student view, or I will go to the student side and not know whether it is the same. I lose trust the moment I cannot tell whether students see what I see.",
  "severity": 4,
  "heuristic": "Match between system and the real world",
  "suggestion": "Give me a 'Students tomorrow' view: the course day exactly as a student gets it (announcement, parts, tests, homework with deadline, pinned question), reading only, with a visible label that says 'preview, not sent' and the class and day it shows."
 },
 {
  "id": "F57",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.07 — who is on a course / US.C.87.06 / US.C.87.04.02",
  "what_i_tried": "I looked for 'who is behind' and opened 'who is on a course' first, because it is the only title that sounds like a list of my students.",
  "what_happened": "By its title (and the invitation and removal wording in US.C.87.08.03) it is a roster of who is in the class, not who is behind; in US.C.87.06 the only way to see a student's state is to read one student's homework at a time, and students' taps on 'Unclear' or their 'What did you learn?' entries 'are seen by the teacher' (US.C.87.04.02) with no place named where I read them.",
  "problem": "To answer 'who is behind' I would have to open every student separately and remember what I saw; with 30 minutes in the evening I will not finish, and the signals students send me have no home that I can find.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "One class page: a row per student with homework sent / missing against the deadline, parts they ticked, and the parts where they tapped 'Unclear', sorted so the ones who are behind are on top, each row opening that student's detail."
 },
 {
  "id": "F58",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.06.01 — keep the register for a session (taught, present, room checks) and US.C.87.04 'What is on screen'",
  "what_i_tried": "I tried to mark a session as 'taught' tonight, to see whether the tests that come 'after each section once taught' would then show up in the student view.",
  "what_happened": "The student's screen is described as showing 'sections with a test after each once taught', so the register entry I make as a record also decides what students see; the structure does not say that marking 'taught' releases anything, and US.C.87.08.05 lists taking back a verdict, an extension or a mark, but not a register entry.",
  "problem": "I could change what the whole class sees tomorrow while I think I am only checking or keeping notes, and I may not find a way to take it back, which is exactly the accident I am afraid of.",
  "severity": 3,
  "heuristic": "User control and freedom",
  "suggestion": "Say on the register control what it does ('Marking this taught shows the tests to the class now') and ask once before it does; add a 'Take back' next to the entry that returns the class to the state before."
 },
 {
  "id": "F59",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.04 — A student's course day (opened on my phone)",
  "what_i_tried": "I opened the course on my phone to check how it looks, then tapped 'Check me' and 'Unclear' on a part to see what the student gets.",
  "what_happened": "Every action on that screen is a real action for a real person: 'Unclear' means 'the teacher sees it', 'Ask about this' is answered 'with their last answer there', homework and 'What did you learn?' are sent, and US.C.87.08.01 shows that anything posted lands in front of the whole class; nothing says that I am viewing as a teacher or that my taps are not counted.",
  "problem": "I cannot test the student side without perhaps leaving answers, flags or remarks in the class record, and I do not know whether my own taps will appear in the numbers I use to see who is behind.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "When a teacher opens the student surface, show a fixed bar 'Preview: nothing you tap here is saved or shown to the class', and make the send buttons do nothing but show a sample result."
 },
 {
  "id": "F60",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.06 — A teacher runs the class (surface S.71.F.13) and US.C.87.04 (before the session, after it, in the evening)",
  "what_i_tried": "I opened the register and homework for the class and wanted to be sure which class, which session and which day I was working on.",
  "what_happened": "US.C.87.06 says only 'during and after a session week'; no situation names a current class, session or day on screen, and the student's day (US.C.87.04) changes between before the session, after it and in the evening, so it matters which of those I am looking at.",
  "problem": "I forget which class or day I am on, and at 9 pm tonight the course would look different from how it looks tomorrow morning. I could read, mark or extend a deadline in the wrong session and not notice.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Keep the class name, session number and date, and 'now' or 'tomorrow morning' fixed in the header of every teacher screen, and ask me to confirm the session before I extend a deadline or mark."
 },
 {
  "id": "F61",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.06.04 — set the course's levels and rules, with US.C.82.02 (you correct it) and US.C.87.04.01 (Simpler · Core · Deeper)",
  "what_i_tried": "I set the course's levels and rules and then wanted to confirm that each student will get what I set.",
  "what_happened": "US.C.82.02 says the student's own correction outranks the machine and the moderator, and the student surface switches Simpler, Core, Deeper per person, so two students may see different text for the same part; nothing in what I am shown lets me see how many students differ from what I set.",
  "problem": "'What students see tomorrow' is not one thing; I cannot tell which version of a part I am checking, and I cannot tell who has overridden my settings, so I cannot trust that I have seen the same as them.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "On each part in the manager, show the level the class gets by default, and a line 'N students see another level or changed their own setting' that opens the list; let me switch my preview between Simpler, Core and Deeper."
 },
 {
  "id": "F62",
  "persona": "petra",
  "run": 2,
  "screen": "TS.C.87.04.01 — A release event, and US.C.89.03 greenlight a send window",
  "what_i_tried": "I prepared a part for tomorrow tonight and tried to find out whether publishing it sends it to students now, in tomorrow's morning message, or only when I say.",
  "what_happened": "The task spec says a piece with no deadline 'is raised on publish' and 'appears in the next window', and the only control I can see is 'greenlight a send window' in C.89; neither says what the next window is or what time it falls on.",
  "problem": "Publishing is outward-facing and I cannot take back a message that has already gone to students' Telegram; if I do not know when it goes I will either not publish or publish too early.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Before and after I publish, show 'Goes to students at 08:30 tomorrow (Telegram and in the course)' with the time, and let me hold it or change the time until it is sent."
 },
 {
  "id": "F63",
  "persona": "petra",
  "run": 2,
  "screen": "US.C.87.08 — A mistake is taken back (UA.US.C.87.08.09 and the 'Tests' list)",
  "what_i_tried": "I tried to find where I read what happened on the course and recover something a correction replaced.",
  "what_happened": "The way is written as `GET /kernels/{k}/actions`, errors are written as '(403)' and 'not found', and the text says what cannot be undone yet is 'written down' in G1–G6 and that 13 gaps are known, in a document, not in the app.",
  "problem": "I am a teacher, not a developer: I cannot use a web address and I do not know which mistakes I can still take back. 'Nothing is lost' does not help me if I cannot reach it and I do not know the 13 exceptions when I am about to tap something.",
  "severity": 3,
  "heuristic": "Help users recognize, diagnose, and recover from errors",
  "suggestion": "A 'History' page in the teacher's panes, in plain words ('You extended Anna's deadline on Tuesday, 18:02 - Take back'), and a short list there of what can still not be taken back, shown next to the action that is affected."
 }
]