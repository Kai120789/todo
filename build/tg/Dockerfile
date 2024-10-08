# Указываем базовый образ
FROM golang:1.23.1-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY ../go.mod ../go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o ./tgbot cmd/tg/tgbot.go

# Указываем команду запуска контейнера
CMD ["./tgbot"]