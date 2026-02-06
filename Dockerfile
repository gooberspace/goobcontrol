FROM golang:alpine AS build-stage
WORKDIR /workdir
COPY . ./
RUN apk add --no-cache make git
RUN make get-deps build-linux-amd64
RUN addgroup -S goobcontrol \
    && adduser -S goobcontrol -G goobcontrol

FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /etc/passwd /etc/passwd
USER goobcontrol
COPY --from=build-stage /workdir/bin/goobcontrol_linux_amd64 /goobcontrol
ENTRYPOINT ["/goobcontrol"]

