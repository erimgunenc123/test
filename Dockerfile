FROM golang:1.21
ENV DEBIAN_FRONTEND noninteractive
WORKDIR /project/src/genericAPI
COPY . /project/src/genericAPI
RUN go mod download && go mod verify
RUN cd cmd && go build -o main
#RUN ./main todo add db and local config file to the container
#EXPOSE 3000/tcp