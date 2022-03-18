FROM golang:1.17-alpine as builder

WORKDIR /src/
COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o server


FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /src/server /app/server

RUN chmod +x /app/server

EXPOSE 80

ENTRYPOINT [ "/app/server" ]