FROM golang:1.21
ENV DEBIAN_FRONTEND noninteractive
WORKDIR /usr/src/genericAPI
COPY .. /genericAPI
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/app ./...

CMD go run /usr/local/bin/app
EXPOSE 3000/tcp