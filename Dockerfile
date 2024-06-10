FROM golang:latest

WORKDIR /home/app

COPY . /home/app/

RUN go mod download

RUN go build -o server

EXPOSE 8080

ENTRYPOINT [ "./server" ]