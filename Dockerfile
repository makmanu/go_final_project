FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/todo_list

CMD ["/app/todo_list"]