# syntax = docker/dockerfile-upstream:master-labs
#-*-mode:dockerfile;indent-tabs-mode:nil;tab-width:2;coding:utf-8-*-
# vi: ft=dockerfile tabstop=2 shiftwidth=2 softtabstop=2 expandtab:

# ─── PREBUILD IMAGE ALIASES ─────────────────────────────────────────────────────
# hadolint ignore=DL3007
FROM fjolsvin/just:latest as just
# hadolint ignore=DL3007
FROM fjolsvin/convco:latest as convco
# hadolint ignore=DL3007
FROM fjolsvin/clog:latest as clog
# ────────────────────────────────────────────────────────────────────────────────
FROM gitpod/workspace-full:build-branch-master
USER root
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3009
RUN \
  yes | add-apt-repository ppa:neovim-ppa/unstable \
  && apt-get update \
  && curl https://i.jpillora.com/chisel! | bash \
  && curl -fsSL https://starship.rs/install.sh | bash -s -- --force
ARG APT_PACKAGES="\
  dropbear \
  cmake \
  xfonts-utils \
  neovim \
  jq \
  curl \
  libssl-dev \
  fonts-symbola \
  fzf \
  "
RUN install-packages ${APT_PACKAGES}
ARG PYTHON_PACKAGES="\
  pynvim \
  pre-commit \
  pylint \
  yapf \
  autoflake \
  isort \
  coverage \
  "
ENV PIP_USER="false"
# hadolint ignore=DL3013
RUN \
  python3 -m pip install --no-cache-dir ${PYTHON_PACKAGES}
ARG NODEJS_PACKAGES="\
  remark \
  remark-cli \
  remark-stringify \
  remark-frontmatter \
  wcwidth \
  prettier \
  bash-language-server \
  dockerfile-language-server-nodejs \
  "
# hadolint ignore=DL3059
RUN \
  yarn global add --prefix /usr/local ${NODEJS_PACKAGES}
# hadolint ignore=DL3059
RUN \
  curl -sLf https://spacevim.org/install.sh | bash \
  && find "${HOME}" \
  -not -group "$(id -g gitpod)" \
  -not -user "$(id -u gitpod)" \
  -print0 \
  | xargs -0 -r -P "$(nproc)" \
  chown --no-dereference "$(id -u gitpod):$(id -g gitpod)" ;
ENV EDITOR="nvim"
ENV VISUAL="nvim"
ENV PIP_USER="false"
# hadolint ignore=SC2016
RUN \
  sed -i \
  -e '/^\s*#/d' \
  -e '/gitpod/d' \
  -e '/PIP_USER/d' \
  -e '/nvim/d' \
  -e '/local/d' \
  -e '/starship/d' \
  -e '/just/d' \
  -e '/^\s*$/d' \
  /etc/bash.bashrc \
  && ( \
  echo "export EDITOR=\"${EDITOR}\"" ; \
  echo "export VISUAL=\"${VISUAL}\"" ; \
  echo "export PIP_USER=\"${PIP_USER}\"" ; \
  echo 'export PATH="${PATH}:/usr/local/bin"' ; \
  echo "[ -d ${HOME}/.cargo/bin ] && export PATH=\"\${PATH}:${HOME}/.cargo/bin\";"; \
  echo "[ -d ${HOME}/.pyenv/bin ] && export PATH=\"\${PATH}:${HOME}/.pyenv/bin\";"; \
  echo "[ -d ${HOME}/.pyenv/shims ] && export PATH=\"\${PATH}:${HOME}/.pyenv/shims\";"; \
  echo "[ -d ${HOME}/go ] && export GOROOT=\"${HOME}/go\";"; \
  echo "[ -d ${HOME}/go ] && export GOPATH=\"${HOME}/go-packages\";"; \
  echo "[ -d ${HOME}/go ] && export PATH=\"\${PATH}:\${GOROOT}/bin:\${GOPATH}/bin\";"; \
  echo '[ $(command -v gp) ] && eval $(gp env | sed -e "s/\(=[[:blank:]]*\)\(.*\)/\1\"\2\" ;/" -e "s/^/export /g")' ; \
  echo '[ $(command -v starship) ] && eval "$(starship init bash)";' ; \
  echo '[ $(command -v just) ] && alias j="just";' ; \
  ) | tee -a /etc/bash.bashrc > /dev/null ;
USER gitpod
RUN \
  sed -i \
  -e '/^\s*#/d' \
  -e '/PIP_USER/d' \
  -e '/cargo/d' \
  -e '/^\s*$/d' \
  "${HOME}/.profile" \
  && ( \
  echo "export EDITOR=\"${EDITOR}\"" ; \
  echo "export VISUAL=\"${VISUAL}\"" ; \
  echo "export PIP_USER=\"${PIP_USER}\"" ; \
  ) | tee -a "${HOME}/.profile" > /dev/null ;
COPY "./.SpaceVim.d" "${HOME}/.SpaceVim.d"
# hadolint ignore=SC2015
RUN \
  curl -sLf https://spacevim.org/install.sh | bash \
  && timeout 150 \
  nvim --headless \
  +VimProcInstall \
  +'call dein#install()'  \
  +UpdateRemotePlugins  \
  +qall || true ;
COPY --from=just /workspace /usr/local/bin
COPY --from=convco /workspace /usr/local/bin
COPY --from=clog /workspace /usr/local/bin
# ─── CLEAN UPS ──────────────────────────────────────────────────────────────────
WORKDIR /workspace
