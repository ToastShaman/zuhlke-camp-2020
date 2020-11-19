FROM golang:1.15.5 AS build-stage

ARG VERSION

RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /src
COPY . .

RUN go build \
    -trimpath \
    -ldflags "-s -w -extldflags '-static' -X main.version=${VERSION}" \
    -o /out/todo-server-linux-amd64

RUN strip /out/todo-server-linux-amd64
RUN upx -q -9 /out/todo-server-linux-amd64

FROM scratch AS export-stage
COPY --from=build-stage /out/ /
