<p align="center"><img src="https://raw.githubusercontent.com/go-ruby-format/brand/main/social/go-ruby-format.png" alt="go-ruby-format/docs" width="720"></p>

# go-ruby-format/docs

Versioned documentation for [go-ruby-format](https://github.com/go-ruby-format),
built with [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) and
versioned with [mike](https://github.com/jimporter/mike). Published to the
`gh-pages` branch and served at <https://go-ruby-format.github.io/docs/>.

The organization landing page ([go-ruby-format.github.io](https://go-ruby-format.github.io))
links here.

## Local preview

```bash
python -m venv .venv && . .venv/bin/activate
pip install -r requirements.txt
mkdocs serve                       # http://localhost:8000 (current sources)
mike serve                         # preview the versioned site
```

## Releasing a new docs version

```bash
mike deploy --push --update-aliases <version> latest
mike set-default --push latest
```
