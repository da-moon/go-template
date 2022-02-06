# syntax = docker/dockerfile:1.3-labs
#-*-mode:dockerfile;indent-tabs-mode:nil;tab-width:2;coding:utf-8-*-
# vim: filetype=dockerfile softtabstop=2 tabstop=2 shiftwidth=2 fenc=utf-8 fileformat=unix expandtab

#
#  ──── COMPRESSION LAYER ─────────────────────────────────────────────
#
FROM alpine:edge as upx
# ────────────────────────────────────────────────────────────────────────────────
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# ────────────────────────────────────────────────────────────────────────────────
USER root
RUN apk add --no-cache "bash~=5.1"
# ────────────────────────────────────────────────────────────────────────────────
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# ────────────────────────────────────────────────────────────────────────────────
ARG BASE_PACKAGES="\
  coreutils~=8.32 \
  findutils~=4.8 \
  binutils~=2.35 \
  file~=5.40 \
  "
RUN \
  IFS=' ' read -a packages <<< $BASE_PACKAGES ; \
  ( \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/v3.14/main" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/v3.14/community" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/v3.13/main" ; \
  echo "http://dl-cdn.alpinelinux.org/alpine/v3.13/community" ; \
  ) | tee /etc/apk/repositories > /dev/null  \
  && apk add --no-cache "${packages[@]}" \
  || ( \
  sed -i -e 's/dl-cdn/dl-4/g' /etc/apk/repositories \
  && apk add --no-cache "${packages[@]}" \
  )
ARG UPX_DEPS="\
  curl~=7.79 \
  jq~=1.6 \
  xz~=5.2 \
  "
RUN \
  IFS=' ' read -a packages <<< $UPX_DEPS ; \
  apk add --no-cache --virtual .deps-upx "${packages[@]}" \
  || ( \
  sed -i -e 's/dl-cdn/dl-4/g' /etc/apk/repositories \
  && apk add --no-cache --virtual .deps-upx "${packages[@]}" \
  )
ARG REPO="upx/upx"
ARG LATEST_RELEASE_ENDPOINT="https://api.github.com/repos/${REPO}/releases/latest"
RUN \
  tag_name="$(curl -sL ${LATEST_RELEASE_ENDPOINT} | jq -r '.tag_name')"; \
  architecture="$(apk --print-arch)"; \
  case "$architecture" in \
    x86_64|amd64) \
      architecture="amd64" \
    ;; \
    aarch64) \
      architecture="arm64" \
    ;; \
    *) \
      echo >&2 "[ WARN ] compression utilities are not available: $architecture"; \
      exit 0 \
    ;; \
  esac ; \
  version="$(echo ${tag_name} | sed 's/v//g')"; \
  download_url="https://github.com/upx/upx/releases/download/${tag_name}/upx-${version}-${architecture}_linux.tar.xz"; \
  rm -rf \
    /tmp/{upx.tar,upx.tar.xz} \
    /usr/local/bin/upx \
  && echo "$download_url" > /tmp/dl \
  && curl -fsSLo /tmp/upx.tar.xz "${download_url}" \
  && xz -d -c /tmp/upx.tar.xz \
  | tar \
    -xOf - upx-${version}-${architecture}_linux/upx > /usr/local/bin/upx
#  ────────────────────────────────────────────────────────────────────
COPY <<-"EOT" /usr/local/bin/compress
#!/usr/bin/env bash
set -euo pipefail
apkArch="$(apk --print-arch)";
case "$apkArch" in
  x86_64|aarch64)
     find . \
    -type f \
    -executable \
    -exec sh -c "file -i '{}' |  grep -q 'charset=binary'" \; \
    -print \
    | xargs \
        -P `nproc` \
        --no-run-if-empty \
        bash -c '
      for pkg do
          strip "$pkg" || true ;
          upx "$pkg" || true ;
      done' bash ;
    ;;
  *)
    echo >&2 "[ WARN ] compression utilities are not available: $apkArch";
    exit 1
  ;;
esac;
find . \
  -mindepth 2 \
  -type f \
  -executable \
  -exec sh -c "file -i '{}' | grep -q 'charset=binary'" \; \
  -print \
  | xargs -P `nproc` -I {} --no-run-if-empty \
    mv "{}" ./
find . -mindepth 1 -maxdepth 1 -type d -exec rm -r "{}" \;
EOT
#  ────────────────────────────────────────────────────────────────────

# [ NOTE ] deleting dependencies on aarch64 leads to the following error
# PROT_EXEC|PROT_WRITE failed.
RUN \
  chmod a+x /usr/local/bin/* \
  && upx --version \
  [ "$(apk --print-arch)" == "aarch64" ] && apk del .deps-upx || true ; \
  rm -rf \
    /var/cache/apk/* \
    /var/tmp/* \
    /tmp/* ;
WORKDIR "/workspace"
# ─── BUILD STAGE ────────────────────────────────────────────────────────────────
FROM golang:alpine AS build-stage
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018
RUN \
  apk add --no-cache \
    "make~=4.3" \
    "ncurses~=6.2" \
    "git~=2.32" \
    "build-base~=0.5" \
    "bash~=5.1"
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# ─── INSTALL MAGE ───────────────────────────────────────────────────────────────
RUN \
  set -ex && \
  git clone "https://github.com/magefile/mage" "/tmp/mage"
WORKDIR /tmp/mage
RUN \
  --mount=type=cache,target="/root/.cache/go-build" \
  --mount=type=cache,target="/go/pkg/mod" \
  --mount=type=tmpfs,target="/go/src/" \
  set -xeu ; \
  go run bootstrap.go
# ─── BUILD ──────────────────────────────────────────────────────────────────────
WORKDIR /workspace
COPY . .
#  ────────────────────────────────────────────────────────────────────
WORKDIR /workspace
# hadolint ignore=SC2239
RUN \
  --mount=type=cache,target="/root/.cache/go-build" \
  --mount=type=cache,target="/go/pkg/mod" \
  --mount=type=tmpfs,target="/go/src/" \
  set -xeu ; \
  mage -v build \
  && echo "nobody:x:65534:65534:Nobody:/:" > "/etc_passwd"
# ─── COMPRESSION STAGE ──────────────────────────────────────────────────────────
FROM upx AS compression-stage
COPY --from=build-stage "/workspace/bin" "/workspace"
RUN \
  compress ;
# ─── MAIN LAYER ─────────────────────────────────────────────────────────────────

# [ NOTE ] => alpine is used for cases where we are running a server and
# want to take advantage of Docker's built in health checks.
#
# e.g :
# HEALTHCHECK --interval=5m --timeout=5s CMD curl -f "http://${API_ADDR}/health" || exit 1
#
# If you are note building a server, you can use 'scratch' image
FROM alpine:3.14
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018
RUN apk add --no-cache "curl~=7.79"
COPY --from=build-stage "/etc_passwd" "/etc/passwd"
USER nobody
WORKDIR /workspace
COPY --from=build-stage "/etc/ssl/certs/ca-certificates.crt" "/etc/ssl/certs/"
# [ TODO ] => change this to match your binary name
COPY --from=compression-stage --chown=65534:0 "/workspace" "/workspace"

ENTRYPOINT [ "/entrypoint" ]
CMD [ "version" ]
