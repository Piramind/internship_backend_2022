FROM golang:latest

WORKDIR /internship_backend_2022

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]