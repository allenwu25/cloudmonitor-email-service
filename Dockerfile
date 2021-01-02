FROM golang:1.12.0-alpine3.9
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk update && apk add git
RUN go get -d -v github.com/gorilla/mux
RUN go build -o main .
EXPOSE 8080
CMD ["/app/main"]
