# syntax = docker/dockerfile-upstream:master-labs
#-*-mode:dockerfile;indent-tabs-mode:nil;tab-width:2;coding:utf-8-*-
# vi: ft=dockerfile tabstop=2 shiftwidth=2 softtabstop=2 expandtab:

# ─── PREBUILD IMAGE ALIASES ─────────────────────────────────────────────────────
# hadolint ignore=DL3007
FROM fjolsvin/vscode-golang-alpine:latest as vscode-golang-alpine
# ─── BUILDER ────────────────────────────────────────────────────────────────────
FROM vscode-golang-alpine
USER root
ARG IMAGE_SPECIFIC_PACKAGES="\
  just \
  just-bash-completion \
  "
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
# hadolint ignore=DL3018,SC2086
RUN \
  set -xeu ; \
  apk add --no-cache ${IMAGE_SPECIFIC_PACKAGES} || \
  (sed -i -e 's/dl-cdn/dl-4/g' /etc/apk/repositories && \
  apk add --no-cache ${IMAGE_SPECIFIC_PACKAGES}) \
  && touch "/var/run/docker.sock" && \
  chown "$(id -u ${USER}):docker" "/var/run/docker.sock" && \
  find . \
  -not -group $(id -g) \
  -not -user $(id -u) \
  -print0 | xargs \
  -P 0 \
  -0 \
  --no-run-if-empty \
  chown --no-dereference "$(id -u ${USER}):$(id -g ${USER})" && \
  find ~/ \
  -not -group $(id -g) \
  -not -user $(id -u) \
  -print0 | xargs \
  -P 0 \
  -0 \
  --no-run-if-empty \
  chown --no-dereference "$(id -u ${USER}):$(id -g ${USER})"
USER "${USER}"
WORKDIR "${WORKDIR}"
ENTRYPOINT [ "" ]
