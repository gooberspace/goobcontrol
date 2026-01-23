FROM golang:1.25-alpine AS build-stage
WORKDIR /workdir
COPY . ./
RUN go mod download
RUN go mod tidy
RUN GOOS=linux go build -o ./goobcontrol

FROM alpine:latest AS release-stage
ARG USERNAME=goobcontrol
RUN adduser -D ${USERNAME}
WORKDIR /app
COPY --from=build-stage /workdir/goobcontrol ./goobcontrol
USER goobcontrol:goobcontrol
ENTRYPOINT ["./goobcontrol"]

