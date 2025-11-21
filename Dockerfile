# dockerfile
FROM golang:latest AS builder

WORKDIR /app
COPY . .

# Statisches Linux-Binary für amd64 bauen
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o gnpmBin .
RUN chmod +x /app/gnpm

FROM alpine:latest AS npm

WORKDIR /root
RUN apk add --no-cache ca-certificates bash
COPY --from=builder /app/gnpmBin /usr/local/bin/gnpm
RUN chmod +x /usr/local/bin/gnpm && chown root:root /usr/local/bin/gnpm

COPY ./exampleApps/npm .
# Quick checks before executing
RUN /usr/local/bin/gnpm install

FROM alpine:latest AS pnpm

WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/gnpmBin /usr/local/bin/gnpm
RUN chmod +x /usr/local/bin/gnpm && chown root:root /usr/local/bin/gnpm

COPY ./exampleApps/pnpm .
RUN /usr/local/bin/gnpm install

FROM alpine:latest AS yarn

WORKDIR /root
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/gnpmBin /usr/local/bin/gnpm
RUN chmod +x /usr/local/bin/gnpm && chown root:root /usr/local/bin/gnpm

COPY ./exampleApps/yarn .
RUN /usr/local/bin/gnpm install

#FROM ubuntu:latest AS deno
#
#WORKDIR /root
#COPY --from=builder /app/gnpmBin /usr/local/bin/gnpm
#RUN chmod +x /usr/local/bin/gnpm && chown root:root /usr/local/bin/gnpm
#
#COPY ./exampleApps/deno .
#
#RUN /usr/local/bin/gnpm install

FROM scratch

COPY --from=npm / /
#COPY --from=pnpm / /
#COPY --from=yarn / /
#COPY --from=deno / /
CMD ["/bin/sh"]
