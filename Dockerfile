# syntax=docker/dockerfile:1
FROM node:16-alpine AS assets
WORKDIR /embed
COPY web ./web
RUN mkdir public
RUN npm --prefix ./web/ install
RUN npm --prefix ./web/ run build

FROM golang:1.18.4-bullseye AS base
WORKDIR /src
COPY go.* ./
RUN go mod download -x
COPY . .
RUN rm -rf public

FROM base AS build
copy --from=assets /embed/public ./public
RUN mkdir /out
RUN go build -v -o /out/felice-n-franz .

FROM debian:bullseye-slim
COPY --from=build /out/felice-n-franz /
EXPOSE 9191
ENTRYPOINT ["/felice-n-franz","-e"]
