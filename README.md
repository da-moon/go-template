# go-template

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/da-moon/go-template)

- [go-template](#go-template)
  - [overview](#overview)
  - [setup](#setup)
    - [Rust ToolChain](#rust-toolchain)
    - [Convco](#convco)
    - [Just](#just)
    - [Bootstrap Repository Go Tools](#bootstrap-repository-go-tools)
  - [docker images](#docker-images)
    - [prerequisites](#prerequisites)
    - [build scripts](#build-scripts)
    - [Github Actions](#github-actions)
  - [Spacevim](#spacevim)
  - [git](#git)

## overview

The purpose of this template repository is to help with quickly bootstrapping a
Golang project. It has the following features :

- Sane [`Spacevim`][spacevim-url] configuration
- Virtualized Development Environment
  - Sane [`VSCode-remote-containers`][vsc-rc-url] configuration and
Alpine Backed Docker container with support for `AMD64` and `AARCH64` CPUs.
  - Vagrantfile Debian 10 based development
environment configured to work with `Virtualbox,`
`Hyper-V,` `Libvirt` and `Google-Cloud` providers.
  - Custom Gitpod Image that enables one to use the Gitpod
as a backend for [`vsc-rs-url`][vsc-rs-url]
- Common [`revive`][revive-url] and [`golang-ci`][golang-ci-url] Linter configuration
- `Makefile` with targets for
Semantic Versioning
- [`Justfile`][just-url] targets for bootstrapping base dependencies,
Semantic Versioning and development environment-related tasks.
- [pre-commit][pc-url]: configuration and some common hooks
- Common dev tools, managed in [tools.go][go-tools-url]
- Github actions workflows for building docker images

## setup

The virtualized environments have all the needed tools pre-provisioned.
If you are using this repo to bootstrap projects out of the virtualized environments,
you need to have the following installed.

- OS: Linux
- Go Toolchain >= 1.16.3
- `Rust` Toolchain for building `convco` and `just`
- `Ripgrep`
- Python for installing `pre-commit.
- Spacevim
- NodeJs and Yarn for installing Spacevim Language Servers and
some pre-commit hooks

I won't cover NodeJS, Python, Go-Toolchain, and Pre-commit
setup in this document. I would explain how to install rust-based dependencies.
Make sure to install `Ripgrep` through your favorite package manager

### Rust ToolChain

- install the stable toolchain and configure `~/.bashrc`

```bash
curl \
  --proto '=https' \
  --tlsv1.2 -sSf \
  https://sh.rustup.rs \
  | sh -s -- \
    -y \
    --no-modify-path \
    --default-toolchain stable \
    --profile default  ;
if ! grep -q "cargo" ~/.bashrc; then
  echo '[ -r $HOME/.cargo/env ] && . $HOME/.cargo/env' >> ~/.bashrc ;
  source ~/.bashrc ;
fi
rustup default stable
```

### Convco

- make sure `cmake` is installed.
- run `cargo install convco`

### Just

- run `cargo install just`

### Bootstrap Repository Go Tools

- ensure `just` is installed
- run `just bootstrap`

## docker images

all images come with a build script
that would prefer to build the image with `buildx`
command.

In case the `buildx` plugin is not available,
then the script would fall back on the basic `build` command

### prerequisites

- update `IMAGE_NAME` variable in build scripts

### build scripts

- login to your container registry of choice
- Docker image used with VScode remote-containers extension

```bash
bash contrib/docker/devcontainer/alpine/build.sh
```

- Gitpod Docker Image

```bash
bash .gp/build.sh
```

### Github Actions

- add `DOCKER_PASSWORD` and `DOCKER_USERNAME` to your repository's secrets

## Spacevim

- I would encourage use `neovim >= 0.5`. Since the
release is not available in most package managers, either
install it from the `Snap` store's edge channel or
build it from the source.
- Spacevim needs Python and Nodejs
- Spacevim has some font package dependencies. for instance,
on Debian-based distributions, install `xfonts-utils` and `fonts-symbola` .
- Install Spacevim

```bash
curl -sLf https://spacevim.org/install.sh | bash
```

- install extra dependencies. Make sure `yarn` is installed before running
the following command

```bash
just spacevim-dependencies
# or
just sd
```

## git

- the following target
ensures pre-commit hooks are installed
and it would run pre-commit hooks.

```bash
just pre-commit
# or
just pc
```

- use the following target for commits. this target
would run `pre-commit` target as a dependency.

```bash
just target
# or
just c
```

- use the following target submitting a patch-release and
generate release changelog.

```bash
just patch-release
# or
just pr
```

- use the following target submitting a minor-release and
generate release changelog.

```bash
just minor-release
# or
just mir
```

- use the following target submitting a major-release and
generate release changelog.

```bash
just major-release
# or
just mar
```

[spacevim-url]: https://spacevim.org/quick-start-guide
[vsc-rc-url]: https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers
[vsc-rs-url]: https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh
[revive-url]: https://github.com/mgechev/revive
[golangci-url]: https://github.com/golangci/golangci-lint
[just-url]: https://github.com/casey/just
[pc-url]: https://pre-commit.com
[go-tools-url]: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
