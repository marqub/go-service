FROM golang:1.11.2 as build

COPY . /go-service
WORKDIR /go-service/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:edge
COPY --from=build /go-service/go-service /go-service

EXPOSE 8080
ENTRYPOINT ["/go-service"]   