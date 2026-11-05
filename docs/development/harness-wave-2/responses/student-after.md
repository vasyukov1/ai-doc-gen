# Mira — simulated after-change rewalk

This is a fictional phone review, not a real participant or independent live browser run (U8). Evidence: inspected after-phone.png and student-status.png; read index/job/sections/document templates, setup.ts and profile.ts; read checks/walk.log. The log reports C1–C3 PASS and C5 student completed with 17 baseline sections, zero model sections and 17 unresolved gaps. It does not separately document every edit/rebuild action.

| Step | Result and evidence |
|---|---|
| Choose PZ | PASS, screenshot/source: visible “FKN coursework / ПЗ” shortcut; Russian guidance explains implementation and rationale and asks for source ТЗ. Log reports purpose and Russian guidance PASS. |
| Mistaken PDF requirements | PASS for recovery guidance in source: specification panel explicitly says to copy readable PDF/Typst text; upload accepts txt/md. The panel is closed in the supplied phone screenshot, so actual rejected-file recovery remains unverified here. |
| Expect key-free finished prose | PASS for honest outcome: screenshot says outline placeholders only, zero model-improved sections, unresolved gaps, and “Outline PDF”. Creation guidance also distinguishes outline from prose in Russian. |
| Leave and restore text | PASS by source plus reported C3 check: explicit opt-in allowlist restores text and model IDs; feedback asks to reselect files and enter a key. Screenshot shows the remember/forget controls. I did not independently navigate away and restore. |
| Inspect/edit/rebuild | PASS for discoverability: status screenshot exposes available sections; source exposes review/rebuild link and warns that saving edits does not update PDF. Actual saved-edit-to-new-PDF behavior is not established by these screenshots or the terse log. |
| Invariant: keys excluded from browser draft | PASS within reviewed feature: setup.ts stores only explicitly allowed text/model fields, with no credential fields, files or confirmations. C3 log reports secret exclusion. This is not a blanket audit of all browser persistence. |

Remaining simulated hypotheses (defer; no new accepted changes):

- S-A1: Creation on phone / find restoration before leaving / remember control sits below all provider cards in the long screenshot / Mira may leave before discovering opt-in and lose text / severity 2 / Recognition rather than recall / place a compact save-state control near the first text input, while preserving explicit opt-in (I1, I4).
- S-A2: Status and section editing / interpret and correct the outline / core warnings and save/edit/rebuild explanations remain mostly English despite Russian document language / Mira with limited English may miss what still needs work / severity 2 / Match between system and the real world / extend Russian guidance to outcome and editing actions (I2).

The bounded workflow makes source-to-editable-document steps visible; superiority over a general coding assistant and real-user success are unmeasured. No model API or application provider key was used.
