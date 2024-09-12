FROM golang:alpine as build

COPY . .
RUN go mod download
RUN go build -o /server cmd/server/main.go

FROM alpine:latest

COPY --from=build server /app/server
RUN ls -la && sleep 10

ARG PORT=50051
ARG DISCORD_BOT_TOKEN

ENV PORT=${PORT}
ENV DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN}
EXPOSE ${PORT}
CMD ["/app/server"]
