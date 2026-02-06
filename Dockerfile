FROM golang:1.25-alpine AS build-stage
WORKDIR /workdir
COPY . ./
RUN apk add --no-cache make git
RUN make get-deps build-linux-amd64

FROM scratch
COPY --from=build-stage /workdir/bin/goobcontrol_linux_amd64 /goobcontrol
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/goobcontrol"]

