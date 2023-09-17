FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

EXPOSE 8080

CMD  ./geo-service --config config.yaml database migrate -m ./migrations && \
     ./geo-service --config config.yaml serve \