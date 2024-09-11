FROM golang:latest

COPY . .
RUN go mod download
RUN go build -o server cmd/server/main.go

ARG PORT=50051

ENV PORT=${PORT}
EXPOSE ${PORT}
CMD ["./server"]
