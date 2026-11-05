"""Recompute simulated-review agreement from the saved raw responses (no API calls)."""
import json
from pathlib import Path

root = Path(__file__).parent
report = {"kind": "simulated review agreement, not real-user frequency", "personas": {}}
for persona in ("beginner", "maintainer", "reviewer"):
    runs = []
    for run in (1, 2):
        findings = json.loads((root / "responses" / f"{persona}-{run}.json").read_text())
        pairs = set()
        for finding in findings:
            screen = finding["screen"]
            screen = "sections" if "sections" in screen else "creation" if "index" in screen else "status"
            pairs.add((screen, finding["heuristic"]))
        runs.append(pairs)
    shared, union = runs[0] & runs[1], runs[0] | runs[1]
    report["personas"][persona] = {
        "shared_pairs": sorted(shared), "union_pairs": sorted(union),
        "jaccard": len(shared) / len(union) if union else None,
    }
print(json.dumps(report, indent=2))

# Triage is analyst-authored; these assignments are not model confidence scores.
assignments = {
    "beginner-1-1": ("DEFER", "R5"), "beginner-1-2": ("ACCEPT", "R1"),
    "beginner-1-3": ("ACCEPT", "R3"), "beginner-2-1": ("DEFER", "R5"),
    "beginner-2-2": ("ACCEPT", "R1"), "beginner-2-3": ("ACCEPT", "R3"),
    "maintainer-1-1": ("ACCEPT", "R2"), "maintainer-1-2": ("ACCEPT", "R2"),
    "maintainer-1-3": ("DEFER", "R2"), "maintainer-2-1": ("ACCEPT", "R2"),
    "maintainer-2-2": ("ACCEPT", "R2"), "maintainer-2-3": ("DEFER", "R2"),
    "reviewer-1-1": ("ACCEPT", "R3"), "reviewer-1-2": ("ACCEPT", "R5"),
    "reviewer-1-3": ("ACCEPT", "R3"), "reviewer-2-1": ("ACCEPT", "R3"),
    "reviewer-2-2": ("ACCEPT", "R5"), "reviewer-2-3": ("ACCEPT", "R3"),
}
indexed = []
for persona in ("beginner", "maintainer", "reviewer"):
    for run in (1, 2):
        rows = json.loads((root / "responses" / f"{persona}-{run}.json").read_text())
        for number, row in enumerate(rows, 1):
            key = f"{persona}-{run}-{number}"
            decision, change = assignments[key]
            indexed.append(dict(row, id=key, persona=persona, run=run, decision=decision,
                                change=change, reason="See triage.md for the invariant, mitigation and deferred scope."))
(root / "findings.json").write_text(json.dumps({"as_of": "2026-11-05", "findings": indexed,
    "spread": report["personas"]}, indent=2, ensure_ascii=False) + "\n")
