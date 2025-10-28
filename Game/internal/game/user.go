package game

import (
	"github.com/gorilla/websocket"
	"time"
	"sync"
)

type User struct {
	ID    string
	Conn  *websocket.Conn
	Snake *Snake
	SendChan  chan any //Канал для отправки сообщений клиенту
	mutex    sync.RWMutex
    LastPing time.Time
}

