# syntax=docker/dockerfile:1
FROM golang:1.23.1 AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Собираем статически скомпилированное приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o task-planner cmd/todo_app/main.go

# Используем минимальный образ для запуска
FROM gcr.io/distroless/static

# Копируем бинарный файл из builder
COPY --from=build /app/task-planner /task-planner

# Указываем команду запуска
ENTRYPOINT ["/task-planner"]
