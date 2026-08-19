# Neovim extension design

A Neovim plugin at `editors/nvim` with feature parity to `editors/vscode`:
syntax highlighting for `.sysml`/`.kerml` plus an LSP client for `sysml-lsp`.
Like the VS Code extension, it is side-loaded from a repository checkout, not
published to a plugin registry.

## Goals

- Filetype detection, syntax highlighting, and language-server support for
  SysML v2 and KerML in Neovim ≥ 0.10, with no plugin dependencies.
- Highlighting that cannot drift from the lexer: syntax files are generated
  from `internal/core/lexer.Keywords()` and gated by a staleness test, the
  same contract the VS Code TextMate grammars follow.
- One keyword grouping shared by both editors. The groups currently private
  to `editors/vscode/tools/gengrammar` move to `editors/internal/keywords`;
  both generators consume it, so a keyword can never be coloured as a
  declaration in VS Code and a modifier in Neovim.

## Non-goals

- A tree-sitter grammar. Regex syntax files generated from the lexer cover
  highlighting today; tree-sitter is a separate, much larger project.
- Publishing to a plugin manager registry (mirrors the VS Code stance).

## Layout

```
editors/internal/keywords/     keyword groups shared by both editor generators
  keywords.go                  Group type, Groups(), Partition() validation
  keywords_test.go
editors/nvim/
  README.md                    install, settings, generation workflow
  plugin/opensysml.lua         filetype registration, :SysmlRestartServer
  ftplugin/{sysml,kerml}.lua   comment options + LSP attach per buffer
  syntax/{sysml,kerml}.vim     GENERATED — do not edit by hand
  lua/opensysml/init.lua       setup(), server resolution, start/restart
  lua/opensysml/health.lua     :checkhealth opensysml
  tools/gensyntax/             generator + drift-gate tests (mirrors gengrammar)
```

`editors/vscode/tools/gengrammar` is refactored to consume
`editors/internal/keywords`; its generated output stays byte-identical, which
`TestCommittedGrammarsAreCurrent` proves.

## Shared keyword groups

`keywords.Partition()` returns the ordered groups (declaration, control,
modifier, relationship, operator, constant) plus the lexer keywords no group
claims, and errors if a group names a keyword the lexer no longer knows or two
groups claim the same keyword. Each generator maps group names to its own
colour vocabulary (TextMate scopes / Vim highlight groups) and fails
generation on an unmapped group, so adding a group forces both editors to
choose a colour. Unclaimed keywords still highlight via a catch-all group, so
a new lexer keyword needs no generator change.

## Syntax files

`gensyntax` renders one `syntax/<ft>.vim` per language (same rules, per-
filetype group prefix, mirroring the per-scope tmLanguage files). Rules match
the TextMate grammar: keyword groups, namespace segments before `::`, numbers,
operators, the name after `def`/`package`/`namespace`/`alias`/`library`, the
type after `:`/`:>`/`::>`/`:>>`, strings, quoted names, and comments. Vim
resolves same-position overlaps last-defined-wins (TextMate is first-wins), so
the generated file emits rules in reverse priority order with comments and
quoted regions last. Drift gates mirror gengrammar's:
`TestCommittedSyntaxIsCurrent` and `TestEveryKeywordIsHighlighted`, run by the
existing `go test ./editors/...`. `make nvim-syntax` regenerates.

## LSP client

Pure Lua on `vim.lsp.start` (no nvim-lspconfig dependency). The ftplugin
attaches each `sysml`/`kerml` buffer; clients are shared per project root
(`vim.fs.root(buf, '.git')`). Server resolution matches VS Code: configured
path → `<root>/bin/sysml-lsp` → `$PATH`; if none is found a one-time warning
explains `make build`, and highlighting keeps working. Settings mirror the VS
Code ones via `require('opensysml').setup{ server = { path, args, enabled } }`.
`:SysmlRestartServer` stops every client and reattaches the buffers they
served (polling for exit, as nvim-lspconfig's restart does), picking up a
rebuilt binary. `:checkhealth opensysml` reports the resolved server and
running clients.

## Verification

- `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...` — includes
  the new drift gates and proves the gengrammar refactor changed no output.
- Headless smoke test against the local `nvim`: filetype detection, syntax
  groups present, LSP attaches to `examples/` files using `bin/sysml-lsp`.
