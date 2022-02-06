let g:spacevim_custom_plugins = [
  \ ['fatih/vim-go',                            { 'on_ft' : 'go'}],
  \ ['jvirtanen/vim-hcl',                       { 'on_ft' : 'hcl'}],
  \ ['vmchale/just-vim'],
  \ ['PhilRunninger/nerdtree-visual-selection'],
  \ ['mg979/vim-visual-multi'],
  \ ['hashivim/vim-hashicorp-tools'],
  \ ['will133/vim-dirdiff'],
  \ ['tarekbecker/vim-yaml-formatter'],
  \ ['machakann/vim-highlightedyank'],
  \ ['sainnhe/sonokai'],
  \ ['dracula/vim'],
  \ ['sheerun/vim-polyglot'],
  \ ['phanviet/vim-monokai-pro'],
  \ ['cormacrelf/vim-colors-github'],
  \ ]
function! init#before() abort
  let g:spacevim_enable_debug = 1
  let g:spacevim_realtime_leader_guide = 1
  let g:spacevim_enable_tabline_filetype_icon = 1
  let g:spacevim_enable_statusline_display_mode = 0
  let g:spacevim_enable_os_fileformat_icon = 1
  let g:spacevim_buffer_index_type = 1
  call before#cfg#bootstrap()
  call before#coc#common#bootstrap()
  call before#coc#list#bootstrap()
  call before#coc#json#bootstrap()
  call before#themes#bootstrap()
  if has('nvim')
    call before#nvim#bootstrap()
  endif
endfunction

function! init#after() abort
  call after#coc#bootstrap()
  let g:neomake_open_list=0
  set showcmd
  nnoremap <silent> [Window]a :cclose<CR>:lclose<CR>
  let g:spacevim_todo_labels = [
    \'FIXME',
    \'[FIXME]',
    \'[ FIXME ]',
    \'QUESTION',
    \'[QUESTION]',
    \'[ QUESTION ]',
    \'TODO',
    \'[TODO]',
    \'[ TODO ]',
    \'IDEA',
    \'[IDEA]',
    \'[ IDEA ]',
    \]
endfunction
