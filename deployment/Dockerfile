# README run from root dir (./..)
# docker image build -f deployment/DockerFile -t app:1.0 .
# docker container run --publish 8080:8080 --name app app:1.0

FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app
COPY --from=builder /go/src/app .

EXPOSE 8080

CMD ["server"]
