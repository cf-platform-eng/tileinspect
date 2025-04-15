ARG base_image=tas-ecosystem-docker-virtual.usw1.packages.broadcom.com/ubuntu:latest

FROM ${base_image}

LABEL maintainer="Pivotal Platform Engineering ISV-CI Team <cf-isv-dashboard@pivotal.io>"

COPY build/tileinspect-linux /usr/local/bin/tileinspect

ENTRYPOINT [ "tileinspect" ]