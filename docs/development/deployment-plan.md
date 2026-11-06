# Railway deployment plan

Date: 2026-11-06. Scope: the user requested deployment of ai-doc-gen using the supplied Streamlit teaching material as a guide, and selected Railway with web and worker as separate processes in one container plus PostgreSQL. Google sign-in is explicitly deferred. This is a deployment increment, not authorization to start milestone 8 or the architecture handoff waves.

## Design and alternatives

Use the existing Go application and Typst image. A dedicated Railway image starts migrations, then supervises separate web and worker processes sharing one persistent /data volume. PostgreSQL is a separate private service. One replica; a volume deployment can have downtime. Preserve browser ownership and encrypted per-job provider connections; do not transfer local jobs or credentials.

A VPS with Docker Compose would preserve separate containers but adds server administration; the user selected Railway. Separate Railway application services would require a shared-storage redesign; defer object storage and horizontal scaling. Do not copy the unrelated Streamlit example into this repository.

The launcher must fail closed on migration or child failure, forward termination, stop the other child, and exit for platform restart. Initialize the volume as root, then run the application as UID 10001. PORT maps to HTTP_ADDR. Enable Secure owner cookies through explicit deployment configuration; do not trust arbitrary forwarded headers. Retain local defaults.

## Requirement / action / check index

| ID | Action | Acceptance |
| --- | --- | --- |
| DEP1 | Add a Railway image and supervised launcher | Build; separate web/worker; graceful stop; child failure exits container |
| DEP2 | Persist application data and PostgreSQL | Fresh volume is writable as app; restart preserves job and artifact |
| DEP3 | Configure HTTPS runtime and safe packaging | PORT respected; Secure cookie option; private files absent from build context |
| DEP4 | Publish and smoke check | Public HTTPS home/readiness; keyless repository job produces downloadable PDF; no claim of model prose quality |
| DEP5 | Document operations and remaining work | Variables, backups, rollback limits, domain and deferred Google login documented |

## Execution

- [x] Inspect attached files, application lifecycle and official Railway documentation.
- [x] Confirm target and topology; defer Google login by user choice.
- [ ] Implement DEP1–DEP3 and deployment instructions.
- [ ] Run go build and isolated Docker smoke checks, without automated test files.
- [ ] Review changes, commit and push task-owned paths.
- [ ] Configure Railway PostgreSQL and persistent application service; check actual plan limits before chargeable actions.
- [ ] Deploy, validate HTTPS and keyless document flow, record URL and limitations.

## Review focus

Migration failure must prevent serving. Either child exiting must stop its sibling. Termination must reach both children. Mounted volumes begin owned by root. Private local artifacts must never enter the image. No local database or provider secret may be uploaded. Initial keyless output is an outline, not proof of real-model quality.

## Resume point

Implementation pending. Railway GitHub authorization remains on screen. Next local command: `go build ./...`. Deployment requires successful account authorization and available resources; domain purchase is not authorized. Use a Railway domain initially.

## Sources

- https://docs.railway.com/volumes/reference
- https://docs.railway.com/volumes
- https://docs.railway.com/builds/dockerfiles
- https://docs.railway.com/config-as-code/reference
