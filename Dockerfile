FROM golang:1.17 AS build

RUN mkdir -p /internal-api
WORKDIR /internal-api

COPY go.mod go.sum ./

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM scratch
COPY --from=build /internal-api/internal-api .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./internal-api", "run"]
