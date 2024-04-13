FROM golang:1.22 as builder

RUN mkdir "marketplace-id"


COPY . ./marketplace-id

WORKDIR ./marketplace-id

ENV GOPRIVATE=*
ENV GOINSECURE=*

RUN go build -o ./server ./cmd/server && chmod +x ./cmd
CMD ["./server"]