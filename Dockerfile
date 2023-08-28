# We have multi-stage docker build
# The first one is for building image to build the binary file (will be deleted after second second stage complete)
# The second one is for building final image with smaller base image with binary file only

# ==========================================
# 1st Stage
# ==========================================
FROM golang:1.18 AS builder

## Set the working directory
WORKDIR /app

## Copy source
COPY . .

## Compile
RUN make build

# ==========================================
# 2nd Stage
# ==========================================
FROM alpine:latest

ENV APP_NAME=tempo

WORKDIR /app

## Add ssl cert
RUN apk add --update --no-cache ca-certificates

## Add timezone data
RUN apk --no-cache add tzdata

## Copy binary file from 1st stage
COPY --from=builder /app/bin/* ./

## Copy migration files
COPY ./migrations /app/migrations

## Copy env files
COPY ./.env /app/.env
