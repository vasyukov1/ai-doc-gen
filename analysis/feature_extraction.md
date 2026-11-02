# Feature and pain extraction from feedback

Input: one dated user interview excerpt, review, or public post, with source and consent status. Use only the text given. Candidate features: repository import, archive upload, file selection, evidence links, standards profile, institution template, section editing, Typst download, PDF download, privacy controls, and model choice.

For each item, identify features explicitly mentioned and sentiment (`+`, `-`, or `0`). Extract the user's actual job, workaround, desired output, and any stated audience. Never infer age, role, country, or willingness to pay from writing style. Return JSON:

```json
{"features":[],"sentiment":{},"job":null,"workaround":null,"desired_output":null,"group":null,"quote":null}
```

The quote must be an exact excerpt of at most 15 words and must not identify the speaker. A human should label a sample independently and compare agreement by feature before using aggregate counts. Do not treat multiple posts from one user or one marketing campaign as independent demand signals.
