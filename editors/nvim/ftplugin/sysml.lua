-- Comment options for Neovim < 0.12, which lacks the bundled sysml ftplugin.
if vim.b.did_ftplugin then
  return
end
vim.b.did_ftplugin = 1

vim.bo.commentstring = '// %s'
vim.bo.comments = 's1:/*,mb:*,ex:*/,://'
vim.b.undo_ftplugin = 'setlocal commentstring< comments<'
