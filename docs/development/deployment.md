# Railway operations

## Current deployment

Public URL: https://ai-doc-gen-production.up.railway.app

Verified 2026-11-06: cloud build and readiness passed; a public repository job completed and its PDF downloaded through the owner browser. Non-owner downloads return 404. This was a keyless infrastructure check, not a live-model quality evaluation. The service tracks `codex/railway-deploy`; GitHub pushes to that branch trigger deployment. The account is on Limited Trial; no paid subscription was activated. Keep local-only jobs and provider credentials local.

## Topology

Configure the service in the Railway dashboard: one Docker service runs separate Go web and worker processes under tini, with a single persistent volume at `/data`. PostgreSQL runs as a separate private Railway service. Use one application replica. The existing `Dockerfile` and Docker Compose remain the local development path. `deploy/Railway.Dockerfile` pins the same Go, Alpine and Typst versions; update both images together.

The attached Streamlit tutorial supplies platform concepts only. No Python demo or Streamlit authentication has been imported. Google login is deferred by user choice. Documents remain owned by their browser cookie: losing that cookie loses access. This initial deployment is not the final security/testing milestone.

## Service setup

1. Create a Railway project and PostgreSQL service. Do not expose PostgreSQL publicly.
2. Connect `vasyukov1/ai-doc-gen` at the reviewed deployment revision. Root is `/`. Select Dockerfile builder and set Dockerfile Path to `deploy/Railway.Dockerfile`. Set healthcheck to `/readyz`, restart policy to On Failure with 10 retries, and leave the start-command override empty. New Railway services cannot opt into the deprecated config-as-code format; no railway.toml is required.
3. Attach a persistent volume to the application at `/data`. Keep one replica and disable sleeping/serverless behavior because the worker must run without HTTP traffic.
4. Set `DATABASE_URL` to Railway's reference to the PostgreSQL private connection URL, normally `${{Postgres.DATABASE_URL}}` (use the actual service name).
5. Set `DATA_DIR=/data`, `COOKIE_SECURE=true`, `RETENTION_DAYS=7`. Leave `UPLOAD_DIR`, `SOURCE_DIR` and `STATE_DIR` unset: the launcher rejects overrides to keep key initialization and volume ownership consistent. Railway supplies `PORT`. The launcher maps it to `HTTP_ADDR` unless explicitly set.
6. The launcher initializes directory ownership and drops to UID 10001. It runs migrations once before starting children, then generates a master encryption key on the persistent disk before either child can use it. Alternatively supply one stable base64 32-byte `LLM_ENCRYPTION_KEY` through Railway's secret variables. Never commit or log its value. Back up the generated key with the volume; changing or losing it prevents decryption of stored provider connections.
7. Generate a Railway HTTPS domain targeting the assigned port. Validate `/healthz`, `/readyz`, and a keyless document job. A PDF generated without a provider is an outline, not validated model prose.

No server provider key or local model is required. Users may supply their own supported connection through the application. Do not transfer local jobs, cookies, provider keys or private source archives to bootstrap deployment. Docker ignores local artifacts, scratch directories and secret files.

## Failure and restart behavior

Migration failure prevents startup. Either child exiting causes the sibling to terminate and the container to exit with failure, allowing Railway to restart it. SIGTERM is forwarded to both children. Database leases retain job recovery semantics; persisted uploads, sources, renders and key survive container replacement. Configure the platform's termination grace to allow at least 15 seconds. A deployment with a volume may have downtime; this topology is not horizontally scalable.

`/readyz` checks database connectivity; it does not prove model availability. Worker death is handled by container supervision. Inspect both web and worker logs. Never log credentials, source bodies or database connection strings.

## Backup and rollback

Enable PostgreSQL and application volume backups in Railway, subject to the selected plan and cost. Save matching snapshots; database rows reference files and encrypted credentials. Restore into a separate environment and verify owner access and downloads before cutover. The initial deployment does not implement backup automation.

Roll back the application to a known compatible commit/image through Railway. Do not automatically run down migrations: newer database schemas may not be backward compatible. Preserve both volumes when deleting/recreating a service. Record deployed commit, environment and backup identifiers without secrets.

## Domain and Google login

Use the generated HTTPS domain first. For an owned domain, add it to Railway and create exactly the DNS records shown by the current UI; verify HTTPS before using it. Do not switch domains casually: browser ownership cookies do not move between domains. Buying a domain requires separate user selection and budget.

Google sign-in is a later integration, not enabled by the tutorial's environment variables. It needs verified OIDC login, state/nonce handling, secure sessions, a persisted user identity and an explicit migration/linking flow for existing browser-owned jobs. Do not replace ownership checks with a login screen alone.

## Platform evidence

Verified 2026-11-06 using official documentation:
- [Volumes and runtime ownership](https://docs.railway.com/volumes)
- [Volume limits and deployment downtime](https://docs.railway.com/volumes/reference)
- [Dockerfile deployment](https://docs.railway.com/builds/dockerfiles)
- [Configuration as code](https://docs.railway.com/config-as-code/reference)

Check the account's actual resource limits and billing before provisioning. The documented Free plan permits only one volume per project; this topology needs two (application and PostgreSQL). Trial or paid capacity may be necessary; no recurring purchase is implicitly authorized.
