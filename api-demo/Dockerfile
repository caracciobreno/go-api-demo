FROM golang:1.14-alpine

WORKDIR /go/src/app

COPY . .

# Install the package
RUN go install -v ./...

# Not a nice thing to do, this is to wait for PG to be ready. Ideally you'd have some real checkings
CMD ["/bin/sh", "-c", "sleep 5; api-demo-service"]