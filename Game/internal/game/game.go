package game

import (
	"sync"
	"github.com/gin/websocket"
)

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Snake *Snake
	Send  chan []byte //Канал для отправки сообщений клиенту
}


type Game struct{
	ID string //Room ID

	clients_mutex sync.RWMutex
	Clients map[string]*Client //ID: *Client

	TickRate int64 //update per second 
}
