# AI Doc Gen

AI Doc Gen is a planned web application that turns a GitHub repository or a source archive into editable Typst and a compiled PDF. Users choose a document profile and may add a short project description.

The [product and niche analysis](docs/analysis/README.md) records evidence, assumptions, and product decisions. The [development guide](docs/development/README.md) contains the incremental implementation plan. The service accepts a public GitHub URL or ZIP, inspects files in a separate worker, and lets users confirm the file selection. Evidence extraction and document generation are not enabled yet.

Run `make run` to start PostgreSQL, migrations, the website, and the worker. Open `http://localhost:8080`. Run `make down` to stop the stack without deleting its data. `make build` builds the binary, `make test` runs the Go test suite when tests are added, and `make mocks` reports that mock generation is not configured yet. Use `make db-generate` after changing migrations to regenerate the checked-in Jet schema. The [development guide](docs/development/README.md) has local commands and architecture details.
