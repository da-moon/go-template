#-*-mode:hcl;indent-tabs-mode:nil;tab-width:2;coding:utf-8-*-
# vi: ft=hcl tabstop=2 shiftwidth=2 softtabstop=2 expandtab:

# [ NOTE ] => clean up buildx builders
# docker buildx ls | awk '$2 ~ /^docker(-container)*$/{print $1}' | xargs -r -I {} docker buildx rm {}
# [ NOTE ] create a builder for this file
# docker buildx create --use --name "go-template" --driver docker-container
# [ NOTE ] run build without pushing to dockerhub
# LOCAL=true docker buildx bake --builder go-template

variable "LOCAL" {default=false}
variable "DOCKER_IMAGE" {default="fjolsvin/go-template"}
variable "TAG" {default=""}
group "default" {
    targets = [
      "prd",
    ]
}
# LOCAL=true docker buildx bake --builder go-template prd
target "prd" {
    context="."
    dockerfile = "contrib/docker/prd/Dockerfile"
    tags = [
        "${DOCKER_IMAGE}:latest",
        notequal("",TAG) ? "${DOCKER_IMAGE}:${TAG}": "",
    ]
    platforms = ["linux/amd64"]
    cache-from = ["type=registry,ref=${DOCKER_IMAGE}:cache"]
    cache-to   = ["type=registry,mode=max,ref=${DOCKER_IMAGE}:cache"]
    output     = [equal(LOCAL,true) ? "type=docker" : "type=registry"]
}
