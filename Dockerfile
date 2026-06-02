FROM golang:1.26.3

RUN apt-get update && apt-get install -y \
    python3 \
    g++

WORKDIR /app

COPY . .

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]