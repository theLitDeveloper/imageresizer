FROM golang:1.16-alpine AS builder

LABEL maintainer=thelitdeveloper

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /app

ARG aws_access_key_id
ARG aws_secret_access_key
ARG aws_region="eu-central-1"
ARG aws_bucket
ARG redirect_host
ARG port="4321"

ENV CGO_ENABLED=0 
ENV AWS_ACCESS_KEY_ID=${aws_access_key_id} AWS_SECRET_ACCESS_KEY=${aws_secret_access_key}
ENV REDIRECT_HOST=${redirect_host} AWS_BUCKET=${aws_bucket} AWS_REGION=${aws_region} PORT=${port}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN cd cmd/service && go build -v -o main .

FROM alpine:3.14

COPY --from=builder /app/cmd/service/main /main

EXPOSE ${PORT}

CMD ["/main"]
