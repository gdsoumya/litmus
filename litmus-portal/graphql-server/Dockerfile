# BUILD STAGE
FROM golang:1.16 AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETOS=linux
ARG TARGETARCH

ADD . /gql-server
WORKDIR /gql-server

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}
    
RUN go env
RUN CGO_ENABLED=0 go build -o /output/server -v

## DEPLOY STAGE
FROM alpine:latest

LABEL maintainer="LitmusChaos"

COPY --from=builder /output/server /
COPY --from=builder /gql-server/manifests/. /manifests

RUN addgroup -S litmus && adduser -S -G litmus 1001 
USER 1001

CMD ["./server"]

EXPOSE 8080
