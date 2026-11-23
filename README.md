# matteosalvatoreg

Personal CV site published at `matteosalvatore.dev` (see `CNAME`), now also served by a small Go API. The page is for **Matthew Giancola**; there is no "Alex Chen" placeholder data left in the code.

Loopback and container-publishing rules for all six demos: [`../docs/LOCAL-DEMOS.md`](../docs/LOCAL-DEMOS.md). Runnable scripts and a loopback-only Compose file: [`../examples`](../examples/README.md).

## What is here

| Path | Purpose |
| --- | --- |
| `index.html` | The single-page CV with inline CSS and JavaScript. |
| `img.png` | Profile image. |
| `CNAME` | GitHub Pages custom domain: `matteosalvatore.dev`. |
| `main.go` | Gin server that embeds and serves `index.html` and `img.png`, and exposes a JSON API. |
| `main_test.go` | Router-level tests for the site and API. |
| `go.mod`, `go.sum` | Pinned dependencies (Gin). |

### The CV page

- Sections: summary, experience, education, skills, credentials and professional status, projects, and contact.
- The summary's "years of experience" field is calculated: the original `15+ years` claim in 2026 gives a start year of **2011**, and the current count is derived from that start year on page load. The JSON API uses the same calculation.
- EN / FR / IT language switcher driven by `data-i18n` attributes.
- Buttons that build a one-page PDF résumé from the DOM with jsPDF.
- Links to GitHub, LinkedIn, HackerRank, Salesforce Trailblazer, HackerOne, and a booking calendar.
- A project accordion linking VENI, VICI, VIDI, VINI, and this site.
- Third-party assets from Bootstrap 4.4.1, jQuery 3.4.1, Popper, Font Awesome, and jsPDF CDNs.

### The Go server

```bash
go test ./...
go run .                 # http://127.0.0.1:8081
```

| Variable | Default | Meaning |
| --- | --- | --- |
| `CV_HOST` | `127.0.0.1` | Bind address. |
| `CV_PORT` | `8081` | Bind port. |

| Endpoint | Purpose |
| --- | --- |
| `GET /` | The CV page, embedded in the binary. |
| `GET /img.png` | Profile image, embedded in the binary. |
| `GET /health` | Service status. |
| `GET /api/cv` | Real CV data: name, contact, summary, experience, education, skills, credentials, links. |
| `GET /api/about` | Site metadata: domain, languages, static-hosting note. |

The `/api/cv` response is built from the same content as `index.html`. The server sends `nosniff`, frame-deny, and referrer headers. Tests call the real Gin router rather than reimplementing handlers inline.

## Tests

`go test ./...` passes. There are no JavaScript tests for the page's translation or PDF code.

## Remaining gaps

- CDN dependencies have no Subresource Integrity hashes and require network access; there is no offline fallback.
- The Go API and the page share content manually. A future change could let the page fetch `/api/cv` so there is one runtime source of truth.
- There is no Dockerfile or deployment manifest; the site is designed for static hosting via GitHub Pages.
- No favicon, privacy policy, or analytics configuration is included.
- The old duplicate `matteosalvatoreg-master.zip` was removed; the working files are the source of truth.
