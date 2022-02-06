function! before#nvim#bootstrap() abort
  call SpaceVim#logger#info("[ before#nvim ] bootstrap function called")
  call SpaceVim#logger#info("[ before#nvim ] display result of incremental commands (ex. :%s/pat1/pat2/g) ")
  set inccommand=split
  call SpaceVim#logger#info("[ before#nvim ] enter terminal buffer in Insert instead of Normal mode")
  autocmd TermOpen term://* startinsert
endfunction
