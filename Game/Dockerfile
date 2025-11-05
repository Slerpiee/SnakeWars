FROM golang:1.24

WORKDIR /app
COPY . .

RUN go mod download

WORKDIR /app/server

RUN go build main.go

#Подставляет PORT из Dockercompose, который мы передали в environment
CMD ./main --port ${GAME_PORT}