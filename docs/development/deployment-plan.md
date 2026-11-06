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
- [x] Implement DEP1–DEP3 and deployment instructions.
- [x] Run go build and isolated Docker smoke checks, without automated test files.
- [x] Review changes, commit and push task-owned paths.
- [x] Configure Railway PostgreSQL and persistent application service; check actual plan limits before chargeable actions.
- [x] Deploy, validate HTTPS and keyless document flow, record URL and limitations.

## Review focus

Migration failure must prevent serving. Either child exiting must stop its sibling. Termination must reach both children. Mounted volumes begin owned by root. Private local artifacts must never enter the image. No local database or provider secret may be uploaded. Initial keyless output is an outline, not proof of real-model quality.

## Resume point

Deployment and local/cloud checks complete. Public URL: https://ai-doc-gen-production.up.railway.app. Next read-only check: `curl -fsS https://ai-doc-gen-production.up.railway.app/readyz`. Do not purchase a subscription/domain or start Google login without the corresponding request.

## Sources

- https://docs.railway.com/volumes/reference
- https://docs.railway.com/volumes
- https://docs.railway.com/builds/dockerfiles
- https://docs.railway.com/config-as-code/reference

## Manual evidence (local, 2026-11-06)

- `go build ./...`: exit 0.
- Dedicated Docker image built successfully; hosted cache mounts removed after independent review.
- Fresh isolated PostgreSQL and application volume: migrations reached version 11; `/readyz` returned `ready`; process listing showed separate app-owned web and worker children.
- Invalid submission set HttpOnly, Secure, SameSite=Strict owner cookie.
- Public `hse-coursework-docs` keyless run reached `completed/rendered`, 4/4 baseline sections, 0 model sections, explicit outline outcome. Downloaded PDF has 9 pages.
- Container restart preserved owner access and PDF download. Key file is app-owned mode 600.
- Sending SIGTERM to the worker stopped the container with exit 1; failed database startup also exited 1 without serving.
- Independent review found incompatible hosted Docker cache mounts and STATE_DIR key-path inconsistency. Removed cache mounts and reject directory overrides. Final image rebuild exited 0; it served readiness and the identical retained PDF. SIGTERM stopped it with exit 0. Custom STATE_DIR was rejected.
- First two manual submission attempts used incorrect metadata field names and were rejected with an actionable 400; corrected to metadata.project.

Cloud deployment and remote checks completed. The connected service uses dashboard settings because new services cannot enable deprecated config-as-code. Local evidence does not establish a published site.

## Cloud checkpoint (2026-11-06)

- Project `b16556ec-d8ef-4f39-b6f4-5aae99476896`, environment `ad0c3524-ff9f-4e57-a8a6-92bf7f7b39a2`.
- App service `888df803-b588-45ed-915e-665b44e76789`; PostgreSQL service `f94dbfb0-a185-430e-aa3b-6e70bd9567bc`. Both observed Online.
- Deployment `a6eaa7b6-5e5b-425a-9712-5a344dc9a3a5` built revision `2a8cd10`; AMD64 Typst checksum passed and `/readyz` healthcheck succeeded at 07:23:55 Europe/Madrid.
- Dashboard configuration: codex/railway-deploy branch, Dockerfile builder, deploy/Railway.Dockerfile, /readyz, one replica, no serverless, On Failure/10 retries.
- Variables: private Postgres DATABASE_URL reference, DATA_DIR=/data, COOKIE_SECURE=true, RETENTION_DAYS=7, PORT=8080. No provider keys supplied.
- Separate application and PostgreSQL volumes; database has no public TCP access.
- Account displays Limited Trial, 30 days or $5 credit. No subscription purchased. External repository/provider connectivity remains unverified.
- User explicitly confirmed public access; https://ai-doc-gen-production.up.railway.app is live on port 8080.
- Browser job `423ab957-bbad-4c55-8433-c004ff0183a0` fetched the public hse-coursework-docs repository and completed rendering with 4 baseline sections, no provider calls, and an explicit outline label. PDF downloaded through the browser successfully. This confirms GitHub access but not LLM-provider connectivity on Limited Trial.
- Independent HTTP checks: /healthz -> ok, /readyz -> ready, unowned PDF request -> 404. The browser-owned PDF download is a valid PDF.
- Google login, custom purchased domain, automated backups and live-provider verification on the deployed host remain deferred. No subscription purchased.
