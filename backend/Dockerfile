# Используем официальный образ Go в качестве базового
FROM golang:1.21.1-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN go build -o url-shortener ./cmd/shortener/main.go

# Определяем порт, который будет открыт
EXPOSE 8080

# Команда для запуска приложения
CMD ["./url-shortener"]
