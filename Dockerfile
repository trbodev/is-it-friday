FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o main .

ENV HOST=0.0.0.0

CMD ["/app/main"]
