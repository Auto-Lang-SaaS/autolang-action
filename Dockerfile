FROM golang:1.20

WORKDIR /app
COPY . .

RUN go build -o translate-action main.go

ENTRYPOINT ["./translate-action"]
