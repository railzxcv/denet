# TODO: добавить мультистейдж сборку
FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o denet ./cmd/denet 
EXPOSE 8080
CMD ["/app/denet", "--env=prod"]