"""Normalize reviewer screen/heuristic pairs; agreement is not user prevalence."""
import json
from pathlib import Path

# Analyst normalization of retained verbatim reviews, not additional model findings.
pairs = {
    'student': [
        [('creation', 'error recovery'), ('creation', 'consistency'), ('creation', 'error prevention'), ('sections', 'recognition')],
        [('creation', 'error recovery'), ('creation', 'consistency'), ('creation', 'error prevention'), ('sections', 'recognition')],
    ],
    'freelancer': [
        [('creation', 'real-world match'), ('creation', 'consistency'), ('creation', 'help'), ('sections', 'recognition')],
        [('creation', 'real-world match'), ('creation', 'consistency'), ('creation', 'help'), ('sections', 'recognition')],
    ],
    'employee': [
        [('creation', 'error prevention'), ('status', 'visibility'), ('sections', 'consistency'), ('sections', 'visibility'), ('sections', 'recognition')],
        [('creation', 'error prevention'), ('status', 'visibility'), ('sections', 'consistency'), ('sections', 'consistency'), ('status', 'error recovery')],
    ],
}
result = {}
for persona, runs in pairs.items():
    first, second = map(set, runs)
    result[persona] = {'raw_findings_per_run': list(map(len, runs)), 'normalized_pairs': runs,
                       'intersection': sorted(first & second), 'union': sorted(first | second),
                       'agreement': len(first & second) / len(first | second)}
Path(__file__).with_name('agreement.json').write_text(json.dumps(result, indent=2) + '\n')
