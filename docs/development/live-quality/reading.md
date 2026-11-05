# Q5 — Stable section reading

Date: 2026-11-06. Part of the five-change [live quality plan](plan.md), not a new milestone.

## Requirement, action and check index

| ID | Requirement and action | Manual check |
|---|---|---|
| Q5.1 | Read intermediate saved sections without forced navigation. Remove the three-second meta refresh; show explicit latest-version links instead. | Open a running job's sections, read below the first section for more than three seconds, and verify no navigation or scroll reset. Explicitly refresh to see worker updates. |
| Q5.2 | Jump to a named section using keyboard or touch. Add a named navigation landmark, native heading anchors, section-specific textarea labels and return links. | Activate a chapter link with the keyboard, inspect heading focus and viewport position; return to the list. Repeat at laptop and phone widths. |
| Q5.3 | Explicit refresh must preserve reading position and never silently discard pending prose. Guard dirty forms and in-flight actions; checkpoint only section ID and scroll coordinates in the current history entry. | Edit two sections; refresh must stay on the page and retain both values. Save one; refresh must still be blocked. Save the other, then refresh from a section footer; verify the visible section and relative position remain stable. |
| Q5.4 | Separate saved availability, quality and rendered output. Explain that saved text can be an outline or contain gaps, and that saving prose requires a separate PDF rebuild. | Inspect running, outline and completed pages. Ensure section count never implies verified content and the PDF rebuild link remains explicit. |

## Decisions

- Use explicit refresh instead of live content replacement. Appending or replacing growing prose would move the reader and require reinitializing editor/save behavior. The worker continues independently, while the user chooses when to load saved changes. No new endpoint or polling loop is needed.
- Keep the server's existing editing rules: processing sections remain read-only. The page becomes editable only after an explicit refresh observes a stopped run. The helper reads the existing form module's `data-unsaved` and `aria-busy` state; it does not clear either state or change save requests.
- Store only the current section heading ID and numeric coordinates in the history entry. No prose, uploaded data or credentials are persisted by this helper. Consume the checkpoint after reload. If browser history writes are restricted, ordinary browser reload/scroll restoration remains the fallback.
- Preserve native anchors and refresh URLs without JavaScript. JavaScript adds dirty-form protection and position restoration; a browser with JavaScript disabled uses ordinary explicit navigation and HTML form submission. No automatic refresh exists in either mode.
- Long navigation lists remain a simple responsive list. Sidebar layouts, collapsible chapter grouping, live quality assessment and version history stay deferred.

## Verification evidence

- `npm run typecheck` in `web/`: exit 0, `tsc --noEmit` on 2026-11-06.
- `go build ./...`: exit 0 on 2026-11-06.
- `curl -s -I http://localhost:8083/`: HTTP 200 from the isolated preview. A Chrome walkthrough of `/jobs` showed the empty browser-owned document list. It did not establish Q5 browser acceptance: the existing service embeds the prior template/assets, and this browser owner has no retained preview job.
- No model calls, keys, service restarts, automated test files, suites or commits were used for Q5.

## Integration and remaining checks

The primary agent integrates `import { initSectionsReading } from "./sections"` and calls it after `initForms()` and `initTextareas()` in `web/src/main.ts`, then rebuilds static assets and the preview binary. The primary owns service restarts.

Phone/laptop screenshots, keyboard navigation, running-page dwell, dirty-edit blocking, completed-save refresh and restored scroll remain **unverified in the integrated browser** until that rebuild. Perform the Q5.1–Q5.4 actions above against disposable owned preview jobs. Deliberately exercise late save responses and preserve newer edits before attempting refresh. Screenshots should contain only synthetic data.

Resume command: `npm --prefix web run typecheck`, then run the indexed manual checks against the primary agent's rebuilt preview.

## Integrated checks (2026-11-06)

The primary rebuilt assets and the isolated 8084 preview. [Browser checks](checks/browser.log) passed at 1366×900 and 390×844: no meta refresh, a stationary reading interval after anchor scrolling settled, explicit refresh retained the reading position, and no horizontal overflow. [Guard checks](checks/reading-guards.log) confirmed two unsaved textareas survive blocked refresh, an in-flight form blocks refresh, and a JavaScript-disabled footer link performs a fresh document request. In-flight state was injected in the DOM for this focused check; it is not evidence of a real slow provider/save race. No edited text was submitted in these checks. Earlier wave-1 evidence covers save-response/edit-revision behavior.

Screenshots are private local artifacts in `tmp/live-quality/screens/` because the generated prose/specification contains the supplied fixture's identifying text; they are not committed. No unseen browser outcome is claimed.
