FROM golang:latest
RUN mkdir /app
ADD ./ /app
WORKDIR /app
RUN go build -o main .

ENTRYPOINT go run api/main.go