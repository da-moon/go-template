function! after#coc#bootstrap() abort
  call SpaceVim#logger#info("[ after#coc ] bootstrap function called")
  for extension in g:coc_extensions
    call SpaceVim#logger#info("[ after#coc ] installing [ " . extension . " ] coc extension")
    call coc#add_extension(extension)
  endfor
    call SpaceVim#logger#info("[ after#coc ] setting coc preference")
  call coc#config('coc.preferences', g:coc_preferences)
endfunction
