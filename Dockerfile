FROM golang:1.25-alpine AS build-stage
WORKDIR /workdir
COPY . ./
RUN apk add --no-cache make git
RUN make get-deps build-linux-amd64

FROM alpine:3 AS release-stage
ARG USERNAME=goobcontrol
RUN adduser -D "${USERNAME}"
WORKDIR /app
COPY --from=build-stage /workdir/bin/goobcontrol_linux_amd64 ./goobcontrol
USER goobcontrol:goobcontrol
ENTRYPOINT ["./goobcontrol"]

