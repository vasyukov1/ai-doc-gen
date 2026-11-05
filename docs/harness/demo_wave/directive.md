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
