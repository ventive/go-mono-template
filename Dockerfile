FROM golang:1.25 as build

ARG GIT_COMMIT=""
ARG APP_VERSION=""
ARG GO_SERVICE=""

RUN mkdir -p /build/
WORKDIR /build/

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
# COPY go.sum .

# This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction below this line)
RUN go mod download

COPY . .

RUN go mod tidy
# RUN go get -v -d ./...

# Build the app and write the executable to /usr/bin/app_service
RUN CGO_ENABLED=0 go build -mod=mod -a -installsuffix cgo --ldflags "-s -w -X ventive.tech/services/$GO_SERVICE/version.GitCommit=$GIT_COMMIT -X ventive.tech/services/$GO_SERVICE/version.AppVersion=$APP_VERSION" -o /usr/bin/app_service cmd/$GO_SERVICE/*.go

# Build the healthchecker and write the executable to /usr/bin/healthchecker

FROM alpine:3.9

ARG GO_SERVICE=""

RUN apk --no-cache add ca-certificates

# Copy app executable
COPY --from=build /usr/bin/app_service /root/
# Copy configuration file
COPY --from=build "/build/services/${GO_SERVICE}/.config.yml" "/root/"

ENV LOGGER_LEVEL="debug"
ENV LOGGER_FORMAT="json"

EXPOSE 6060
EXPOSE 80
ENV PORT 80
WORKDIR /root/

# handle system signals
CMD ["./app_service", "start"]
