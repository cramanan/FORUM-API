FROM golang:latest

WORKDIR /home/app

COPY . /home/app/

EXPOSE 8080

EXPOSE 8081