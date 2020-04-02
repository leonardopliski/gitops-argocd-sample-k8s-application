FROM golang:1.14-alpine3.11

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o server .

EXPOSE 3000

CMD ["./server"]