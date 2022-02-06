#!/usr/bin/env bash
set -xeuo pipefail
pushd "$WD" > /dev/null 2>&1
if [ -z ${IMAGE_NAME+x} ] || [  -z ${IMAGE_NAME+x} ];then
  IMAGE_NAME="fjolsvin/gitpod-$(basename "$WD")"
fi
IMAGE_NAME="${IMAGE_NAME}/alpine-devcontainer"
CACHE_NAME="${IMAGE_NAME}:cache"
export DOCKER_BUILDKIT=1
BUILD="docker"
if [[ $(docker buildx version 2> /dev/null ) ]]; then
  builder="$(echo "$IMAGE_NAME" | cut -d/ -f2)"
  BUILD+=" buildx build"
  BUILD+=" --file .gp/Dockerfile"
  BUILD+=" --cache-from type=registry,ref=${CACHE_NAME}"
  BUILD+=" --cache-to type=registry,mode=max,ref=${CACHE_NAME}"
  BUILD+=" --tag ${IMAGE_NAME}:latest"
  BUILD+=" --progress=plain"
  BUILD+=" --push"
  docker buildx use "${builder}" || docker buildx create --use --name "${builder}"
else
  BUILD+=" build"
  BUILD+=" --file .gp/Dockerfile"
  BUILD+=" --tag ${IMAGE_NAME}:latest"
  BUILD+=" --cache-from type=registry,ref=${CACHE_NAME}"
  BUILD+=" --progress=plain"
  BUILD+=" --pull"
fi
$BUILD $WD
if [[ $(docker buildx version 2> /dev/null ) ]]; then
  docker buildx use default
else
  PUSH="docker push"
  PUSH+=" ${IMAGE_NAME}:latest"
  $PUSH
fi
popd > /dev/null 2>&1
