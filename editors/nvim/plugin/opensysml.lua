if vim.g.loaded_opensysml then
  return
end
vim.g.loaded_opensysml = 1

-- Neovim 0.12 detects these itself; older versions need the mapping.
vim.filetype.add({
  extension = {
    sysml = 'sysml',
    kerml = 'kerml',
  },
})

-- Attach via autocmd, not ftplugin/: Neovim 0.12's bundled sysml/kerml
-- ftplugins can win the b:did_ftplugin race.
vim.api.nvim_create_autocmd('FileType', {
  pattern = { 'sysml', 'kerml' },
  group = vim.api.nvim_create_augroup('opensysml', {}),
  callback = function(event)
    require('opensysml').start(event.buf)
  end,
})

vim.api.nvim_create_user_command('SysmlRestartServer', function()
  require('opensysml').restart()
end, { desc = 'Restart the OpenSysML language server' })
