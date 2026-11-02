# AI Doc Gen

AI Doc Gen is a planned web application that turns a GitHub repository or a source archive into editable Typst and a compiled PDF. Users choose a document profile and may add a short project description.

The [product and niche analysis](analysis/README.md) records evidence, assumptions, and product decisions. The [development guide](development/README.md) contains the incremental implementation plan. The first implemented step is a local Go website with a disabled input preview; document submission and generation are not enabled yet.

Run the site with `go run ./cmd/ai-doc-gen`, then open `http://localhost:8080`. PostgreSQL for future job processing is available through `docker compose up -d postgres`. The [backend package map](development/structure.md), [Go style](development/code_style.md), and [job lifecycle](development/job_lifecycle.md) guide the next steps.
