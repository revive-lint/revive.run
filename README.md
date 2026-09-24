# Website of Revive

Source of <https://revive.run/>, a [Hugo](https://gohugo.io/) site using the [Lotus Docs](https://github.com/colinwilson/lotusdocs) theme.

Most of the content is generated from [revive-lint/revive](https://github.com/revive-lint/revive):

- `/docs` from `README.md`
- `/r` from `RULES_DESCRIPTIONS.md`
- `/images` from `assets/`

The generated site is committed to `docs/` and served by GitHub Pages from the `master` branch.

## Publishing a release

The `REVIVE_VERSION` file holds the revive release the website was built from.
The [Build workflow](.github/workflows/build.yaml) runs twice a month (and on `workflow_dispatch`):
if `revive-lint/revive` has a newer release, it checks out that tag, regenerates the site, and commits `docs/` and `REVIVE_VERSION` to `master`.
Otherwise it does nothing.

To publish hand-written changes (landing page, guides, layouts, config) without a new revive release, run the workflow manually with the `force` input checked.

## Local development

Requires [Go](https://go.dev/) and [Hugo extended](https://gohugo.io/installation/).

1. Get the revive documentation:

```sh
git clone --depth 1 https://github.com/revive-lint/revive
```

2. Generate the content:

```sh
go run ./scripts/website
```

3. Start the local server:

```sh
hugo server -d /tmp/revive.run
```

4. Open <http://localhost:1313/> in your browser.

Always pass `-d` to `hugo server`: without it the dev build (with `localhost:1313` URLs) overwrites the committed `docs/`.

To rebuild `docs/` the way CI does:

```sh
hugo --gc --minify --cleanDestinationDir
```

## Layout

Hand-written:

- `hugo.yaml`, `go.mod`, `go.sum` — site and theme configuration
- `data/landing.yaml` — homepage
- `content/docs/rule.md`, `content/docs/formatter.md` — custom rule/formatter guides
- `content/docs/api.md` — sidebar entry redirecting to [pkg.go.dev](https://pkg.go.dev/github.com/revive-lint/revive)
  (the theme sidebar only links local pages, see `layouts/_default/redirect.html`)
- `layouts/partials/head.html`, `layouts/partials/docs/head.html` — theme head partials without Lotus Docs' hardcoded description/author metadata
- `assets/images/logos/` — logo used in the header
- `static/` — favicons, `site.webmanifest` and `sw.js`
- `scripts/website/` — content generator
- `REVIVE_VERSION` — revive release the published site is built from, updated by the Build workflow

Generated (git-ignored):

- `content/docs/documentation.md` (`/docs/`), `content/docs/rules.md` (`/r/`), `static/images/`

Published (committed): `docs/` — marked `linguist-generated` in `.gitattributes`, so GitHub collapses it in PR diffs.

## Notes

- `static/sw.js` is a self-destroying service worker: the previous Gatsby site registered one,
  and without this file returning visitors would keep seeing the cached old site. Keep it deployed.
- Links of the form `/r#<rule>` and `/r/#<rule>` (the ones printed by revive) keep working:
  all rules live on the single `/r/` page and its heading anchors use GitHub-style slugs.
