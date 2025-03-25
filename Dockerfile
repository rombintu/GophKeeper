# Используем официальный образ Golang
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарный файл для сервиса
# ARG SERVICE_PATH - путь до main.go сервиса
ARG SERVICE_PATH
RUN go build -o app ./${SERVICE_PATH}

# Команда для запуска приложения
CMD ["./app"]