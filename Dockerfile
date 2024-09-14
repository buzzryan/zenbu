FROM golang:1.23-alpine AS build

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install ./cmd/server

FROM alpine:latest AS app

COPY --from=build /go/bin/* /go/bin/

ENV PATH="${PATH}:/go/bin"

ENTRYPOINT [ "/go/bin/server" ]
