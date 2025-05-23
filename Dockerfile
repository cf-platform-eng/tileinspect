ARG base_image=tas-ecosystem-docker-virtual.usw1.packages.broadcom.com/ubuntu:latest

FROM ${base_image}

LABEL maintainer="Pivotal Platform Engineering ISV-CI Team <cf-isv-dashboard@pivotal.io>"

COPY tileinspect-build/tileinspect-linux-amd64 /usr/local/bin/tileinspect

ENTRYPOINT [ "tileinspect" ]
