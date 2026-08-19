-- Drives the plugin in a headless Neovim: filetype, syntax groups, LSP
-- attach and restart. Run by TestPluginSmoke; the group positions are tied
-- to editors/vscode/examples/demo.sysml.
local repo = assert(arg and arg[1], 'usage: nvim -l smoke.lua <repo-root>')

-- prepend, as a plugin manager would
vim.opt.runtimepath:prepend(repo .. '/editors/nvim')

local failures = 0
local function check(name, ok, detail)
  if ok then
    print('ok   ' .. name)
  else
    failures = failures + 1
    print('FAIL ' .. name .. (detail and (' -- ' .. detail) or ''))
  end
end

-- plugin/ files on a runtimepath changed after startup are not auto-sourced
vim.cmd.source(repo .. '/editors/nvim/plugin/opensysml.lua')
vim.cmd.edit(repo .. '/editors/vscode/examples/demo.sysml')
local buf = vim.api.nvim_get_current_buf()

check('filetype', vim.bo[buf].filetype == 'sysml', vim.bo[buf].filetype)
check('syntax loaded', vim.b[buf].current_syntax == 'sysml', tostring(vim.b[buf].current_syntax))
check('commentstring', vim.bo[buf].commentstring == '// %s', vim.bo[buf].commentstring)

vim.cmd('syntax on')
vim.cmd('redraw')

local function group_at(lnum, col)
  return vim.fn.synIDattr(vim.fn.synID(lnum, col, true), 'name')
end
for _, g in ipairs({
  { 'line comment', 1, 1, 'sysmlLineComment' },
  { 'declaration keyword', 3, 1, 'sysmlDeclaration' },
  { 'declaration name', 3, 9, 'sysmlDeclarationName' },
  { 'modifier', 4, 5, 'sysmlModifier' },
  { 'relationship', 4, 13, 'sysmlRelationship' },
  { 'namespace', 4, 20, 'sysmlNamespace' },
  { 'block comment', 6, 7, 'sysmlBlockComment' },
  { 'typing operator', 7, 28, 'sysmlTypeOperator' },
  { 'type reference', 7, 31, 'sysmlTypeReference' },
  { 'number', 10, 41, 'sysmlNumber' },
  { 'constant', 11, 36, 'sysmlConstant' },
}) do
  local name, lnum, col, want = unpack(g)
  local got = group_at(lnum, col)
  check(name, got == want, got)
end

local server = repo .. '/bin/sysml-lsp'
if vim.fn.executable(server) == 1 then
  local attached = vim.wait(10000, function()
    return #vim.lsp.get_clients({ bufnr = buf, name = 'opensysml' }) > 0
  end, 100)
  check('lsp attached', attached)

  if attached then
    local client = vim.lsp.get_clients({ bufnr = buf, name = 'opensysml' })[1]
    check('lsp root', client.root_dir == repo, tostring(client.root_dir))
    check('lsp cmd', client.config.cmd[1] == server, tostring(client.config.cmd[1]))

    local responded = false
    local params = vim.lsp.util.make_position_params(0, client.offset_encoding)
    local handler = function(err)
      responded = err == nil
    end
    -- Client.request is a method from Neovim 0.11 and a closure before.
    if vim.fn.has('nvim-0.11') == 1 then
      client:request('textDocument/hover', params, handler, buf)
    else
      client.request('textDocument/hover', params, handler, buf)
    end
    vim.wait(10000, function()
      return responded
    end, 100)
    check('lsp hover', responded)

    local before = client.id
    vim.cmd('SysmlRestartServer')
    local restarted = vim.wait(10000, function()
      local now = vim.lsp.get_clients({ bufnr = buf, name = 'opensysml' })
      return #now > 0 and now[1].id ~= before
    end, 100)
    check('restart command', restarted)
  end
else
  print('skip lsp checks: ' .. server .. ' not built')
end

local ok, err = pcall(vim.cmd, 'checkhealth opensysml')
local text = ok and table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), '\n') or tostring(err)
check('checkhealth', ok and text:find('opensysml', 1, true) ~= nil, text)

if failures > 0 then
  print(failures .. ' failure(s)')
  vim.cmd('cquit')
end
print('all checks passed')
vim.cmd('quitall!')
