FROM golang:1.23-alpine

RUN apk add --no-cache curl git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/bin/migrate

RUN chmod +x /app/wait-for.sh
RUN chmod +x /app/start.sh

EXPOSE 3000

ENTRYPOINT ["/app/start.sh"]

CMD [ "/app/main" ]