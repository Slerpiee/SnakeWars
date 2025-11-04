FROM golang:1.24

WORKDIR /app
COPY . .

RUN go mod download


RUN go build -o main ./server

#Подставляет PORT из Dockercompose, который мы передали в environment
CMD ./main --port ${GAME_PORT}