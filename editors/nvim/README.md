# SysML v2 for Neovim (OpenSysML)

Syntax highlighting and language support for `.sysml` and `.kerml` files,
backed by OpenSysML's `sysml-lsp` server: diagnostics, hover, go-to-definition,
document symbols and typed completion.

Requires Neovim 0.10+ and has no plugin dependencies. Like the VS Code
extension, the plugin is loaded from a repository checkout and deliberately
**not published** to a plugin registry.

## Install

Build the server first, from the repo root:

```bash
make build          # builds bin/sysml-lsp
```

Then add `editors/nvim` to your runtimepath with your plugin manager:

```lua
-- lazy.nvim, from a local checkout
{ dir = '/path/to/OpenSysML/editors/nvim' }
```

```vim
" vim-plug, cloning the repository and loading the subdirectory
Plug 'Open-MBEE/OpenSysML', { 'rtp': 'editors/nvim' }
```

Or link it as a native package:

```bash
mkdir -p ~/.local/share/nvim/site/pack/opensysml/start
ln -s /path/to/OpenSysML/editors/nvim ~/.local/share/nvim/site/pack/opensysml/start/opensysml
```

Then open any `.sysml` file. The plugin finds the server in this order:

1. `server.path`, if set;
2. `bin/sysml-lsp` under the buffer's project root (a repo checkout that ran
   `make build`);
3. `sysml-lsp` on `$PATH`.

If none exist, highlighting still works and a one-time warning explains how to
build the server. `:SysmlRestartServer` restarts it after a rebuild, and
`:checkhealth opensysml` shows what the plugin resolved.
`editors/vscode/examples/demo.sysml` is a highlighting smoke-test file.

## Settings

Settings mirror the VS Code extension's and are all optional:

```lua
require('opensysml').setup({
  server = {
    path = '',      -- absolute path to sysml-lsp
    args = {},      -- extra server arguments
    enabled = true, -- false keeps syntax highlighting only
  },
})
```

| Setting | Default | Meaning |
| --- | --- | --- |
| `server.path` | `''` | Absolute path to `sysml-lsp`; empty falls back to the project root's `bin/`, then `$PATH`. |
| `server.args` | `{}` | Extra server arguments. |
| `server.enabled` | `true` | Set to `false` for highlighting without a server. |

Call `setup` before opening SysML buffers (with lazy.nvim, `opts = {}` does).
Changed settings apply to servers started afterwards; `:SysmlRestartServer`
applies them to running ones.

## Syntax generation

`syntax/*.vim` are generated — do not edit them by hand. The keyword grouping
is shared with the VS Code extension (`editors/internal/keywords`) and the
keyword list comes from `internal/core/lexer.Keywords()`, so highlighting
cannot drift from the lexer:

```bash
make nvim-syntax       # regenerate
go test ./editors/...  # fails if the committed files are stale
```
