# Browser interface

Server-rendered HTML with a small TypeScript layer. Go renders the pages; Vite
builds the stylesheet and script that the binary embeds.

## Layout

- `templates/` — Go page templates and shared partials.
- `src/styles/` — design tokens, base styles, components, page styles.
- `src/*.ts` — theme selection, file inputs, auto-growing textareas.
- `static/` — build output. It is committed because `go:embed` reads it at
  compile time, so `go build ./...` works without Node installed.
- `pages.go` — embeds `templates/` and `static/`, parses pages, and exposes the
  small template helpers used by the views.

## Building

```sh
make ui          # install and build the assets
npm run build    # build only
npm run typecheck
npm run dev      # Vite dev server proxied to localhost:8080
```

Run `make ui` after changing anything in `src/` and commit the refreshed
`static/app.css` and `static/app.js` together with the source change.

## Design notes

The interface follows a warm, editorial direction: paper-like backgrounds,
serif headings, thin rules, and one quiet ochre accent. Light and dark themes
share every token; the system preference is the default and the header menu
stores an explicit choice in `localStorage`.

Pages work without JavaScript. Scripts add theme selection, file-name previews,
auto-growing textareas, submit feedback, and live status updates.

While a job runs, the status sidebar draws a small canvas scene: scraps of paper
and letters circle a sheet and settle into lines as sections complete. It is
decorative, hidden from assistive technology, static under
`prefers-reduced-motion`, and never captures pointer input. With JavaScript the
status page polls `/api/jobs/{id}` and updates in place instead of reloading;
without it, the page keeps the slower `<noscript>` refresh.
