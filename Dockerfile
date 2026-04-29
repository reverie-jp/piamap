FROM golang:1.25.6

WORKDIR /piamap-api

RUN apt-get update \
  && apt-get install -y curl xz-utils postgresql-client \
  && rm -rf /var/lib/apt/lists/* \
  && curl -fsSL https://nodejs.org/dist/v24.10.0/node-v24.10.0-linux-arm64.tar.xz \
  | tar -xJ -C /usr/local --strip-components=1

RUN go install github.com/air-verse/air@v1.64.0
RUN go install github.com/bufbuild/buf/cmd/buf@v1.63.0
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
RUN go install github.com/fullstorydev/grpcui/cmd/grpcui@v1.4.3
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
