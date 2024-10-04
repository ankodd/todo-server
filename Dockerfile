FROM bitnami/golang:1.23.2-debian-12-r0

WORKDIR /app

COPY . .

EXPOSE 8080

# for metrics
EXPOSE 8081

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod tidy

RUN go mod download

RUN mkdir "storage"

RUN mkdir "build"

RUN go build -o ./build/main ./cmd/todo-service/main.go

CMD ["./build/main"]