FROM golang:1.22rc1-bullseye AS build

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o meeting-server

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main main

CMD ["/main"]
