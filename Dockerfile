FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY . .

RUN go build ./cmd/proxy/main.go

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get -y install postgresql postgresql-contrib ca-certificates sudo

USER postgres

COPY /scripts /opt/scripts

RUN service postgresql start && \
        psql -c "CREATE USER admin WITH superuser login password 'admin';" && \
        psql -c "ALTER ROLE admin WITH PASSWORD 'admin';" && \
        createdb -O admin vk && \
        psql -f ./opt/scripts/sql/migrations.sql -d vk

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /proxy
COPY --from=builder /build/main .

COPY . .

EXPOSE 8080
EXPOSE 8000
EXPOSE 5432

ENV PROXY_PORT=8080
ENV REPEATER_PORT=8000
ENV DB_USER=user
ENV DB_NAME=Requests

RUN bash "/opt/scripts/gen.sh"

CMD service postgresql start && ./main