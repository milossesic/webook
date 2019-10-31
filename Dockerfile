# Please update your base container regularly for bug fixes and security patches.
# See https://git.corp.adobe.com/ASR/bbc-factory for the latest BBC releases.
FROM docker-asr-release.dr.corp.adobe.com/asr/golang:2.11-alpine

ENV GOWORKDIR="/go/src/git.corp.adobe.com/dc/notifications_load_test"

COPY container/root /

COPY config ${GOWORKDIR}/config

COPY main ${GOWORKDIR}/main

WORKDIR ${GOWORKDIR}
