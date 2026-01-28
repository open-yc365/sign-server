# build stage
FROM golang:1.23-alpine3.21 AS build
WORKDIR /app
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux
COPY . .
RUN pwd
RUN go version
RUN go mod download
RUN go build -o sign-server ./main.go

# release stage
FROM alpine:3.21 AS release

WORKDIR /app
# install timezone data
RUN apk add -U tzdata
# set timezone to UTC
ENV TZ=UTC

# exe and config file
COPY --from=build /app/sign-server sign-server
COPY --from=build /app/config.json config.json

RUN ls -lh
CMD ["sign-server", ""]
