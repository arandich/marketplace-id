FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./server ./migrator

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY ./migrator/migrations ./migrator/migrations
COPY --from=builder /app/server .

CMD ["./server"]