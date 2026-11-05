# Source intake

The worker resolves a public GitHub repository to its default branch commit SHA, downloads the ZIP for that commit, or uses the uploaded ZIP. The source is unpacked into a private work directory. Uploaded code is never run.

Limits for this stage are 20 MiB compressed, 100 MiB declared expanded size, 2,000 archive entries, 2 MiB per extracted text file, 512 bytes per path, and 20 path segments. The extractor rejects absolute or parent paths, backslashes, duplicate paths, symlinks, special files, and unreadable ZIP entries. Nested archives are listed with `nested_archive` and never unpacked.

Each file has an explainable reason code. Dependencies, build output, generated files, binary files, likely secrets, nested archives, and oversized files are excluded. Lock files and repository metadata are reported as `tooling_metadata` and excluded. Build and deployment files are eligible as `build_config` but are not selected by default. Test files remain eligible. They are selected by default for `technical-overview` and left unselected for `academic`. Manifests such as `go.mod` and `package.json` stay selected because evidence extraction derives structured facts from them. At most 500 files are selected by default.

The worker confirms the default selection automatically after saving the inventory, so intake advances to evidence review without a manual step. The file page remains readable with every reason code. If automatic confirmation fails, the job stays at `needs_input` / `selecting_files` and the user can confirm or adjust 1 to 500 eligible files manually. If every file is excluded, the inventory remains visible and asks for a different project input.

The inventory stores the source digest, optional GitHub commit SHA, private workspace path, file paths, sizes, content digests for eligible files, reason codes, and selected flags. Confirmation changes the job from `needs_input` to `input_ready` with stage `source_selected`. Evidence extraction begins in milestone 4.

Set `GITHUB_TOKEN` on the worker to use an optional server-side token for public GitHub API requests. The token is never sent to the browser or to the archive redirect host. GitHub rate limits and permanent source errors are shown promptly; transient failures use bounded retries.
