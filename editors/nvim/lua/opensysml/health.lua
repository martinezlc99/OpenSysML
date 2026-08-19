local M = {}

function M.check()
  local health = vim.health
  local opensysml = require('opensysml')

  health.start('opensysml')

  if vim.fn.has('nvim-0.10') == 1 then
    health.ok('Neovim >= 0.10')
  else
    health.error('Neovim >= 0.10 is required for the language client')
  end

  if not opensysml.options.server.enabled then
    health.info('server.enabled is false; syntax highlighting only')
    return
  end

  local root = vim.fs.root(vim.fn.getcwd(), '.git')
  local command = opensysml.resolve(root)
  if command then
    health.ok(('language server: %s'):format(command))
  else
    health.warn(
      'sysml-lsp not found (checked server.path, <project root>/bin, $PATH); '
        .. 'build it with `make build` in an OpenSysML checkout or point '
        .. 'server.path at the binary'
    )
  end

  local clients = vim.lsp.get_clients({ name = 'opensysml' })
  if #clients == 0 then
    health.info('no running language server clients')
  else
    for _, client in ipairs(clients) do
      health.ok(('client %d running for %s'):format(client.id, client.root_dir or '(single file)'))
    end
  end
end

return M
