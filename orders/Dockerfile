# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /orders ./app/*.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/.env /.env
COPY --from=build /orders /orders

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/orders"]