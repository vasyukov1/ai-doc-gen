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
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I opened the app between classes and looked at the first screen for today's parts and the 5-minute check.",
  "what_happened": "The first card is a long block of text titled \"Chunk map retention raised to 14 days for the audit work\", dated 2026-08-11. The next ones are dated 08-14, 08-18 and 08-24. Nothing says \"today\", \"new\" or \"check\", and the only status is a small grey dashed \"Not yet checked\" chip.",
  "problem": "I cannot tell what is new or what I should do now. It reads like a list of old technical news in words I do not know (chunk map, config server, rollback), and I do not know whether any of it is today's part. I would close the app after two seconds.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Put a big \"Today\" card at the very top: \"2 new parts · 5-minute check\" with one clear button. Show dates like \"Today\" or \"Yesterday\" instead of 2026-08-11 07:40, and put the newest first."
 },
 {
  "id": "F02",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I skimmed the cards to see which one was mine.",
  "what_happened": "The second card is tiny grey text with raw labels: \"nour\", \"worked: the sharding exercise sheet  finished: 3 of 5\", \"askedFor: an easier sheet\". Below it is another card from \"twin-ada\" with a purple chip \"0 experts · 1 twin\" and round avatars.",
  "problem": "The first card looks like a debug log. It seems to be another student's data, which makes me worry about privacy and stops me trusting the app. \"Experts\", \"twin\" and the initials in circles mean nothing to me, and I skip grey small text anyway.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Write these cards as sentences with a face and a name, for example \"Nour asked for an easier sheet\", or hide them. Replace \"0 experts · 1 twin\" with plain words or an icon that opens a short explanation when I tap it."
 },
 {
  "id": "F03",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I wanted to tap the most obvious green button at the bottom to get help.",
  "what_happened": "\"Ask the library\" is pale green with white text, so it looks switched off. Its small dark ▾ part on the right looks like the real, active button. The last card in the list is cut off by the box.",
  "problem": "I am not sure which part to tap. I would hit the dark ▾ because it is the strongest colour, and I do not know what it opens. I also lose the card behind the box.",
  "severity": 2,
  "heuristic": "Visibility of system status",
  "suggestion": "Make it clearly disabled (grey) until I type, then fully green. Hide the ▾ menu behind a \"more\" option, and add some space at the bottom of the list so the last card is not cut off."
 },
 {
  "id": "F04",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_stream.jpg (the bottom bar is the same on all three screens)",
  "what_i_tried": "I looked for \"homework\" or \"check\" in the bottom bar with my thumb.",
  "what_happened": "There are six icons with no words: books, lines, person, a question box, a speech bubble, three dots. They carry numbers 11, 8, 8 and 4, and I do not know what the numbers count.",
  "problem": "I have to guess what each icon is and tap around. The numbers look like things I owe (11 unread? 8 tasks?) but they do not say what. When I get interrupted I forget which icon was which.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Put a short word under each icon (Course, Stream, Me, Check, Chat, More). Keep a number only on the one that means \"to do\", and make it say what it counts."
 },
 {
  "id": "F05",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I tapped the course icon to see what I should read for today.",
  "what_happened": "A big heading \"Sharding and data placement\" with a paragraph, then \"Kept by the infra community\", \"Revision r14 · Sep 9, 2026\", \"12 elements in 4 chapters\", \"KEEPERS\" with two circles, and tabs \"Content kernel 12 / Experts / Rooms / From this ke…\". Below that is a grey paragraph: \"Every element is checked on its own. There is no score for the kernel as a whole — three confirmed definitions say nothing about a contested caveat.\"",
  "problem": "That is a wall of grey text before I see any content, and I need a dictionary for it (kernel, keepers, revision, elements, contested caveat). I do not know which of the 12 pieces are today's, and the tab row is cut off at the right so I do not know what else is there.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Say it like a course: \"Today's parts (3)\", then the list. Move \"Kept by\" and \"Revision\" into a small \"About\" link. Delete the grey paragraph about the kernel, or show it only on tap. Use plain tab names such as Parts, People, Chats."
 },
 {
  "id": "F06",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I tried to read the first definition card while holding the phone with one hand.",
  "what_happened": "The \"Ask the library\" area (a heading, a grey sentence, an empty text box and a wide pale button) is stuck to the bottom, with the bottom bar under it. The card above is cut off after one sentence. The box has no side margin and its text touches the screen edge.",
  "problem": "Almost half the screen is taken up by the ask box and bar, so I can see only one card at a time. The heading and the button both say \"Ask the library\". The edge-to-edge look makes it seem broken, so I trust it less.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Collapse the ask box into a small round \"Ask\" button, like a chat bubble, that opens the box when I tap it. Add normal side margins and use the freed space for the parts."
 },
 {
  "id": "F07",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I saw the outlined green chip \"I understand everything here\" under the section title and thought this was my way to finish and move on.",
  "what_happened": "It is the first thing under \"What a shard is\" and looks like a button, but nothing says what it does. There is no sign that I get tested, and I cannot see how to undo it.",
  "problem": "I would tap it to get rid of the red \"not yet\" labels and I might mark all 3 as learned without testing myself. I would then trust a green map that is wrong, and I could not tell how to take it back.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Remove the chip, or ask \"Mark all 3 as learned? You can undo this.\" and show an \"Undo\" message after I tap. Better still, put the real check first and only offer this after it."
 },
 {
  "id": "F08",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked on this screen for a way to do the 5-minute check on \"A shard\" and \"The shard key\".",
  "what_happened": "Each piece has a grey \"not yet\" pill, a cut-off line of text (\"A shard is one horizontal slice of a collectio…\") and an underlined link \"Reply\" with a speech bubble. The intro says \"Open a piece to read it, see what its level rests on and correct it: your correction wins…\". There is no \"check\", \"quiz\" or \"test myself\" anywhere.",
  "problem": "\"Reply\" sounds like answering a message, not testing myself, and \"not yet\" does not say not yet what. I cannot find the one thing I came for, and the intro paragraph is more small grey text that I skip. I would leave the app before the next class starts.",
  "severity": 4,
  "heuristic": "Match between system and the real world",
  "suggestion": "Put one big green button at the top: \"5-minute check · 3 questions\". On each piece, replace \"Reply\" with \"Test me\" and change \"not yet\" to \"Not checked\". Cut the intro paragraph to one short line."
 },
 {
  "id": "F09",
  "persona": "ayla",
  "run": 1,
  "screen": "stage1_course.jpg vs stage1_your_map.jpg",
  "what_i_tried": "After the course screen I went looking for \"Your map\", \"Homework\" and \"Your skills\".",
  "what_happened": "On the course screen the pill row is Content kernel / Experts / Rooms / From this ke…. On the map screen the pill row is Stream / Your map / Homework / Your skills / For…, with a smaller title above it. The map screen also repeats \"What a shard is\" as a heading in an empty white box and again as the section title just below. The last pill on each row is cut off, and I see no tap that took me from one row to the other.",
  "problem": "It feels like two different apps. I do not know how I got to \"Your map\" or how to get back. The empty box with the repeated title looks like a loading error. Because I lose my place when I am interrupted, I would be lost.",
  "severity": 2,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one row of tabs for the whole course, in the same place on every screen: Parts · My map · Homework · Skills. Highlight the one I am on, and remove the empty duplicate title box."
 },
 {
  "id": "F10",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I opened the app and looked at the first screen for today's parts and the 5-minute check.",
  "what_happened": "The top card is a big bold title, \"Chunk map retention raised to 14 days for the audit work\", with a paragraph about a config server and a rollback, dated 2026-08-11 and tagged \"Not yet checked\". Nothing says \"today\" and nothing says \"start your check\".",
  "problem": "I don't know any programming, so this reads like a foreign language, and I can't tell which of these cards are today's parts. It is a wall of text, so in two seconds I decide it isn't for me and leave.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Put one card at the very top that says \"Today: 3 new parts, 5-minute check\" with one big button, and write the card titles in plain words a first-year can follow."
 },
 {
  "id": "F11",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scrolled down the stream to see what else is new, tapping nothing, just reading the cards.",
  "what_happened": "One card shows \"nour  worked: the sharding exercise sheet  finished: 3 of 5  askedFor: an easier sheet\", which looks like raw data. Another says \"twin-ada\" and tells me I've opened a chapter three times. The badges say \"0 experts · 1 twin\" and \"! 3 experts · 1 twin\" with little coloured circles.",
  "problem": "I can't tell whose log this is, why I'm seeing it, or what \"experts\" and \"twin\" mean. The red \"!\" makes me think something is wrong with my account, and the \"you opened it three times\" message feels like being watched.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Write each card as a normal sentence with a name and a face, like \"Nour finished 3 of 5 on the sharding sheet\". Use plain badge words, and tap-to-explain for anything like \"twin\"."
 },
 {
  "id": "F12",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I tapped the big green \"Ask the library\" button at the bottom first, because it is the biggest coloured thing on the screen.",
  "what_happened": "The button is a very pale green with white text, so it looks switched off. The dark green arrow part on its right is the strongest colour, so my thumb goes there. The field above says \"Write something…\" and doesn't mention the library.",
  "problem": "I can't tell if the button works, which part of it I should tap, or whether I'm writing a post to the class or asking a question. I might send something to everyone by accident.",
  "severity": 2,
  "heuristic": "Visibility of system status",
  "suggestion": "Make the box say \"Ask the library about this course…\" and make the button a clear, solid colour once I've typed, greyed out only when it's empty. Put a plain label on the arrow part, or drop it."
 },
 {
  "id": "F13",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I opened the course page to find the first part and the 5-minute check.",
  "what_happened": "The first screen is all header: the title, a paragraph, \"Kept by the infra community\", \"Revision r14 · Sep 9, 2026\", \"12 elements in 4 chapters\" and \"KEEPERS\" with initials. Then a grey paragraph says \"Every element is checked on its own. There is no score for the kernel as a whole — three confirmed definitions say nothing about a contested caveat.\" The first real card is barely visible above the \"Ask the library\" box.",
  "problem": "I have to scroll past words like \"kernel\", \"keepers\" and \"revision\" that I don't know, and I skim the small grey text anyway. I still can't see where to start the check, so I give up.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Shrink the header to the title and one line, with a big \"Start today's check\" button right under it. Move revision, keepers and the explainer paragraph behind a small \"i\" icon."
 },
 {
  "id": "F14",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I looked for my homework and deadlines, moving between the course page and the \"Your map\" page.",
  "what_happened": "The course page has tabs \"Content kernel 12 · Experts · Rooms · From this ke…\", cut off at both edges. On the map page the tabs are \"Stream · Your map · Homework · Your skills · Fo…\", also cut off. The bottom bar has six icons with no words, and the badges say 11, 8, 8 and 4.",
  "problem": "Homework is in a different row on a different screen, and I have to remember which row it's in. The icons and numbers don't tell me what they count, and tabs cut off at the edge are easy to miss with one thumb.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one tab row on every screen with Homework and Check first, and keep it short enough to fit 390 px. Put a small word under each bottom icon, and say what each number counts."
 },
 {
  "id": "F15",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked at the piece \"A shard\" and tapped \"Reply\" to start testing myself, since there's no other button on the card.",
  "what_happened": "Each card shows the title, text cut off with \"…\", a grey outlined \"not yet\" pill and an underlined \"Reply\" with a speech bubble. \"A shard is one horizontal slice of a collectio…\" never finishes.",
  "problem": "\"Reply\" sounds like answering a message, not checking myself. I don't know if it opens a quiz, a chat, or a post to others. \"not yet\" looks like a button, but I'm not sure what it would do.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Call the button \"Test myself · 1 min\" and make it a proper solid button. Show \"not learned yet\" as a quiet label, not a pill that looks tappable."
 },
 {
  "id": "F16",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I tapped the outlined pill \"I understand everything here\" under \"What a shard is\", thinking it would just fold or skip the section.",
  "what_happened": "I can't see what it does. It says I understand everything, yet the counter reads \"0 of 3 learned\" and each piece says \"not yet\". There's no sign of what happens if I tap it.",
  "problem": "I could mark three things as learned without ever testing myself, and then the map would be wrong. If I'm interrupted, I couldn't tell whether I tapped it or not.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Rename it \"Mark all as learned\", ask me once to confirm, and show an \"Undo\" message afterwards."
 },
 {
  "id": "F17",
  "persona": "ayla",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked at the top of the map for a quick idea of how far I've got.",
  "what_happened": "Right under the tabs is a white box with \"What a shard is\" and a lot of empty space, and the same title appears again below as \"What a shard is  3 pieces\". Between them are \"0 learned · 0 worked on · 12 not yet\" and a long grey paragraph: \"From your own answers, not a score. Open a piece to read it, see what its level rests on and correct it: your correction wins and stays marked as yours…\".",
  "problem": "The empty white box looks broken or still loading. The paragraph is long, grey and full of abstract words, so I skip it and miss that I can open a piece or correct it. I'd like to see a picture of my progress instead.",
  "severity": 2,
  "heuristic": "Help and documentation",
  "suggestion": "Remove the empty box and show a small progress bar with \"0 of 12 learned\". Cut the explanation to one short line, like \"Tap a piece to read it or fix the level\", with a \"?\" for more."
 },
 {
  "id": "F18",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the first screen for anything that says what is due and when, so I could open it and skip the rest.",
  "what_happened": "I see four cards with no deadline on any of them: a bold post about chunk map retention, a block of raw fields from someone called nour, a note from twin-ada, and a post about graph storage. The bottom bar has six icons with no words and number badges (11, 8, 8, 4), and I ignore badges. Nothing says which icon is homework.",
  "problem": "I opened the app to find what is due, and the first screen can't tell me. I have to guess which unlabeled icon, or which tab on another screen, holds homework, and I give up fast when I have to hunt.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Pin a one-line strip at the top of the stream, such as 'Due: Sharding sheet, Fri 18:00 - open', that goes straight to the homework. Put text labels under the bottom icons."
 },
 {
  "id": "F19",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I opened the course page and looked for the list of what I have to know, and I skipped the text above it.",
  "what_happened": "The first half of the screen holds a description paragraph, 'Kept by the infra community', 'Revision r14', '12 elements in 4 chapters', a KEEPERS row with avatars, and another paragraph about how the kernel is checked and scored. The first real definition, 'A shard', starts at the bottom and the second is cut off.",
  "problem": "None of the header block helps me finish the homework, and it pushes the content below the fold on a phone. Words like 'kernel', 'elements' and 'keepers' are the app's own terms, so I can't tell what I'm looking at without reading the paragraph, and I won't read it.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Collapse the description, revision and keepers into a single line with an expand arrow, and drop the explanatory paragraph or tuck it behind a small '?'. The first definitions should be visible without scrolling."
 },
 {
  "id": "F20",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I looked for the Homework tab on the course page, since the 'Your map' screen has one.",
  "what_happened": "This row shows 'Content kernel 12', 'Experts', 'Rooms' and a clipped 'From this ke...', with a sliver of another pill cut off on the far left. On 'Your map' the row is different: Stream, Your map, Homework, Your skills, 'For...'. Homework is not visible here at all.",
  "problem": "The tab row changes between screens and scrolls sideways, so I can't tell whether Homework is one swipe away or doesn't exist here. I don't expect to swipe a row of tabs to find the one I need, and the clipped labels look broken.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one tab row on every course screen with the same order, starting with Homework, and make sure the labels fit without sideways scrolling, for example by shortening them or moving the secondary ones into a 'More' menu."
 },
 {
  "id": "F21",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I tried to scroll the list of definitions, and I expected the text box at the bottom to be a plain place to type.",
  "what_happened": "A panel headed 'Ask the library' sits permanently at the bottom with the line '- Just you and the library, about Sharding and data placement.', then a 'Write something...' box and a big pale-green 'Ask the library' button with a small split arrow. The same panel is on 'Your map', and the stream has a similar box without the heading. The panel takes about a fifth of the screen, and the button looks disabled.",
  "problem": "It shrinks the content I actually read, and the label appears twice, as heading and as button. The split arrow could be something like 'post to the group' or 'send privately', and I can't tell. In the stream, 'Write something...' could be a message to the whole class, and I assumed it works like Slack.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Collapse the panel into a single 'Ask the library' button or a one-line field that expands when tapped. Drop the duplicate label, and either remove the split arrow or label it with its options, such as 'Ask privately / Post to class'."
 },
 {
  "id": "F22",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I tapped the 'Everyone' dropdown at the top, expecting a filter like in Slack, such as 'only what is due' or 'only mine'.",
  "what_happened": "The dropdown says 'Everyone' and stays the same at the top of all three screens, including 'Your map', which is about my own progress.",
  "problem": "I can't tell what it filters, whether it is who I'm talking to, who sees what I write, or whose progress I see. On a personal screen 'Everyone' contradicts the title, so I don't trust what I'm looking at. It's also the first thing on every screen, so it adds one more thing I have to decide about.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Give it a visible label such as 'Showing: Everyone', and remove it from personal screens like 'Your map', where it doesn't apply."
 },
 {
  "id": "F23",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scanned the stream cards for the ones that concern me and skipped the long text.",
  "what_happened": "The second card is a block of raw fields from someone named nour: 'worked: the sharding exercise sheet', 'finished: 3 of 5', 'askedFor: an easier sheet', plus a timestamp. It has the same weight as the teacher's posts and the twin's note.",
  "problem": "It reads like a debug dump with raw field names. I can't tell whether it is mine, someone else's progress, or information for me, and if it is another student's data I'd be uneasy seeing it. It makes the stream look noisy.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Show it as a plain sentence, for example 'Nour finished 3 of 5 on the sharding sheet', or hide other students' activity from the stream unless I turn it on."
 },
 {
  "id": "F24",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I saw 'I understand everything here' right under the 'What a shard is' section and went to tap it to skip the three pieces, which I already know from other systems.",
  "what_happened": "The button looks like a small pill and sits next to the section title. Nothing on the screen says what happens after I tap it, whether the three pieces turn 'learned', or whether I can take it back. The text above says the map is 'from your own answers, not a score', which doesn't match a button that marks things as understood.",
  "problem": "One tap could mark a whole section as learned, and I can't see an undo. I might hit it by accident, or tap it and then not trust the 'learned' count.",
  "severity": 2,
  "heuristic": "User control and freedom",
  "suggestion": "After the tap, show a short toast like '3 pieces marked learned - Undo'. Show the effect in the button label, for example 'Mark all 3 as learned'."
 },
 {
  "id": "F25",
  "persona": "rune",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I tried to read the shard definition quickly from its card, and then I tried 'Reply' to see what it does.",
  "what_happened": "Each card shows its title, a 'not yet' pill and one line cut off with '...' ('A shard is one horizontal slice of a collectio...'), with an underlined 'Reply' link with a chat icon below. The text says 'Open a piece to read it', so reading needs a tap. A sticky white block with 'What a shard is' sits above the counts and looks empty.",
  "problem": "I can't read even a one-sentence definition without opening it. 'Reply' sounds like answering a thread in Slack, and I don't know whether it is for answering a quiz question or for writing to the teacher. The 'not yet' pill says nothing about what I'm supposed to do next.",
  "severity": 2,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Let the cards expand in place to show the whole definition with a one-tap self-check. Rename 'Reply' to what it does, such as 'Test myself', and make the 'not yet' pill say what is missing."
 },
 {
  "id": "F26",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the bottom bar for whatever opens my homework and deadlines, and had no idea which of the six icons that is.",
  "what_happened": "The bar has six icons with no words, and four of them carry green number badges (11, 8, 4, 8). The only hint of where I am is the pale green highlight on the second icon.",
  "problem": "I came for one thing, what is due. I would have to tap icons one by one to learn what each does. I ignore badges, so the numbers tell me nothing. On a phone on the train I would give up after two wrong taps.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Put a short word under each icon (Stream, Homework, Ask, ...), and show a due-soon count only on the Homework item."
 },
 {
  "id": "F27",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scanned the first screen of the Stream for the next deadline.",
  "what_happened": "The top of the stream is a news-style card about chunk map retention, then a classmate's status, then a message from twin-ada about a chapter I supposedly keep stopping at, then a card about graph storage. The newest item is dated 2026-08-24. Nothing says what is due or when.",
  "problem": "This reads like a team activity feed, not a course timeline. Everything is the same weight, so I have to read each card to find out whether it matters to me. None of it answers my question.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Pin a one-line strip at the top of the stream, such as \"Due: sharding exercise sheet, Fri 18:00, needs: A shard, The shard key\", with a link to the homework. Let the rest of the feed scroll below it."
 },
 {
  "id": "F28",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I tried to read the card from \"nour\" to see whether it is something I have to act on.",
  "what_happened": "It shows raw field names in bold: \"worked: the sharding exercise sheet  finished: 3 of 5  askedFor: an easier sheet\". There is no sentence and no button.",
  "problem": "It looks like a debug dump or a log line, so I do not trust it. I also cannot tell why a classmate's status is in my feed, or whether it is meant for me.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Show it as a plain sentence (\"Nour finished 3 of 5 on the sharding exercise sheet and asked for an easier one\"), or leave it out of my stream unless it affects me."
 },
 {
  "id": "F29",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scrolled down the stream to read the fourth card, about graph storage.",
  "what_happened": "The \"Ask the library\" composer is fixed at the bottom and takes about a quarter of the screen. The graph storage card is cut off in the middle of its timestamp and badges, and the feed behind it is squeezed. The big green button, with a ▾ half beside it, looks disabled because it is pale.",
  "problem": "I lose a quarter of a phone screen to a box I did not come to use, and the content I am reading is clipped behind it. I also cannot tell if the button works or what the ▾ opens. I would guess it behaves like GitHub's split button, which is a guess.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Collapse the composer to a single \"Ask the library\" icon or bar until I tap it. Keep the feed full height."
 },
 {
  "id": "F30",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I opened the course and looked for the things I must know for the homework.",
  "what_happened": "The first screen is a header: a title, a three-line description, \"Kept by the infra community\", \"Revision r14 · Sep 9, 2026\", \"12 elements in 4 chapters\" and a KEEPERS row. Under the tab chips there is another paragraph (\"Every element is checked on its own. There is no score for the kernel as a whole...\"). The first real element, \"A shard\", appears only after that.",
  "problem": "That is two paragraphs and a block of metadata that I skip, and I still have to scroll past them. I do not care who keeps it or which revision it is. The definition I need sits at the fold.",
  "severity": 3,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Cut the header to the title and one line of metadata, and move the explanatory paragraph behind a small \"?\" link. Have the first element visible without scrolling."
 },
 {
  "id": "F31",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I looked for \"Homework\" among the chips under the course title. It was one tap away on the other screen, so I expected it here too.",
  "what_happened": "This chip row reads \"Content kernel 12 · Experts · Rooms · From this ke...\", with a clipped chip on the left and a cut-off label on the right. The \"Your map\" screen has a different row: \"Stream · Your map · Homework · Your skills · Fo...\". Stream is also in the bottom bar. Homework does not appear in this row at all.",
  "problem": "The same place in the layout holds different tabs on different screens, and both rows run off the edge. So I cannot tell if Homework is anywhere on this screen or only appears after I switch to \"Your map\". It also means I would have to remember where each tab lives.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one tab row for the course, in the same order on every screen: Stream, Course, Your map, Homework. Keep it short enough to fit the screen width without scrolling."
 },
 {
  "id": "F32",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I wanted to check myself on \"A shard\", so I looked for a quick test or check button on the piece. The only action I could see was \"Reply\".",
  "what_happened": "Each piece shows its title, a text cut off with \"…\", a grey \"not yet\" tag and an underlined \"💬 Reply\". The page text says \"Open a piece to read it, see what its level rests on and correct it\".",
  "problem": "Where I come from, Reply means answering in a thread, as in Slack. Here I cannot tell if it starts a test, corrects the level, or messages someone. There is no obvious way to go from \"not yet\" to \"learned\". The definitions are cut off, so I cannot even see what I need to know without opening each one.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Replace \"Reply\" with a clear one-tap action such as \"Check me\", which opens a one-question test on that piece. Show the full definition, or its first sentence, instead of cutting it off at the first line."
 },
 {
  "id": "F33",
  "persona": "rune",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I would tap \"I understand everything here\", as I would \"Mark all as viewed\" on GitHub, to clear the three \"not yet\" tags in that section and get to what I really need.",
  "what_happened": "The pill sits under the section title with no hint of what it will change, and I see no confirmation or undo next to it. The page text says the map is \"from your own answers, not a score\", which does not match a button that lets me declare it myself.",
  "problem": "One tap could set three pieces to learned when I have not been checked on any of them. If I misjudge, my map shows a wrong picture and I do not know how to undo it. That is a risk for me, since I tend to move fast and skip the text.",
  "severity": 2,
  "heuristic": "Error prevention",
  "suggestion": "Make the button say what it does, such as \"Mark these 3 as known (without a check)\". Show a short Undo bar for a few seconds after the tap, and keep those pieces visibly marked as \"claimed\" rather than \"learned\"."
 },
 {
  "id": "F34",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scrolled the feed from the top to find the part of the day I missed, and read the second card, the grey one from \"nour\".",
  "what_happened": "The card shows what looks like raw field names: \"worked: the sharding exercise sheet  finished: 3 of 5\" and \"askedFor: an easier sheet\", with a classmate's name and a time.",
  "problem": "It reads like debug output, not like something written for me, and \"askedFor\" is one squashed word that I have to decode. It also shows me that another student asked for an easier sheet, so I immediately think my own \"I do not understand\" could appear like this in front of everyone, and I would stop using the app for questions.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Write these cards as a plain sentence (\"Nour finished 3 of 5 exercises\"), and do not show what one student asked for to other students at all, or show it without a name."
 },
 {
  "id": "F35",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I read the \"twin-ada\" card and tried to understand the little badges under each card (\"0 experts · 1 twin\", \"! 3 experts · 1 twin\", \"Not yet checked\").",
  "what_happened": "The card talks to me (\"You have opened the placement chapter three times ... without finishing it\") and is signed \"twin-ada\"; the badges mix a check mark, an exclamation mark, a dashed circle, a target-like circle and coloured initials, with no explanation.",
  "problem": "I do not know what a \"twin\" is (a person? a bot? a copy of Ada?), what an \"expert\" has done to a card (checked it? wrote it?), or whether \"!\" is a warning about my work or about the card. The same words appear with different symbols, so I cannot trust any of them. The card also tells me how often I got stuck, which feels like being watched.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Replace \"twin\" and \"experts\" with plain words, such as \"checked by 3 teachers\" and \"suggested by your study assistant\", and give each symbol one meaning, with a one-line legend I can open from the badge."
 },
 {
  "id": "F36",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the \"Everyone\" box at the top and at the sentence under \"Ask the library\" on the other two screens, to learn who would see what I write.",
  "what_happened": "The top of every screen says \"Everyone\" in a drop-down, while the bottom of the course and map screens says \"Just you and the library\"; on the stream screen the bottom has no such sentence at all.",
  "problem": "I take labels literally. \"Everyone\" next to a place where I type means that everyone sees it, and one screen later \"Just you\" says the opposite. I cannot tell whether \"Everyone\" is a filter for what I read or the audience for what I write, so I do not dare to type the thing I did not understand.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Give the drop-down a label that says what it does (for example \"Show cards from: Everyone\") and keep it separate from the writing box, which always carries its own visible line about who can see the text."
 },
 {
  "id": "F37",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I clicked into \"Write something...\" to type a question about the compound-key sentence, next to the pale \"Ask the library\" button with the small arrow.",
  "what_happened": "The box has no line saying who reads my text; the button is pale green like a disabled control and has a \"▼\" part whose options I cannot see; the box sits right over the last card so that card is cut off.",
  "problem": "\"Write something\" in a feed headed \"Everyone\" could be a post to the whole class. I would not know if I am sending it to the library or to the class before I press, and I cannot tell whether the pale button is disabled or just waiting for text. I would rather not type anything than send it to the wrong place.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Show a fixed line above the box on the stream screen (\"Only you and the library will see this\"), make the button look active when it can be pressed, and put the other send options as visible words, not behind \"▼\"."
 },
 {
  "id": "F38",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I read the top of the course page, then the tabs, to understand what \"Content kernel 12\" holds and whether it is the day I missed.",
  "what_happened": "I see \"Kept by the infra community\", \"Revision r14\", \"12 elements in 4 chapters\", \"KEEPERS\", a tab called \"Content kernel\" and the sentence \"There is no score for the kernel as a whole — three confirmed definitions say nothing about a contested caveat.\"",
  "problem": "These are terse labels and idioms that I read slowly in English: kernel, keepers, element, revision r14, \"contested caveat\". I do not know if \"element\" is the same thing as the \"pieces\" on the map screen, or why a missing score is explained with a sentence that I have to translate twice. I read carefully, so I lose time on words that should be obvious.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use the course's own words (\"Course content\", \"pieces\", \"maintained by ...\"), keep one word for one thing across screens, and rewrite the explanation sentence in short plain English without \"contested caveat\"."
 },
 {
  "id": "F39",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I looked at the bottom bar to go to my homework and to the place where my questions would be, and tried to scroll the tab row (Content kernel, Experts, Rooms, \"From this ke...\").",
  "what_happened": "The bottom bar has six icons with no words and number badges (11, 8, 8, 4) that mean nothing to me; the tab row is cut off at both edges, so \"From this ke...\" is unreadable and there is a hidden tab on the left; the course page also has its own \"Stream\" tab while the bottom bar has a stream icon.",
  "problem": "I open many things at once and lose which one I was answering, and with icon-only navigation I cannot tell the question-mark icon from the chat-bubble icon. I do not know if the badges are things I must do, or just counts. There seem to be two ways to reach the stream, so I am not sure where I am.",
  "severity": 2,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Put a short word under each bottom icon, say what each number counts (for example \"8 new\"), and show the tabs in a way that the full list is visible, with a clear marker for the current screen."
 },
 {
  "id": "F40",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I read the summary \"0 learned · 0 worked on · 12 not yet\" and the grey \"not yet\" tags on each piece to see where I am in the day.",
  "what_happened": "Every piece says \"not yet\", the summary has three states, and the explanation says \"From your own answers, not a score... your correction wins and stays marked as yours. Fold the sections you are done with; the map remembers.\"",
  "problem": "\"not yet\" could mean not read, not understood, not tested, or not checked, and the stream screen used \"Not yet checked\" for something else, so one label seems to mean two things. \"worked on\" is vague. I read the explanation three times and still do not know what my \"correction\" is or what is being corrected.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Name the three states with full phrases (\"not read yet\", \"read, not tested\", \"tested and understood\"), and replace the paragraph with two short sentences and one example of a correction."
 },
 {
  "id": "F41",
  "persona": "thao",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked for a way to say that \"A shard\" is unclear to me, and found \"Reply\" under each piece and \"I understand everything here\" under the section title.",
  "what_happened": "Each piece has a \"Reply\" link with a speech bubble, and the section has one pill \"I understand everything here\"; there is no matching \"I do not understand\" control and no sentence about who sees a reply.",
  "problem": "\"Reply\" suggests I am answering a message from someone, and I do not see a message. I would not guess that this is where I say \"unclear\". The only clearly worded button is the one that says I understand everything, and I do not see how to undo it for all 3 pieces; I fear that not pressing it, or pressing the other thing, shows the teacher that I am failing.",
  "severity": 3,
  "heuristic": "User control and freedom",
  "suggestion": "Change \"Reply\" to \"I am not sure about this\" with a small line \"Only you and the library see this\", and give \"I understand everything here\" a visible \"Undo\" after I press it."
 },
 {
  "id": "F42",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I started to read the stream from the top to find the day I missed, and I read the small grey card from \"nour\" and the card from \"twin-ada\".",
  "what_happened": "The nour card shows \"worked: the sharding exercise sheet  finished: 3 of 5  askedFor: an easier sheet\" with her name and a time, written like a database record. The twin-ada card tells me I opened a chapter three times \"without finishing it\".",
  "problem": "I now think that if I press anything unclear, a card like this will appear for me and everybody (and the teacher) will read \"askedFor: an easier sheet\". Nour's struggle is on display. Also the card says the app counts how many times I open a chapter. This is exactly what I am afraid of, so I would probably not ask at all.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Do not show other students' struggles with their names. Write cards in normal sentences, not field names like askedFor. Say clearly on a card whether it is only for me."
 },
 {
  "id": "F43",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I wanted to write my question in the box \"Write something...\" at the bottom, and first I looked at the \"Everyone\" box at the top.",
  "what_happened": "The top shows \"Everyone\" with an arrow, the box says \"Write something...\", and the pale green button says \"Ask the library\" with a small arrow beside it. Nothing on this screen says who will see what I write.",
  "problem": "\"Everyone\" is the first word on the screen, so I read it literally: what I write goes to everyone. \"Write something\" does not sound like asking the library, and the little arrow hides other choices I cannot see. On the course screen there is a line \"Just you and the library\", but here there is none. So the same box feels private in one place and public in another. I stop and do not write.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Under the box, always say in one plain line who sees the text (\"Only you and the library\" or \"All students\"). Name the \"Everyone\" box (for example \"Show posts from: Everyone\") and show the choices behind the arrow with words."
 },
 {
  "id": "F44",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the small coloured labels under each post to decide how much I can trust the claim (\"Not yet checked\", \"0 experts · 1 twin\", \"! 3 experts · 1 twin\").",
  "what_happened": "I see different small icons (empty circle, circle with dot, exclamation mark) and round initials that overlap each other (AC, BL, AT, LI cut off). The same \"experts · twin\" labels are also on the elements in the course screen, but there with a tick.",
  "problem": "I check claims, so this matters to me. I do not know what a \"twin\" is (a person? a bot? a copy of me?), what \"checked\" means, or why 0 experts is shown with a calm icon but 3 experts with an exclamation mark. The overlapping initials I cannot read. I cannot tell if the first post (\"Not yet checked\") is safe to learn from.",
  "severity": 3,
  "heuristic": "Help and documentation",
  "suggestion": "Explain the labels on tap (\"twin = an AI assistant that checked this\", \"experts = people who confirmed it\"). Use one clear set of words with a short legend, and show the avatars without overlap."
 },
 {
  "id": "F45",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the bottom bar to go from the stream to the course and then to my own map.",
  "what_happened": "The bottom bar has six icons without any words and with numbers on them (11, 8, 8, 4); two icons both show 8. On stage1_your_map.jpg none of the icons is highlighted, and at the top there is a tab called \"Stream\" while the bar also has a stream icon.",
  "problem": "I take things literally, and here there is nothing to take. I do not know what the 8 is (8 new? 8 unread? 8 to do?) and I cannot tell the two 8s apart. On my map screen I do not know where I am in the app. And \"Stream\" tab and the stream icon look like the same thing, but I do not know if they are. Late at night I do not want to guess.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "Put a short word under each icon (Course, Stream, Me, ...), explain the numbers (\"8 new\"), always highlight where I am, and use one name for the stream in only one place."
 },
 {
  "id": "F46",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I wanted to read what I missed, so I looked at the heading \"Sharding and data placement\", the tab \"Content kernel 12\", and the line \"12 elements in 4 chapters\".",
  "what_happened": "The course calls the content a \"kernel\" made of \"elements\" in \"chapters\". The map screen calls the same things \"pieces\" in \"sections\" (\"3 pieces\", \"12 not yet\").",
  "problem": "In data science \"kernel\" means a Jupyter kernel or a kernel function, so I first think this is code. And I count: elements or pieces? chapters or sections? If I cannot be sure that two words mean one thing, I start to doubt the numbers (12 here, 12 there, 4 chapters but I see only a section of 3).",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Use one set of ordinary words everywhere (for example \"parts\" and \"chapters\") and replace \"kernel\" with \"Course content\" or \"What to read\"."
 },
 {
  "id": "F47",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I tried to find the day I missed, the one with my teacher's post, and to read it in order.",
  "what_happened": "The course screen is organised by topic (\"What a shard is\", \"3\") with no days. The stream has posts of 11, 14, 18 and 24 August in one list, mixed with nour's record and a note for me, and no day headings. Some tabs are cut at the edges (\"From this ke...\" and one before \"Content kernel\").",
  "problem": "My goal is to read the day I missed properly. I cannot say \"show me yesterday\". I must scroll through other students' records and hidden tabs and guess by date, and I do not know which parts are new to me. It costs me many minutes of my 60-90.",
  "severity": 3,
  "heuristic": "Flexibility and efficiency of use",
  "suggestion": "Add a clear \"Catch up\" place that lists the parts by day (\"Yesterday: 4 new parts\"), with the unread ones marked, and keep other students' activity out of that list."
 },
 {
  "id": "F48",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I read the text under \"0 learned · 0 worked on · 12 not yet\" to understand what this map is for and what I am supposed to do on it.",
  "what_happened": "The text says \"From your own answers, not a score. Open a piece to read it, see what its level rests on and correct it: your correction wins and stays marked as yours. Fold the sections you are done with; the map remembers.\" Every piece has a grey \"not yet\" label.",
  "problem": "Too many idioms and short phrases for slow reading in English: \"what its level rests on\", \"your correction wins\", \"fold\", \"the map remembers\". Which level? Correct what? I have not answered anything yet, so \"from your own answers\" does not fit. And \"not yet\" on all 12 pieces looks like a failing mark, even if the text says it is not a score.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Write two or three short, literal sentences (\"This shows what you have learned. Tap a part to read it. If the app is wrong about you, change it.\") and use a neutral label like \"not started\" instead of \"not yet\"."
 },
 {
  "id": "F49",
  "persona": "thao",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked for a way to say that \"A shard\" is unclear to me, without anybody thinking I failed, and I saw \"Reply\" under each piece and the button \"I understand everything here\".",
  "what_happened": "\"Reply\" has a speech bubble but I do not see what I reply to. The \"Ask the library\" panel at the bottom stays open and says \"about Sharding and data placement\", the whole topic, not the piece. The single button \"I understand everything here\" sits above the three pieces, before I have read any.",
  "problem": "I open several things at once, so I will press Reply on \"The shard key\" and then lose which piece my text belongs to, because the panel at the bottom is about the whole topic. I also do not know if Reply goes to the library only, to the teacher, or to the group. There is no plain \"this is unclear\" choice, and one slip on \"I understand everything here\" would mark all three as learned, which is not true.",
  "severity": 3,
  "heuristic": "Recognition rather than recall",
  "suggestion": "When I press Reply, show the title of the piece inside the panel (\"Asking about: A shard\") and say who sees it. Add a private \"This is unclear\" action, and ask for confirmation, with an undo, before \"I understand everything here\"."
 },
 {
  "id": "F50",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the grey box at the top that says \"Everyone\" and expected it to be a way to pick a class or to see the page as a student sees it.",
  "what_happened": "It is a bare dropdown with the word \"Everyone\" and no label; the same box sits unchanged at the top of all three screens, and nothing says whether it filters the posts, switches the class, or switches who I am.",
  "problem": "I cannot tell which class or audience I am looking at, and I cannot tell whether what I see is exactly what students see tomorrow. That is the one thing I need to trust, so I would stop trusting the page.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Replace the bare dropdown with a labelled bar such as \"Viewing: Class 2B as a student sees it\" (with the class name written out), and make it obvious how to switch to the teacher view and back."
 },
 {
  "id": "F51",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I scrolled the stream for what students will see tomorrow and for any sign of which day and which course these posts belong to.",
  "what_happened": "The posts are stamped 2026-08-11 to 2026-08-24, they are about chunk map retention and graph storage while the course screen is about sharding, and nothing marks a post as \"today\", \"tomorrow\" or \"not yet published\".",
  "problem": "I forget which day and which class I am in, and here I have no way to find out; I cannot separate what is already visible to students from what comes tomorrow, so I cannot do my evening check.",
  "severity": 4,
  "heuristic": "Visibility of system status",
  "suggestion": "Show the course name on every post, group the stream by day with a clear \"Tomorrow, Thursday\" section, and mark posts that students do not see yet as \"scheduled\"."
 },
 {
  "id": "F52",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I read the grey card from \"nour\" to find out who is behind with the homework.",
  "what_happened": "It shows the raw lines \"worked: the sharding exercise sheet  finished: 3 of 5  askedFor: an easier sheet\" with a user name in lower case and no heading, like a data record.",
  "problem": "It looks like debug output, not a message, so I am not sure whether it is a note to me, a log, or something the whole class also reads; \"askedFor: an easier sheet\" about a named student worries me if everyone can see it.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Write it as a sentence a person would read (\"Nour finished 3 of 5 on the sharding exercise sheet and asked for an easier one\") and say on the card who can see it (\"only you\")."
 },
 {
  "id": "F53",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I read the card from \"twin-ada\" to see whether this is something a student receives or something about a student.",
  "what_happened": "The text speaks to \"you\" (\"You have opened the placement chapter three times this fortnight without finishing it\") but the card is labelled twin-ada and shows \"0 experts · 1 twin\", so I cannot tell who \"you\" is.",
  "problem": "If this is a private nudge to Ada, I should not be reading it in the same stream as the course posts; if it is addressed to me, it is wrong. I do not know whether students see each other's cards, and I do not know what a \"twin\" is.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Say on the card whom it is for and who sees it (\"Private note to Ada Torres, visible only to her and to you\"), and explain \"twin\" once in plain words the first time it appears."
 },
 {
  "id": "F54",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the small pills under each post (\"Not yet checked\", \"0 experts · 1 twin\", \"! 3 experts · 1 twin\" with coloured initials) to learn whether a post is reliable before students read it.",
  "what_happened": "The pills use \"experts\" and \"twin\" with no explanation, the exclamation mark on the last one is not explained, and \"Not yet checked\" does not say by whom or what happens next; the avatars overlap and are cut off at the bottom edge.",
  "problem": "I cannot tell if \"Not yet checked\" means students are already reading unchecked material, and I do not know if the red \"!\" with 3 experts is good or bad news.",
  "severity": 2,
  "heuristic": "Match between system and the real world",
  "suggestion": "Use plain words (\"Checked by 3 teachers; 1 student asked about it\") and a short \"What do these mean?\" link under the first pill; for \"Not yet checked\" say whether students can see it."
 },
 {
  "id": "F55",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I opened \"Your map\" to see who in the class is behind.",
  "what_happened": "It says \"0 learned · 0 worked on · 12 not yet\" and \"From your own answers, not a score\", so it is my own empty map as if I were a student, with no class or student list.",
  "problem": "I read the zeros first as the class being completely behind; only after the second sentence do I understand it is about me. The screen does not answer my question and I cannot find where the class overview is.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "For a teacher, put a class overview first (\"6 of 24 students have not started this chapter\", with names) and keep the personal map behind a clearly named \"My own map\" tab."
 },
 {
  "id": "F56",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I wanted to open pieces and look at them, but I was afraid to touch anything: there is \"I understand everything here\" under the section title and a \"Reply\" link on each piece.",
  "what_happened": "Both are live-looking controls with no sign that I am only previewing; the text says \"your correction wins and stays marked as yours\", which tells me my clicks are recorded.",
  "problem": "My goal is to check without editing anything by accident; a wrong tap could mark a section as understood or post a reply under my name, and I would not know how to undo it.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Add a visible \"Preview, nothing you do here is saved\" banner in teacher view and disable or ask for confirmation on \"I understand everything here\" and \"Reply\"."
 },
 {
  "id": "F57",
  "persona": "petra",
  "run": 1,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I tapped into the \"Write something...\" box at the bottom of the stream, wondering whether it would post to the class or to the library.",
  "what_happened": "On the stream the box has no caption and the green \"Ask the library\" button has a small arrow beside it; on the two course screens the same box has the line \"Just you and the library, about Sharding and data placement\", which the stream does not have.",
  "problem": "I cannot tell whether my text would go to all students or only to the library, and the arrow may hide another choice such as posting to everyone; one slip and a half-written teacher remark is in front of the class.",
  "severity": 3,
  "heuristic": "Consistency and standards",
  "suggestion": "Show the same caption on the stream as on the other screens, and when the arrow offers \"Post to everyone\" label it plainly and ask \"Send to the whole class?\" before posting."
 },
 {
  "id": "F58",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I looked at the grey 'Everyone' dropdown at the top first, to find out which class and which day I am looking at.",
  "what_happened": "It only says 'Everyone', with no label. The stream below has entries dated 2026-08-11, 2026-08-14, 2026-08-18 and 2026-08-24, and nowhere does it say which class or which session day this is.",
  "problem": "I forget which class or day I am on, and here nothing reminds me. 'Everyone' could be a class, an audience filter or a student picker. I cannot tell whether I am looking at tomorrow's session or at last month.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Replace the bare 'Everyone' with a fixed header line such as 'Class DS402-A · Session of 30 Oct · viewing as teacher', so the class, the day and my role are always readable."
 },
 {
  "id": "F59",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I tried to read the stream as a list of who is behind, starting with the entry from 'nour'.",
  "what_happened": "The entry reads 'nour  worked: the sharding exercise sheet  finished: 3 of 5  askedFor: an easier sheet'. It uses lowercase names and raw field names such as 'askedFor', and it is in a different style from the cards around it. A card from 'twin-ada' says 'You have opened the placement chapter three times this fortnight without finishing it'.",
  "problem": "The raw labels look like a database record, not like something meant for a teacher. The 'twin-ada' card is written to a student ('You have opened...'), so I cannot tell whether I am seeing a student's private nudge or the shared stream. If students can see each other's notes I have a privacy problem. If they cannot, I am not seeing what they see.",
  "severity": 3,
  "heuristic": "Match between system and the real world",
  "suggestion": "Write these entries in plain sentences ('Nour finished 3 of 5 on the sharding sheet and asked for an easier one'). Mark every entry with who can see it, for example 'Only Ada sees this' or 'Whole class'."
 },
 {
  "id": "F60",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_stream.jpg",
  "what_i_tried": "I wanted to scroll to the bottom to see the newest entries, and I looked at the 'Write something...' box and the 'Ask the library' button.",
  "what_happened": "The newest card ('Graph storage moved to per-record documents', 3 experts · 1 twin) is cut off by the pinned composer. The composer has one big button, 'Ask the library', with a small arrow. The box is labelled 'Write something...' while the stream header says 'Everyone'.",
  "problem": "I cannot tell whether typing here posts to the stream for everyone, or sends a private question to the library. The arrow suggests other actions but does not say which. Writing to the whole class by accident is exactly what I am afraid of.",
  "severity": 3,
  "heuristic": "Error prevention",
  "suggestion": "Make the composer say where the text goes ('Private question to the library' or 'Post to the whole class'). Make the other action a clearly named second button, or ask me to confirm before anything is posted to everyone."
 },
 {
  "id": "F61",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I looked for what students will see tomorrow, and read 'Revision r14 · Sep 9, 2026' and '12 elements in 4 chapters'.",
  "what_happened": "The page shows a revision number and a date, and 'Kept by the infra community' with initials under KEEPERS. It does not say whether r14 is what students have now or what they get tomorrow, or whether there are unpublished changes.",
  "problem": "I assume that what I see is what students see, and this page does nothing to confirm it. 'Revision' and 'Kept by' read like a wiki, not like a teacher's session. If r15 goes out tonight I would not know.",
  "severity": 3,
  "heuristic": "Visibility of system status",
  "suggestion": "Add a line such as 'Students see revision r14 now. Tomorrow's session shows r14 (no changes since)'. Give me a clear 'See as a student' switch with a banner saying so while it is on."
 },
 {
  "id": "F62",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_course.jpg",
  "what_i_tried": "I read the grey text under the tabs and the green badges on each element ('3 experts · 1 twin') to judge whether the content is ready.",
  "what_happened": "The text says 'Every element is checked on its own. There is no score for the kernel as a whole — three confirmed definitions say nothing about a contested caveat.' The badges show '✓ 3 experts · 1 twin' with several overlapping round initials, such as AC, BL, AT and LI, that run into each other.",
  "problem": "I know the content perfectly, but I do not know what a 'kernel', a 'twin' or a 'contested caveat' mean here, or whether '3 experts' is good enough for tomorrow. The overlapping initials cannot be read.",
  "severity": 2,
  "heuristic": "Help and documentation",
  "suggestion": "Use plain words ('checked by 3 teachers', 'a student's digital assistant' or whatever twin means). Add a small 'What does this mean?' link next to the badge. Let me tap the badge to see the full names."
 },
 {
  "id": "F63",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I opened 'Your map' hoping it would show who in the class is behind.",
  "what_happened": "It shows '0 learned · 0 worked on · 12 not yet' and the text 'From your own answers, not a score.' So it is my own map, and it is empty. The tab row (Stream, Your map, Homework, Your skills, 'Fo...') has nothing for the class.",
  "problem": "My main goal tonight is to see who is behind. This screen is built for a student, and I get my own blank progress instead, which looks as if I had never studied my own course. I can find no view of the class anywhere in the three screens.",
  "severity": 4,
  "heuristic": "Match between system and the real world",
  "suggestion": "For a teacher, show the class on this tab: each element with how many students have learned it, and a short list of the students who are behind, which I can open one by one."
 },
 {
  "id": "F64",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I considered tapping 'I understand everything here' and 'Reply' on 'A shard' to see what a student sees.",
  "what_happened": "The page says 'Open a piece to read it, see what its level rests on and correct it: your correction wins and stays marked as yours.' There is an 'I understand everything here' button and a 'Reply' link on every piece. Nothing says that I am only previewing.",
  "problem": "I am afraid to click anything. A tap on 'I understand everything here' might mark all 3 pieces as learned, and a 'correction' might change something the students see. I cannot tell what is undoable.",
  "severity": 3,
  "heuristic": "User control and freedom",
  "suggestion": "Give teachers a read-only preview where these buttons are greyed out and a banner says 'Preview: nothing here is saved'. Add a visible 'Undo' after any action that marks or corrects something."
 },
 {
  "id": "F65",
  "persona": "petra",
  "run": 2,
  "screen": "stage1_your_map.jpg",
  "what_i_tried": "I looked at the top of the page to see how the sections are organised.",
  "what_happened": "A white box with only 'What a shard is' sits above the counts, and the same heading comes back below with a triangle and '3 pieces · 0 of 3 learned'. The tabs are cut off at the right edge ('Fo...'). The text on each piece is truncated ('A shard is one horizontal slice of a collectio...').",
  "problem": "It looks broken or duplicated, so I trust it less. I cannot see all the tabs, and I cannot read a whole definition without opening it, which is slow when I want to check many pieces tonight.",
  "severity": 2,
  "heuristic": "Aesthetic and minimalist design",
  "suggestion": "Show the section heading once. Let the tab row wrap or show a 'more' menu so no tab is hidden. Show the full first sentence of each piece, or give me an 'expand all' control."
 }
]