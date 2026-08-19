-- LSP client for the OpenSysML language server (sysml-lsp).
local M = {}

local defaults = {
  server = {
    path = '', -- absolute path to sysml-lsp; '' falls back to bin/, then $PATH
    args = {}, -- extra arguments passed to the language server
    enabled = true, -- start the server; false keeps syntax highlighting only
  },
}

M.options = vim.deepcopy(defaults)

function M.setup(opts)
  M.options = vim.tbl_deep_extend('force', vim.deepcopy(defaults), opts or {})
end

local executable = vim.fn.has('win32') == 1 and 'sysml-lsp.exe' or 'sysml-lsp'

-- resolve prefers the configured path, then the project root's bin/, then
-- $PATH.
function M.resolve(root)
  local configured = M.options.server.path
  if configured ~= '' then
    return vim.fn.executable(configured) == 1 and configured or nil
  end
  if root then
    local candidate = vim.fs.joinpath(root, 'bin', executable)
    if vim.fn.executable(candidate) == 1 then
      return candidate
    end
  end
  local found = vim.fn.exepath(executable)
  return found ~= '' and found or nil
end

local warned = false

-- start attaches bufnr to the server for its project root; a missing server
-- warns once per session.
function M.start(bufnr)
  if not M.options.server.enabled then
    return
  end
  bufnr = bufnr or vim.api.nvim_get_current_buf()
  local root = vim.fs.root(bufnr, '.git')
  local command = M.resolve(root)
  if not command then
    if not warned then
      warned = true
      vim.notify(
        ('OpenSysML: could not find %s. Build it with `make build` and set '
          .. 'require("opensysml").setup{ server = { path = ... } }, or put it '
          .. 'on your PATH. Syntax highlighting still works.'):format(executable),
        vim.log.levels.WARN
      )
    end
    return
  end

  local cmd = { command }
  vim.list_extend(cmd, M.options.server.args)
  vim.lsp.start({
    name = 'opensysml',
    cmd = cmd,
    root_dir = root,
  }, { bufnr = bufnr })
end

local function attach_loaded_buffers()
  for _, bufnr in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_loaded(bufnr) then
      local ft = vim.bo[bufnr].filetype
      if ft == 'sysml' or ft == 'kerml' then
        M.start(bufnr)
      end
    end
  end
end

-- restart stops every client and reattaches buffers once it exits, picking up
-- a rebuilt binary or changed setup.
function M.restart()
  warned = false
  local clients = vim.lsp.get_clients({ name = 'opensysml' })
  if #clients == 0 then
    attach_loaded_buffers()
    return
  end
  for _, client in ipairs(clients) do
    -- Client.stop is a method from Neovim 0.11 and a closure before.
    if vim.fn.has('nvim-0.11') == 1 then
      client:stop()
    else
      client.stop()
    end
  end

  local tries = 0
  local function poll()
    if #vim.lsp.get_clients({ name = 'opensysml' }) == 0 then
      attach_loaded_buffers()
      return
    end
    tries = tries + 1
    if tries >= 40 then
      vim.notify(
        'OpenSysML: the language server did not stop; run :SysmlRestartServer again.',
        vim.log.levels.ERROR
      )
      return
    end
    vim.defer_fn(poll, 50)
  end
  vim.defer_fn(poll, 50)
end

return M
