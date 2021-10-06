# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.17-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /webservice

EXPOSE $PORT

CMD ["/webservice"]


# ... the rest of the Dockerfile is ...
# ...   omitted from this example   ...