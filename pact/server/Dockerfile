FROM golang:1.15.5c -alpine3.12 AS build-stage

WORKDIR /src
COPY . .

ARG VERSION
RUN GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o /out/otp-darwin-amd64
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o /out/otp-linux-amd64
RUN GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o /out/otp-windows-amd64.exe

FROM scratch AS export-stage
COPY --from=build-stage /out/ /
