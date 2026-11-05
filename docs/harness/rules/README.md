# How this rules pack works

A short set of rules for building a small product, and the data analysis behind it, with coding agents.
It condenses the rules our team works by every day: the analytics rules and the development rules.

| File | What it is | Loaded |
|---|---|---|
| `AGENTS.md` | the starter for your repo root: product directive, always-on core, "read when" table | every turn |
| `CLAUDE.md` | one line, `@AGENTS.md`, so Claude Code follows the same file | every turn |
| `ANALYTICS_RULES.md` | A1–A14: data, numbers, findings, recommendations | on demand |
| `DEVELOPMENT_RULES.md` | D1–D15: requirements, plans, code, tests, commits | on demand |
| `USER_SIMULATION.md` | U-rules: simulated users critique your product in waves, under a fixed directive | on demand |
| `SOURCES.md` | which of our team's rules each student rule condenses | for the teacher |

## What a rule file is

A rule file is a Markdown file that the agent reads as standing instructions.
It is not code: nothing enforces it except the agent reading it and you checking.
Each rule in this pack has four parts:

- an **id** (A3, D8), so the agent can cite it and you can search for it;
- the **rule**, one imperative line;
- **why**, one line, so the agent can apply the rule to cases it did not foresee;
- a **check**, one line, so you or the agent can verify that it was followed.

## Always on, or loaded on demand

Every line an agent reads takes context (the model's working memory) on every turn.
So rules come in two kinds:

- **Always on.** Read at the start of every session and present in every turn. Keep this part short:
  the product directive, five to eight core rules, and a map of the rest.
- **On demand.** Read only when the task needs it. There are three triggers:
  - **by description**: the file carries a one-line description, and the agent opens it when the task matches;
  - **by file glob** (a file pattern like `*.ipynb`): it loads when the agent touches a matching file;
  - **manually**: you name the file in your message.

Our team's loading rule estimates that loading by task saves 30–40% of the rule text per request. Yet in practice
24 of our 36 development rule files are still marked always on: the always-on part grows unless someone prunes it.

| | Claude Code | Codex, and others reading `AGENTS.md` | Cursor | Windsurf | GitHub Copilot |
|---|---|---|---|---|---|
| Always on | `CLAUDE.md` at the root, and every file it imports with `@path` | `AGENTS.md` at the root | `.cursor/rules/*.mdc` with `alwaysApply: true` | `.windsurf/rules/*.md` with `trigger: always_on` | `.github/copilot-instructions.md` |
| By description | a skill: `.claude/skills/<name>/SKILL.md` with `description:` | none: use the "read when" table | `description:` with `alwaysApply: false` | `trigger: model_decision` with `description:` | none: use the "read when" table |
| By file glob | `.claude/rules/*.md` with `paths:`, or a `CLAUDE.md` inside that folder | an `AGENTS.md` inside that folder | `globs:` | `trigger: glob` with `globs:` | `.github/instructions/*.instructions.md` with `applyTo:` |
| Manually | `@rules/ANALYTICS_RULES.md` in your message | name the file in your prompt | `@` and the rule's name | `trigger: manual`, then `@` and its name | attach the file to the chat |

Five details that catch people out:

- An `@path` import in `CLAUDE.md` is not on demand. It loads at the start, like `CLAUDE.md` itself.
- Recent Claude Code versions read `AGENTS.md` when there is no `CLAUDE.md`. The one-line `CLAUDE.md` works in all versions.
- Cursor ignores a plain `.md` file in `.cursor/rules/`. The extension must be `.mdc`.
- ChatGPT in the browser sees no repo. Paste `AGENTS.md` into the project's instructions and attach the rules files.
- Windsurf is now Devin Desktop: it reads `.devin/rules/` and falls back to `.windsurf/rules/`; one rule file holds at
  most 12,000 characters.

## The navigation file

The always-on part should be a map, not the territory.
Our team's always-on navigation file says "load this file first", sorts the task into a kind
(new feature, bug, tests, interface, shared infrastructure, review), and lists the rules to load for each kind.
Two lessons from running it:

- Trigger on intent ("add a chart"), not on file names. People rarely name files.
- Prefer loading a rule you did not need to missing one you did. A wasted read costs seconds; a missed rule costs hours.

In this pack the map is the "Read when" table in `AGENTS.md`.
It also makes on-demand loading work in tools with no description trigger:
the agent reads the table every turn and opens the file it points to.

## Install it in your repo, in three steps

1. Copy `AGENTS.md` and `CLAUDE.md` to the root of your repo. Copy `ANALYTICS_RULES.md`, `DEVELOPMENT_RULES.md`
   and `USER_SIMULATION.md` into a `rules/` folder.
2. Write the product directive in `AGENTS.md` yourself, not the agent, and fill in the project facts. Commit.
3. Open a new agent session and ask: *"Which rules in this repo apply if I ask you to add a chart of sign-ups per
   week? Cite their ids."* Expect something like A3, A12, D7 and D8. If the agent cites nothing, it did not load the
   file: check the file's name and place.

Optional: load the two rules files on demand in your tool instead of through the table. One wrapper per file:

```markdown
---
description: Data analysis rules. Load when collecting data, computing numbers or writing findings.
alwaysApply: false
---
Follow rules/ANALYTICS_RULES.md.
```

That is `.cursor/rules/analytics.mdc`. For Windsurf, save it as `.windsurf/rules/analytics.md` and replace
`alwaysApply: false` with `trigger: model_decision`. For Claude Code, save it as `.claude/rules/analytics.md` with
`paths: ["**/*.ipynb", "analysis/**"]` instead of the two keys. Do the same for the development rules.

## How to tell the rules work

- Make the agent cite the id of every rule it applies, e.g. "(D8)". `AGENTS.md` already asks for this.
- Read the citations for a week. A rule that is never cited is dead weight: it costs context every turn and
  changes nothing. Delete it, or move it from always on to on demand.
- A rule that is cited but still broken needs a sharper check, not louder wording.
- Keep the always-on part short. Our team caps any rule file at about 300 lines; for your repo, keep `AGENTS.md`
  under 60.
- After every rule change, ask the question from step 3 again and compare the answers.

## How to change a rule

- One rule, one reason. Do not rewrite the pack because one thing went wrong.
- Change the rule in the same commit as the change that made it necessary (the bug fix, the incident),
  and say why in the commit message. Our team's newest rules open with the request that created them, quoted and dated.
- When several people share the rules, one person edits the master copy, and everyone copies it unchanged.
- Once a month, review the always-on part: drop what is never cited, and add a trigger where a rule was missed.
