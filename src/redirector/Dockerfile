FROM golang:1.24-bookworm AS builder

WORKDIR /redirector

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o redirector .

FROM debian:bookworm-slim

COPY --from=builder /redirector/redirector /redirector

EXPOSE 8080

CMD ["/redirector"]
