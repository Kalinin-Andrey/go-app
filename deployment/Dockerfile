# README run from root dir (./..)
# docker image build -f deployment/DockerFile -t rbt:1.0 .
# docker container run --publish 8080:8080 --name rbt rbt:1.0
#
#
# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.14

# Copy the local package files to the container's workspace.
#ADD . /go/src/rbt

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go install rbt

# Run the outyet command by default when the container starts.
#ENTRYPOINT /go/bin/rbt

# Document that the service listens on port 8080.
#EXPOSE 8080

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server"]