# !/usr/bin/env -S just --justfile
# vim: filetype=just softtabstop=2 tabstop=2 shiftwidth=2 fenc=utf-8 fileformat=unix expandtab

set positional-arguments := true
set dotenv-load := true
set shell := ["/bin/bash", "-o", "pipefail", "-c"]

project_name := `basename $PWD`
docker_image := `/bin/grep '"image":' .devcontainer/devcontainer.json | /bin/sed -e 's/"image": //g' -e 's/"//g' -e 's/[[:space:],]*//g'`

# ────────────────────────────────────────────────────────────────────
default:
    @just --choose

# ─── DEPENDENCIES ───────────────────────────────────────────────────────────────

alias b := bootstrap

bootstrap: dependencies go-bootstrap vscode-tasks kary-comments format pre-commit
    @echo bootstrap completed

# ────────────────────────────────────────────────────────────────────────────────

alias d := dependencies

dependencies:
    #!/usr/bin/env bash
    set -euo pipefail
    IFS=':' read -a paths <<< "$(printenv PATH)" ;
    [[ ! " ${paths[@]} " =~ " ${HOME}/bin " ]] && export PATH="${PATH}:${HOME}/bin" || true ;
    if command -- sudo pip3 -h > /dev/null 2>&1 ; then
      if ! command -- pre-commit -h > /dev/null 2>&1 ; then
        curl "https://pre-commit.com/install-local.py" | python3 -
      fi
      if ! command -- jq -h > /dev/null 2>&1 ; then
        python3 -m pip install --user jq
      fi
    fi
    if command -- cargo -h > /dev/null 2>&1 ; then
      if ! command -- convco -h > /dev/null 2>&1 ; then
        cargo install --locked --all-features --git https://github.com/convco/convco.git
      fi
      if ! command -- jsonfmt -h > /dev/null 2>&1 ; then
        cargo install --locked --all-features -j `nproc` jsonfmt
      fi
    fi

# ────────────────────────────────────────────────────────────────────────────────

alias gb := go-bootstrap

go-bootstrap:
    go env -w "GO111MODULE=on"
    go env -w "CGO_ENABLED=0"
    go env -w "CGO_LDFLAGS=-s -w -extldflags '-static'"
    go clean -modcache
    go mod tidy
    go generate -tags tools tools.go

# ────────────────────────────────────────────────────────────────────────────────

alias kc := kary-comments

kary-comments:
    #!/usr/bin/env bash
    set -euo pipefail
    sed -i.bak \
    -e "/case 'yaml':.*/a case 'dockerfile':" \
    -e "/case 'yaml':.*/a case 'just':" \
    -e "/case 'yaml':.*/a case 'hcl':" \
    ~/.vscode*/extensions/karyfoundation.comment*/dictionary.js > /dev/null 2>&1 || true

# ────────────────────────────────────────────────────────────────────────────────

alias vt := vscode-tasks

vscode-tasks:
    #!/usr/bin/env bash
    set -euo pipefail
    if command -- jq -h > /dev/null 2>&1 ; then
      IFS=' ' read -a TASKS <<< "$(just --summary --color never -f "{{ justfile() }}" 2>/dev/null)"
      if [ ${#TASKS[@]} -ne 0  ];then
        mkdir -p "{{ justfile_directory() }}/.vscode"
        json=$(jq -n --arg version "2.0.0" '{"version":$version,"tasks":[]}')
        for task in "${TASKS[@]}";do
          taskjson=$(jq -n --arg task "${task}" --arg command "just ${task}" '[{"type": "shell","label": $task,  "command": $command }]')
          json=$(echo "${json}" | jq ".tasks += ${taskjson}")
        done
        echo "${json}" | jq -r '.' > "{{ justfile_directory() }}/.vscode/tasks.json"
      fi
    fi
    just format-just

# ─── FORMAT ─────────────────────────────────────────────────────────────────────

alias f := format
alias fmt := format

format: format-json format-just format-go
    @echo format completed

# ────────────────────────────────────────────────────────────────────────────────

alias fj := format-json
alias json-fmt := format-json

format-json:
    #!/usr/bin/env bash
    set -euo pipefail
    if command -- jsonfmt -h > /dev/null 2>&1 ; then
      while read file;do
        echo "*** formatting $file"
        jsonfmt "$file" || true
        echo '' >> "$file"
      done < <(find -type f -not -path '*/\.git/*' -name '*.json')
    fi

# ────────────────────────────────────────────────────────────────────────────────

alias go-fmt := format-go
alias gofmt := format-go
alias fg := format-go

format-go:
    #!/usr/bin/env bash
    set -euo pipefail
    gofmt -l -w {{ justfile_directory() }}

# ────────────────────────────────────────────────────────────────────────────────
format-just:
    #!/usr/bin/env bash
    set -euo pipefail
    just --unstable --fmt 2>/dev/null

# ─── GO ─────────────────────────────────────────────────────────────────────────
build:
    #!/usr/bin/env bash
    set -euo pipefail
    mage build

# ────────────────────────────────────────────────────────────────────────────────

lint: golangci-lint

golangci-lint: format-go
    #!/usr/bin/env bash
    set -euox pipefail
    golangci-lint run \
    --print-issued-lines=false \
    --exclude-use-default=false \
    --config "{{ justfile_directory() }}/.golangci.yml"

# ─── GIT ────────────────────────────────────────────────────────────────────────
# Variables

MASTER_BRANCH_NAME := 'tmo/main'
MAJOR_VERSION := `[[ ! -z $(git tag -l | head -n 1 ) ]] && convco version --major 2>/dev/null || echo '0.0.1'`
MINOR_VERSION := `[[ ! -z $(git tag -l | head -n 1 ) ]] && convco version --minor 2>/dev/null || echo '0.0.1'`
PATCH_VERSION := `[[ ! -z $(git tag -l | head -n 1 ) ]] && convco version --patch 2>/dev/null || echo '0.0.1'`

# ────────────────────────────────────────────────────────────────────────────────

alias pc := pre-commit

pre-commit: format-just
    #!/usr/bin/env bash
    set -euo pipefail
    IFS=':' read -a paths <<< "$(printenv PATH)" ;
    [[ ! " ${paths[@]} " =~ " ${HOME}/bin " ]] && export PATH="${PATH}:${HOME}/bin" || true ;
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    if [ -r .pre-commit-config.yaml ]; then
      git add ".pre-commit-config.yaml"
      pre-commit install > /dev/null 2>&1
      pre-commit install-hooks
      pre-commit
    fi
    popd > /dev/null 2>&1

# ────────────────────────────────────────────────────────────────────────────────

alias gf := git-fetch

git-fetch:
    #!/usr/bin/env bash
    set -euo pipefail
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    git fetch -p ;
    for branch in $(git branch -vv | grep ': gone]' | awk '{print $1}'); do
      git branch -D "$branch";
    done
    popd > /dev/null 2>&1

# ────────────────────────────────────────────────────────────────────────────────

alias c := commit

commit: pre-commit git-fetch
    #!/usr/bin/env bash
    set -euo pipefail
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    if command -- convco -h > /dev/null 2>&1 ; then
      convco commit
    else
      git commit
    fi
    popd > /dev/null 2>&1

# ────────────────────────────────────────────────────────────────────────────────

alias mar := major-release

major-release: git-fetch format-just
    #!/usr/bin/env bash
    set -euo pipefail
    IFS=':' read -a paths <<< "$(printenv PATH)" ;
    [[ ! " ${paths[@]} " =~ " ${HOME}/bin " ]] && export PATH="${PATH}:${HOME}/bin" || true;
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    git checkout "{{ MASTER_BRANCH_NAME }}"
    git pull
    git tag -a "v{{ MAJOR_VERSION }}" -m 'major release {{ MAJOR_VERSION }}'
    git push origin --tags
    if command -- convco -h > /dev/null 2>&1 ; then
      convco changelog > CHANGELOG.md
      git add CHANGELOG.md
      if command -- pre-commit -h > /dev/null 2>&1 ; then
        pre-commit || true
        git add CHANGELOG.md
      fi
      git commit -m 'docs(changelog): updated changelog for v{{ MAJOR_VERSION }}'
      git push
    fi
    just git-fetch
    popd > /dev/null 2>&1

# ────────────────────────────────────────────────────────────────────────────────

alias mir := minor-release

minor-release: git-fetch format-just
    #!/usr/bin/env bash
    set -euo pipefail
    IFS=':' read -a paths <<< "$(printenv PATH)" ;
    [[ ! " ${paths[@]} " =~ " ${HOME}/bin " ]] && export PATH="${PATH}:${HOME}/bin" || true;
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    git checkout "{{ MASTER_BRANCH_NAME }}"
    git pull
    git tag -a "v{{ MINOR_VERSION }}" -m 'minor release {{ MINOR_VERSION }}'
    git push origin --tags
    if command -- convco -h > /dev/null 2>&1 ; then
      convco changelog > CHANGELOG.md
      git add CHANGELOG.md
      if command -- pre-commit -h > /dev/null 2>&1 ; then
        pre-commit || true
        git add CHANGELOG.md
      fi
      git commit -m 'docs(changelog): updated changelog for v{{ MINOR_VERSION }}'
      git push
      just git-fetch
    fi
    popd > /dev/null 2>&1

# ────────────────────────────────────────────────────────────────────────────────

alias pr := patch-release

patch-release: git-fetch format-just
    #!/usr/bin/env bash
    set -euo pipefail
    IFS=':' read -a paths <<< "$(printenv PATH)" ;
    [[ ! " ${paths[@]} " =~ " ${HOME}/bin " ]] && export PATH="${PATH}:${HOME}/bin" || true;
    pushd "{{ justfile_directory() }}" > /dev/null 2>&1
    git checkout "{{ MASTER_BRANCH_NAME }}"
    git pull
    git tag -a "v{{ PATCH_VERSION }}" -m 'patch release {{ PATCH_VERSION }}'
    git push origin --tags
    if command -- convco -h > /dev/null 2>&1 ; then
      convco changelog > CHANGELOG.md
      git add CHANGELOG.md
      if command -- pre-commit -h > /dev/null 2>&1 ; then
        pre-commit || true
        git add CHANGELOG.md
      fi
      git commit -m 'docs(changelog): updated changelog for v{{ MINOR_VERSION }}'
      git push
    fi
    just git-fetch
    popd > /dev/null 2>&1

alias gc := generate-changelog

generate-changelog: format-just
    #!/usr/bin/env bash
    set -euo pipefail
    rm -rf "{{ justfile_directory() }}/tmp"
    mkdir -p "{{ justfile_directory() }}/tmp"
    convco changelog > "{{ justfile_directory() }}/tmp/$(basename {{ justfile_directory() }})-changelog-$(date -u +%Y-%m-%d).md"
    if command -- pandoc -h >/dev/null 2>&1; then
    pandoc \
      --from markdown \
      --pdf-engine=xelatex \
      -o  "{{ justfile_directory() }}/tmp/$(basename {{ justfile_directory() }})-changelog-$(date -u +%Y-%m-%d).pdf" \
      "{{ justfile_directory() }}/tmp/$(basename {{ justfile_directory() }})-changelog-$(date -u +%Y-%m-%d).md"
    fi
    if [ -d /workspace ]; then
      cp -f "{{ justfile_directory() }}/tmp/$(basename {{ justfile_directory() }})-changelog-$(date -u +%Y-%m-%d).pdf" /workspace/
      cp -f "{{ justfile_directory() }}/tmp/$(basename {{ justfile_directory() }})-changelog-$(date -u +%Y-%m-%d).md" /workspace/
    fi

snapshot: git-fetch format-just
    #!/usr/bin/env bash
    set -euo pipefail
    sync
    snapshot_dir="{{ justfile_directory() }}/tmp/snapshots"
    mkdir -p "${snapshot_dir}"
    time="$(date +'%Y-%m-%d-%H-%M')"
    path="${snapshot_dir}/${time}.tar.gz"
    tmp="$(mktemp -d)"
    tar -C {{ justfile_directory() }} -cpzf "$tmp/${time}.tar.gz" .
    mv "$tmp/${time}.tar.gz" "$path"
    rm -r "$tmp"
    echo >&2 "*** snapshot created at ${path}"

# ─── DEVCONTAINER SETUP ─────────────────────────────────────────────────────────

alias dcp := dev-container-pull

dev-container-pull:
    #!/usr/bin/env bash
    set -xeuo pipefail
    docker pull "{{ docker_image }}"

# ────────────────────────────────────────────────────────────────────────────────

alias dcc := dev-container-clean

dev-container-clean:
    #!/usr/bin/env bash
    set -xeuo pipefail
    docker rm -f "{{ project_name }}" > /dev/null 2>&1 || true

# ────────────────────────────────────────────────────────────────────────────────

alias dcu := dev-container-up

dev-container-up: dev-container-pull dev-container-clean
    #!/usr/bin/env bash
    set -xeuo pipefail
    if [ ! "$(docker ps -q -f name="{{ project_name }}")" ]; then
      # [ NOTE ] => this code path is just for safety purposes
      # it is probably never going to run.
      if [ "$(docker ps -aq -f status=exited -f name="{{ project_name }}")" ]; then
        docker rm "{{ project_name }}"
      fi
      docker run \
        --detach \
        --tty \
        --name "{{ project_name }}" \
        --volume "{{ justfile_directory() }}:/workspace" \
        --workdir "/workspace" \
        --network "host" \
        --user "$(id -u):$(id -g)" \
        --cap-add SYS_PTRACE \
        --security-opt seccomp=unconfined \
        "{{ docker_image }}";
        sleep 5 ;
    fi

# ────────────────────────────────────────────────────────────────────────────────

alias dce := dev-container-exec

dev-container-exec: dev-container-up
    #!/usr/bin/env bash
    set -xeuo pipefail
    docker exec --interactive --tty "{{ project_name }}" /bin/bash

# ─── GITPOD RELATED TARGET ──────────────────────────────────────────────────────
docker-socket-chown:
    #!/usr/bin/env bash
    set -euo pipefail
    sudo chown "$(id -u gitpod):$(cut -d: -f3 < <(getent group docker))" /var/run/docker.sock

# ────────────────────────────────────────────────────────────────────────────────

alias fo := fix-ownership

fix-ownership: docker-socket-chown
    #!/usr/bin/env bash
    set -euo pipefail
    sudo find "${HOME}/" "/workspace" -not -group `id -g` -not -user `id -u` -print0 | xargs -P 0 -0 --no-run-if-empty sudo chown --no-dereference "`id -u`:`id -g`" || true ;
    # sudo find "/workspace" -not -group `id -g` -not -user `id -u` -print | xargs -I {}  -P `nproc` --no-run-if-empty sudo chown --no-dereference "`id -u`:`id -g`" {} || true ;

# ────────────────────────────────────────────────────────────────────────────────
docker-login-env:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "*** ensuring current user belongs to docker group" ;
    sudo usermod -aG docker "$(whoami)"
    echo "*** ensuring required environment variables are present" ;
    while [ -z "$DOCKER_USERNAME" ] ; do \
    printf "\n❗ The DOCKER_USERNAME environment variable is required. Please enter its value.\n" ;
    read -s -p "DOCKER_USERNAME: " DOCKER_USERNAME ; \
    done ; gp env DOCKER_USERNAME=$DOCKER_USERNAME && printf "\nThanks\n" || true ;
    while [ -z "$DOCKER_PASSWORD" ] ; do \
    printf "\n❗ The DOCKER_PASSWORD environment variable is required. Please enter its value.\n" ;
    read -s -p "DOCKER_PASSWORD: " DOCKER_PASSWORD ; \
    done ; gp env DOCKER_PASSWORD=$DOCKER_PASSWORD && printf "\nThanks\n" || true ;

# ────────────────────────────────────────────────────────────────────────────────

alias dl := docker-login

docker-login: fix-ownership docker-login-env
    #!/usr/bin/env bash
    set -euo pipefail
    echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin ;
    just fix-ownership

# ────────────────────────────────────────────────────────────────────────────────

alias gp := gitpod

gitpod:
    #!/usr/bin/env bash
    set -euxo pipefail
    bash "{{ justfile_directory() }}/.gp/build.sh"

# ────────────────────────────────────────────────────────────────────────────────
ssh-pub-key-env:
    #!/usr/bin/env bash
    set -euo pipefail
    while [ -z "$SSH_PUB_KEY" ] ; do \
    printf "\n❗ The SSH_PUB_KEY environment variable is required. Please enter its value.\n" ;
    read -s -p "SSH_PUB_KEY: " SSH_PUB_KEY ; \
    done ; gp env SSH_PUB_KEY=$SSH_PUB_KEY && printf "\nThanks\n" || true ;

# ────────────────────────────────────────────────────────────────────────────────
ssh-pub-key: fix-ownership ssh-pub-key-env
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p ${HOME}/.ssh ;
    echo "${SSH_PUB_KEY}" | tee ${HOME}/.ssh/authorized_keys > /dev/null ;
    chmod 700 ${HOME}/.ssh ;
    chmod 600 ${HOME}/.ssh/authorized_keys ;
    just fix-ownership
    exit 0

# ────────────────────────────────────────────────────────────────────────────────
chisel: fix-ownership
    #!/usr/bin/env bash
    set -euo pipefail
    [ -f ${HOME}/chisel.pid ] && echo "*** killing chisel server" && kill -9 "$(cat ${HOME}/chisel.pid)" && rm -rf ${HOME}/chisel.pid ;
    pushd ${HOME}/ ;
    echo "*** starting chisel server" ;
    bash -c "chisel server --socks5 --pid > ${HOME}/chisel.log 2>&1 &" ;
    echo "*** chisel was started successfully" ;
    popd ;
    just fix-ownership
    exit 0

# ────────────────────────────────────────────────────────────────────────────────
dropbear: fix-ownership
    #!/usr/bin/env bash
    set -euo pipefail
    [ ! -f ${HOME}/dropbear.hostkey ] && echo "*** generating dropbear host key" && dropbearkey -t rsa -f ${HOME}/dropbear.hostkey ;
    [ -f ${HOME}/dropbear.pid ] && echo "*** killing dropbear server" && kill -9 "$(cat ${HOME}/dropbear.pid)" && rm -rf ${HOME}/dropbear.pid ;
    echo "*** starting dropbear server" ;
    bash -c "dropbear -r ${HOME}/dropbear.hostkey -F -E -s -p 2222 -P ${HOME}/dropbear.pid > ${HOME}/dropbear.log 2>&1 &" ;
    echo "*** dropbear server was started successfully" ;
    just fix-ownership
    exit 0

# ────────────────────────────────────────────────────────────────────────────────

alias ssh := ssh-config

ssh-config: ssh-pub-key
    #!/usr/bin/env bash
    set -euo pipefail
    cat << EOF
    Host $(gp url | sed -e 's/https:\/\///g' -e 's/[.].*$//g')
      HostName localhost
      User gitpod
      Port 2222
      ProxyCommand chisel client $(gp url 8080) stdio:%h:%p
      RemoteCommand cd /workspace && exec bash --login
      RequestTTY yes
      IdentityFile ~/.ssh/id_rsa
      IdentitiesOnly yes
      StrictHostKeyChecking no
      CheckHostIP no
      MACs hmac-sha2-256
      UserKnownHostsFile /dev/null
    EOF
