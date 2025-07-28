FROM golang:1.24.5-alpine

WORKDIR  /app

RUN apk add --no-cache curl tar

RUN go install github.com/air-verse/air@latest

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

EXPOSE 4000

CMD ["air"]



