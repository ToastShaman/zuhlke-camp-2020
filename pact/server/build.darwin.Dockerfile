FROM golang:1.15.5 AS build-stage

ARG VERSION
ARG ARTEFACT_NAME

RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=darwin \
  GOARCH=amd64

WORKDIR /src
COPY . .

RUN go build -v \
    -trimpath \
    -ldflags "-s -w -extldflags '-static' -X main.version=${VERSION}" \
    -o /out/${ARTEFACT_NAME}-${GOOS}-${GOARCH}

FROM scratch AS export-stage
COPY --from=build-stage /out/ /
