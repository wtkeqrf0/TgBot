FROM golang:1.20.2-alpine

WORKDIR /tgbot

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
EXPOSE 3000

CMD ["go", "run","cmd/main/app.go"]